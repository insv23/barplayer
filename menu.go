package main

import (
	"fmt"
	"strings"

	"github.com/egoist/mygo"
)

// trayMenu keeps the items whose label, check or enabled state follows the
// app, so that each change reaches the native menu without rebuilding it.
type trayMenu struct {
	now, toggle, prev, next       *mygo.MenuItem
	dislike, revealTrack          *mygo.MenuItem
	search                        *mygo.MenuItem
	seq, shuffle, one             *mygo.MenuItem
	path, include, choose, reveal *mygo.MenuItem
	settings                      *mygo.MenuItem
}

func (a *app) buildMenu() *mygo.Menu {
	L := current()
	m := &trayMenu{
		dislike: &mygo.MenuItem{ID: "dislike", Label: L.Dislike, Click: func(*mygo.MenuItem, *mygo.Window) {
			a.dislike()
		}},
		revealTrack: &mygo.MenuItem{ID: "reveal.track", Label: L.Reveal, Click: func(*mygo.MenuItem, *mygo.Window) {
			if p := a.currentPath(); p != "" {
				mygo.Shell.ShowItemInFolder(p)
			}
		}},
		toggle: &mygo.MenuItem{ID: "toggle", Accelerator: a.p.PlayPause, Click: func(*mygo.MenuItem, *mygo.Window) {
			a.toggle()
		}},
		prev: &mygo.MenuItem{ID: "previous", Accelerator: a.p.Previous, Click: func(*mygo.MenuItem, *mygo.Window) {
			a.previous()
		}},
		next: &mygo.MenuItem{ID: "next", Accelerator: a.p.Next, Click: func(*mygo.MenuItem, *mygo.Window) {
			a.next()
		}},
		search: &mygo.MenuItem{ID: "search", Label: L.Search + "…", Click: func(*mygo.MenuItem, *mygo.Window) {
			a.openSearch()
		}},
		seq: &mygo.MenuItem{ID: "order.sequential", Type: mygo.MenuItemRadio, Click: func(*mygo.MenuItem, *mygo.Window) {
			a.setOrder(orderSequential)
		}},
		shuffle: &mygo.MenuItem{ID: "order.shuffle", Type: mygo.MenuItemRadio, Click: func(*mygo.MenuItem, *mygo.Window) {
			a.setOrder(orderShuffle)
		}},
		one: &mygo.MenuItem{ID: "order.one", Type: mygo.MenuItemRadio, Click: func(*mygo.MenuItem, *mygo.Window) {
			a.setOrder(orderRepeatOne)
		}},
		path: &mygo.MenuItem{ID: "library.path", Disabled: true},
		include: &mygo.MenuItem{ID: "library.subfolders", Type: mygo.MenuItemCheckbox, Label: L.Include,
			Click: func(item *mygo.MenuItem, _ *mygo.Window) {
				a.setSubfolders(item.Checked)
			}},
		choose: &mygo.MenuItem{ID: "library.choose", Label: L.Choose, Click: func(*mygo.MenuItem, *mygo.Window) {
			a.chooseFolder()
		}},
		reveal: &mygo.MenuItem{ID: "library.reveal", Label: L.Reveal, Click: func(*mygo.MenuItem, *mygo.Window) {
			if a.p.Folder != "" {
				mygo.Shell.ShowItemInFolder(a.p.Folder)
			}
		}},
		settings: &mygo.MenuItem{ID: "settings", Label: L.Settings + "…", Click: func(*mygo.MenuItem, *mygo.Window) {
			a.openSettings()
		}},
	}
	m.seq.Label, m.shuffle.Label, m.one.Label = L.Sequential, L.Shuffle, L.One

	a.trayMenu = m
	// The track row opens a submenu: what acts on the track sits under its
	// name, where a slip of the mouse does not reach Play/Pause.
	m.now = &mygo.MenuItem{ID: "now", Submenu: []*mygo.MenuItem{m.dislike, m.revealTrack}}
	menu := mygo.NewMenu([]*mygo.MenuItem{
		m.now,
		mygo.Separator(),
		m.toggle,
		m.prev,
		m.next,
		m.search,
		mygo.Separator(),
		{Label: L.Order, Submenu: []*mygo.MenuItem{m.seq, m.shuffle, m.one}},
		{Label: L.Library, Submenu: []*mygo.MenuItem{
			m.path,
			m.include,
			mygo.Separator(),
			m.choose,
			m.reveal,
		}},
		m.settings,
		mygo.Separator(),
		{Label: L.Quit, Click: func(*mygo.MenuItem, *mygo.Window) { mygo.App.Quit() }},
	})
	a.syncMenu()
	return menu
}

