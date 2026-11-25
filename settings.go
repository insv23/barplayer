package main

import (
	"fmt"
	"math"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/audio"
	"github.com/egoist/mygo/ui"
)

// clearIcon is the circle with a cross that clears a shortcut, as the erase
// button of macOS fields.
var clearIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round"><circle cx="12" cy="12" r="9"/><path d="m15 9-6 6"/><path d="m9 9 6 6"/></svg>`))

// resetIcon is the counter-clockwise arrow that puts one shortcut back to the
// app's default.
var resetIcon = ui.MustParseSVG([]byte(`<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><path d="M3 12a9 9 0 1 0 9-9 9.75 9.75 0 0 0-6.74 2.74L3 8"/><path d="M3 3v5h5"/></svg>`))

// settingsView is the one window of the app: language, the four shortcuts,
// and the general rows (volume, start at login, the menu bar title), each as
// a row with its label on the left.
func (a *app) settingsView(c *ui.Context) {
	L := current()
	t := c.Theme()

	// Cmd+W hides the window, as the close button does. A key press recorded
	// as a shortcut stays a shortcut.
	if a.recording < 0 {
		c.OnShortcut(ui.Cmd, ui.KeyW, a.win.Hide)
	}

	// The keys held are a fallback for a press that comes before the row has
	// the keyboard: the window reports them whichever element has the focus.
	a.captureHeldKeys(c)

	// The content fits the window, so no scroll bar is wanted. MyGo's overlay
	// indicators would still show while the pointer is over the window;
	// scrolling by wheel and keyboard stays.
	ui.Scroll(c).Fill().Scrollbars(false).Children(func() {
		ui.Column(c).Padding(20, 18).Gap(6).ObserveInput(a.captureInput).Children(func() {
			// The choice is a label of the popup, so it follows the language.
			a.langChoice = langLabel(a.p.Locale)
			a.row(c, L.Language, func() {
				if ui.Select(c, &a.langChoice, langLabels()).Width(150).Changed() {
					a.setLanguage(langIndex(a.langChoice))
				}
			})

			ui.Spacer(c).Height(6)
			section(c, L.Shortcuts)
			def := defaultShortcuts()
			a.shortcutRow(c, 0, L.PlayPause, a.p.PlayPause, def[0])
			a.shortcutRow(c, 1, L.Previous, a.p.Previous, def[1])
			a.shortcutRow(c, 2, L.Next, a.p.Next, def[2])
			a.shortcutRow(c, 3, L.Search, a.p.Search, def[3])
			if a.shortcutErr != "" {
				ui.Text(c, a.shortcutErr).FontSize(11).TextColor(t.Danger)
			}
			ui.Text(c, L.Hint).FontSize(11).TextColor(t.TextMuted)

			ui.Spacer(c).Height(6)
			section(c, L.General)
			a.row(c, L.Volume, func() {
				if ui.StepSlider(c, &a.p.Volume, 0, 2, 0.1).Width(150).Label(L.Volume).Changed() {
					audio.SetVolume(a.p.Volume)
					a.savePrefs()
				}
				ui.Text(c, fmt.Sprintf("%d%%", int(math.Round(a.p.Volume*100)))).
					FontSize(12).Width(40).TextAlign(ui.End).TextColor(t.TextMuted)
			})
			a.row(c, L.StartAtLogin, func() {
				if ui.Switch(c, &a.startAtLogin).Changed() {
					a.setOpenAtLogin(a.startAtLogin)
				}
			})
			if a.loginErr != "" {
				ui.Text(c, a.loginErr).FontSize(11).TextColor(t.Danger)
			}
			a.row(c, L.ShowTitle, func() {
				if ui.Switch(c, &a.p.ShowTitle).Changed() {
					a.savePrefs()
					a.updateTitle()
				}
			})

			ui.Spacer(c).Height(10)
			ui.Divider(c)
			ui.Text(c, versionLine()).FontSize(11).TextColor(t.TextMuted)
		})
	})
}

// versionLine is the app's name, and its version in a packaged build.
func versionLine() string {
	if v := mygo.App.Version(); v != "" {
		return fmt.Sprintf("%s %s", mygo.App.Name(), v)
	}
	return mygo.App.Name()
}

// row puts a label on the left and one control on the right.
func (a *app) row(c *ui.Context, label string, control func()) {
	ui.Row(c).Gap(12).AlignItems(ui.Center).Padding(4, 0).Children(func() {
		ui.Text(c, label).FontSize(13).Grow(1)
		control()
	})
}

