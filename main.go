// Barplayer plays the music files of one folder from the menu bar: the tray
// icon opens a menu with play, pause, next, previous, the play order and the
// library folder, and a small window holds the settings.
package main

import (
	"context"
	_ "embed"
	"log"
	"time"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/audio"
	"github.com/egoist/mygo/ui"

	// MyGo's own control strings, dialog buttons and menu roles.
	_ "github.com/egoist/mygo/locales/zh-Hans"
)

// app is the whole state. Tray and menu callbacks, the media keys and the
// views all run on the main thread; goroutines hand work over with
// win.Update.
type app struct {
	p        prefs
	win      *mygo.Window
	tray     *mygo.Tray
	trayMenu *trayMenu
	quitting bool

	music   []string
	index   map[string]int
	cur     int
	missing bool

	player  *audio.Player
	release func()
	playing bool

	// history is every track that started, in order, with hpos the current
	// one, so that previous returns to the track heard before and next goes
	// forward again through the same tracks.
	history []string
	hpos    int
	queue   []string
	qpos    int

	cancelWatch context.CancelFunc

	langChoice  string
	recording   int
	focusRow    int
	shortcutErr string

	startAtLogin bool
	loginErr     string

	lastTitle string
}

func newApp() *app {
	return &app{cur: -1, hpos: -1, recording: -1, index: map[string]int{}}
}

func main() {
	a := newApp()
	mygo.App.SetName("Barplayer")
	// The languages the app says its texts in; mygo resolves the user's
	// preference against them, and App.SetLocale overrides it.
	mygo.App.SetLanguages("en", "zh-Hans")
	mygo.App.SetActivationPolicy(mygo.ActivationPolicyAccessory)
	mygo.App.OnBeforeQuit(func(*mygo.QuitEvent) { a.quitting = true })
	mygo.App.OnWillQuit(func(*mygo.QuitEvent) { a.shutdown() })
	mygo.App.WhenReady(a.ready)
	if err := mygo.App.Run(); err != nil {
		log.Fatal(err)
	}
}

func (a *app) ready() {
	a.p = loadPrefs()
	if a.p.Locale != "" {
		mygo.App.SetLocale(a.p.Locale)
	}
	audio.SetLatency(audio.LatencyPlayback)

	// A tray icon does not keep the app running, and the app quits when its
	// last window closes: the settings window exists from the start, hidden.
	a.win = mygo.NewWindow(mygo.WindowOptions{
		Title:     current().Settings,
		Width:     440,
		Height:    420,
		MinWidth:  380,
		MinHeight: 380,
		Hidden:    true,
		StateKey:  "settings",
		Content:   ui.View(a.settingsView),
	})
	a.win.OnClose(func(e *mygo.CloseEvent) {
		if !a.quitting {
			e.PreventDefault()
			a.win.Hide()
		}
	})

	var err error
	a.tray, err = mygo.NewTray(mygo.TrayOptions{
		ID:             "main",
		Icon:           trayIcon,
		IconIsTemplate: true,
		ToolTip:        "Barplayer",
	})
	if err != nil {
		log.Fatal("tray: ", err)
	}
	a.tray.SetMenu(a.buildMenu())
	mygo.App.SetMenu(a.buildAppMenu())
	a.refresh()

	// Next and previous come from the media keys and the system's controls;
	// Control, in playPath, takes play, pause and seek.
	mygo.NowPlaying.OnCommand(func(cmd mygo.MediaCommand, _ time.Duration) {
		switch cmd {
		case mygo.MediaNext:
			a.next()
		case mygo.MediaPrevious:
			a.previous()
		}
	})

	a.startAtLogin = mygo.App.OpenAtLogin()
	a.shortcutErr = a.bindShortcuts(a.shortcutList())

	// A language picked in the system while the app runs relabels the menus.
	mygo.App.OnLocaleChanged(func(string) {
		a.tray.SetMenu(a.buildMenu())
		mygo.App.SetMenu(a.buildAppMenu())
		a.win.SetTitle(current().Settings)
		a.refresh()
	})

	if a.p.Folder != "" {
		a.loadLibrary()
		a.watchFolder()
	}
}

func (a *app) shutdown() {
	a.stopWatch()
	a.stopPlayer()
	mygo.NowPlaying.Clear()
	if a.tray != nil {
		a.tray.Destroy()
		a.tray = nil
	}
}

// trayIcon is a band rolled up twice, black on transparent: macOS tints a
// template icon to match the menu bar, in light and dark alike. resources/
// trayicon.png holds the @2x bitmap of the same shape.
//
//go:embed resources/trayicon.png
var trayIcon []byte