// buildAppMenu is the menu bar the app shows while its settings window is
// frontmost: Quit, Settings and the editing roles for its fields.
func (a *app) buildAppMenu() *mygo.Menu {
	L := current()
	return mygo.NewMenu([]*mygo.MenuItem{
		{Label: mygo.App.Name(), Submenu: []*mygo.MenuItem{
			{Role: mygo.RoleAbout},
			mygo.Separator(),
			{Label: L.Settings + "…", Accelerator: "Cmd+,", Click: func(*mygo.MenuItem, *mygo.Window) {
				a.openSettings()
			}},
			mygo.Separator(),
			{Role: mygo.RoleHide},
			{Role: mygo.RoleQuit},
		}},
		{Role: mygo.RoleEditMenu},
	})
}

// syncMenu brings the tray menu and the menu bar title up to date with what
// plays and what the library holds.
func (a *app) syncMenu() {
	m := a.trayMenu
	if m == nil {
		return
	}
	L := current()
	empty := len(a.music) == 0

	switch {
	case empty:
		m.now.SetLabel(L.NothingPlaying)
		m.now.SetSubtitle("")
	case a.cur >= 0:
		m.now.SetLabel(trackTitle(a.music[a.cur]))
		m.now.SetSubtitle(fmt.Sprintf("%d / %d", a.cur+1, len(a.music)))
	}

	label := L.Play
	if a.playing {
		label = L.Pause
	}
	m.toggle.SetLabel(label)
	m.toggle.SetEnabled(!empty)
	m.now.SetEnabled(!empty)
	m.dislike.SetEnabled(!empty)
	m.revealTrack.SetEnabled(!empty)
	m.prev.SetLabel(L.Previous)
	m.prev.SetEnabled(!empty)
	m.next.SetLabel(L.Next)
	m.next.SetEnabled(!empty)
	m.search.SetEnabled(!empty)

	// The shortcuts show beside the four, as macOS draws them, and go away
	// when one is cleared. A menu does not answer them while it is open:
	// macOS runs its own loop then and holds back the global shortcuts too.
	m.toggle.SetAccelerator(a.p.PlayPause)
	m.prev.SetAccelerator(a.p.Previous)
	m.next.SetAccelerator(a.p.Next)
	m.search.SetAccelerator(a.p.Search)

	m.seq.SetChecked(a.p.Order == orderSequential)
	m.shuffle.SetChecked(a.p.Order == orderShuffle)
	m.one.SetChecked(a.p.Order == orderRepeatOne)
	m.include.SetChecked(a.p.Subfolders)

	switch {
	case a.p.Folder == "":
		m.path.SetLabel(L.NoFolder)
	case a.missing:
		m.path.SetLabel(L.Unavailable)
	default:
		m.path.SetLabel(fmt.Sprintf("%s · %s", shortPath(a.p.Folder), fmt.Sprintf(L.Tracks, len(a.music))))
	}
	m.reveal.SetEnabled(a.p.Folder != "" && !a.missing)

	a.updateTitle()
}

// updateTitle writes the track name beside the icon, when the setting asks
// for it and the name changed: a title that changes width moves the icons
// beside it.
func (a *app) updateTitle() {
	if a.tray == nil {
		return
	}
	title := ""
	if a.p.ShowTitle && len(a.music) > 0 && a.cur >= 0 {
		title = truncate(trackTitle(a.music[a.cur]), 24)
	}
	if title == a.lastTitle {
		return
	}
	a.lastTitle = title
	a.tray.SetTitle(title)
}

func (a *app) refresh() {
	a.syncMenu()
}

// shortPath writes a path with the home folder as ~.
func shortPath(p string) string {
	home, err := mygo.App.Path(mygo.PathHome)
	if err == nil && home != "" && strings.HasPrefix(p, home) {
		return "~" + strings.TrimPrefix(p, home)
	}
	return p
}

func truncate(s string, n int) string {
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	return string(runes[:n-1]) + "…"
}