func section(c *ui.Context, label string) {
	ui.Text(c, label).FontSize(11).FontWeight(600).TextColor(c.Theme().TextMuted).Padding(2, 0)
}

// shortcutRow is a recorder field: the keys on the right, centered in a
// bordered box as the macOS recorders of other apps show them. A click
// records new ones. The row's two icons sit either side of the field, the
// arrow left of it for the default and the circle right of it for none; each
// keeps its 18 point slot when it does not show, so the fields line up.
func (a *app) shortcutRow(c *ui.Context, i int, label, value, def string) {
	L := current()
	t := c.Theme()
	recording := a.recording == i

	ui.Row(c).Gap(8).AlignItems(ui.Center).Padding(2, 0).
		Focusable().FocusBind(&a.focusRow, i+1).FocusRing(false).
		Children(func() {
			ui.Text(c, label).FontSize(13).Grow(1)

			if value != def && !recording {
				ui.ButtonBase(c).Width(18).Height(18).Radius(9).Label(L.ResetDefault).
					Tooltip(L.ResetDefault).
					OnClick(func() { a.setShortcut(i, def) }).
					Children(func() {
						ui.Icon(c, resetIcon).FontSize(13).TextColor(t.TextMuted)
					})
			} else {
				ui.Box(c).Width(18)
			}

			shown, ink, edge := displayAccel(value), t.Text, t.Border
			switch {
			case recording:
				shown, ink, edge = L.PressKeys, t.Accent, t.Accent
			case value == "":
				shown, ink = L.Unset, t.TextMuted
			}
			ui.ButtonBase(c).Width(104).Height(24).Radius(6).Border(1, edge).
				Label(fmt.Sprintf("%s: %s", L.Record, label)).
				OnClick(func() { a.startRecording(i) }).
				Children(func() {
					ui.Text(c, shown).FontSize(12).TextColor(ink)
				})

			if value != "" && !recording {
				ui.ButtonBase(c).Width(18).Height(18).Radius(9).Label(L.Clear).
					Tooltip(L.Clear).
					OnClick(func() { a.setShortcut(i, "") }).
					Children(func() {
						ui.Icon(c, clearIcon).FontSize(13).TextColor(t.TextMuted)
					})
			} else {
				ui.Box(c).Width(18)
			}
		})
}

// langLabels are the popup's entries, in the language the app shows.
func langLabels() []string {
	L := current()
	return []string{L.FollowSystem, "简体中文", "English"}
}

func langLabel(locale string) string {
	return langLabels()[localeIndex(locale)]
}

// langIndex reads the choice back from the popup's label.
func langIndex(label string) int {
	for i, l := range langLabels() {
		if l == label {
			return i
		}
	}
	return 0
}

func (a *app) startRecording(i int) {
	// The app's own shortcuts stay out of the way while one is recorded, so
	// that pressing a bound combination does not also play or change track.
	mygo.GlobalShortcut.UnregisterAll()
	a.recording = i
	a.shortcutErr = ""
	a.focusRow = i + 1 // gives the row the keyboard, where the keys arrive
}

func (a *app) stopRecording() {
	a.recording = -1
	a.focusRow = 0
	a.shortcutErr = a.bindShortcuts(a.shortcutList())
}

// captureHeldKeys records the keys a window reports held, focused or not, so
// a click and a key pressed in the same instant are not lost.
func (a *app) captureHeldKeys(c *ui.Context) {
	i := a.recording
	if i < 0 {
		return
	}
	mods := c.Modifiers()
	for _, k := range c.KeysDown() {
		name, ok := keyNameFromKey(k)
		if !ok {
			continue
		}
		if mods == 0 {
			switch k {
			case ui.KeyEscape:
				a.stopRecording()
				return
			case ui.KeyBackspace, ui.KeyDelete:
				a.setShortcut(i, "")
				return
			}
		}
		if acc, ok := accelFromName(mods, name); ok {
			a.setShortcut(i, acc)
		}
		return
	}
}

// captureInput takes the keys pressed while a shortcut is recorded. The
// settings column observes them, so the focus has to be inside it: the row
// being recorded holds it, through FocusBind.
func (a *app) captureInput(ev ui.InputEvent) {
	i := a.recording
	if i < 0 || ev.Kind != ui.InputKeyDown || ev.Repeat {
		return
	}
	switch ev.Code {
	case ui.CodeEscape:
		a.stopRecording()
		return
	case ui.CodeBackspace, ui.CodeDelete:
		a.setShortcut(i, "")
		return
	}
	if acc, ok := acceleratorFrom(ev.Mods, ev.Key, ev.Code); ok {
		a.setShortcut(i, acc)
	}
}

func (a *app) shortcutList() []string {
	return []string{a.p.PlayPause, a.p.Previous, a.p.Next, a.p.Search}
}

// setShortcut binds one row, or clears it with "". A combination another row
// uses or the system refuses leaves the old bindings in place.
func (a *app) setShortcut(i int, acc string) {
	prev := a.shortcutList()
	if acc != "" {
		for j, v := range prev {
			if j != i && v == acc {
				a.recording = -1
				a.focusRow = 0
				a.shortcutErr = current().Conflict
				a.bindShortcuts(prev)
				return
			}
		}
	}
	next := append([]string(nil), prev...)
	next[i] = acc
	a.recording = -1
	a.focusRow = 0
	if msg := a.bindShortcuts(next); msg != "" {
		a.shortcutErr = msg
		a.bindShortcuts(prev)
		return
	}
	a.shortcutErr = ""
	a.p.PlayPause, a.p.Previous, a.p.Next = next[0], next[1], next[2]
	a.p.Search = next[3]
	a.savePrefs()
	a.refresh()
}

// bindShortcuts registers the four shortcuts, and says what went wrong when
// one is taken.
func (a *app) bindShortcuts(list []string) string {
	mygo.GlobalShortcut.UnregisterAll()
	entries := []struct {
		acc string
		fn  func()
	}{
		{list[0], a.toggle},
		{list[1], a.previous},
		{list[2], a.next},
		{list[3], a.toggleSearch},
	}
	for _, e := range entries {
		if e.acc == "" {
			continue
		}
		if err := mygo.GlobalShortcut.Register(e.acc, e.fn); err != nil {
			mygo.GlobalShortcut.UnregisterAll()
			return fmt.Sprintf("%s: %v", displayAccel(e.acc), err)
		}
	}
	return ""
}

func (a *app) openSettings() {
	a.win.SetTitle(current().Settings)
	a.win.Show()
	a.win.Focus()
}

// setLanguage switches the language at once; mygo keeps the choice for the
// next launches.
func (a *app) setLanguage(i int) {
	switch i {
	case 1:
		mygo.App.SetLocale("zh-Hans")
		a.p.Locale = "zh-Hans"
	case 2:
		mygo.App.SetLocale("en")
		a.p.Locale = "en"
	default:
		mygo.App.SetLocale("")
		a.p.Locale = ""
	}
	a.savePrefs()

	// The app's own labels need rebuilding; the roles follow by themselves.
	a.tray.SetMenu(a.buildMenu())
	mygo.App.SetMenu(a.buildAppMenu())
	a.win.SetTitle(current().Settings)
	a.refresh()
}

func (a *app) setOrder(order string) {
	a.p.Order = order
	if order == orderShuffle {
		a.buildQueue()
	}
	a.savePrefs()
	a.refresh()
}

func (a *app) setSubfolders(on bool) {
	a.p.Subfolders = on
	a.savePrefs()
	a.loadLibrary()
	a.watchFolder()
	a.refresh()
}

func (a *app) setOpenAtLogin(on bool) {
	if err := mygo.App.SetOpenAtLogin(on); err != nil {
		a.loginErr = err.Error()
		a.startAtLogin = mygo.App.OpenAtLogin()
		return
	}
	a.loginErr = ""
	a.startAtLogin = on
}

// chooseFolder asks for the library folder. The dialog blocks, so it runs on
// a goroutine and hands the answer back to the main thread.
func (a *app) chooseFolder() {
	title := current().Choose
	start := a.p.Folder
	if start == "" {
		start, _ = mygo.App.Path(mygo.PathMusic)
	}
	go func() {
		paths, err := mygo.Dialog.Open(mygo.OpenDialogOptions{
			Title:       title,
			DefaultPath: start,
			Directory:   true,
		})
		if err != nil || len(paths) == 0 {
			return
		}
		a.win.Update(func() { a.setFolder(paths[0]) })
	}()
}

func (a *app) setFolder(folder string) {
	a.p.Folder = folder
	a.missing = false
	a.savePrefs()

	a.stopPlayer()
	a.cur = -1
	a.history = nil
	a.hpos = -1

	// The first track is selected but nothing plays until the user asks.
	a.loadLibrary()
	a.watchFolder()
	a.refresh()
}
