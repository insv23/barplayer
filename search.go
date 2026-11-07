package main

import (
	"path/filepath"
	"strings"

	"github.com/egoist/mygo/ui"
)

// searchView is the window that lists the library: the field on top filters
// the tracks by name, a click plays one, and the list stays open to scroll.
func (a *app) searchView(c *ui.Context) {
	L := current()

	// Cmd+W hides the window, as the close button does.
	c.OnShortcut(ui.Cmd, ui.KeyW, a.searchWin.Hide)

	// With no query the highlight follows what plays, so a track started from
	// the menu, and the next one after it ends, show here too. With a query
	// the highlight is the user's: the first match, or the row clicked.
	if a.query == "" {
		for i, p := range a.results {
			if p == a.currentPath() {
				a.sel = i
				break
			}
		}
	}

	ui.Column(c).Fill().Children(func() {
		ui.Row(c).Padding(12, 14).Children(func() {
			field := ui.SearchField(c, &a.query).Bind(&a.searchInput).Grow(1).Label(L.Search)
			if field.Changed() {
				a.sel = 0
				a.refilter()
			}
			if field.Submitted() {
				a.playResult(a.sel)
			}
		})
		ui.Divider(c)

		a.searchList.Key = func(i int) any { return a.results[i] }
		a.searchList.Label = func(i int) string { return trackTitle(a.results[i]) }
		a.searchList.Selected = &a.sel
		list := ui.List(c, &a.searchList, len(a.results), func(i int) { a.searchRow(c, i) }).Grow(1).Children(func() {
			ui.Text(c, a.emptyText()).FontSize(12).TextColor(c.Theme().TextMuted).Padding(14)
		})
		if list.Changed() {
			a.playResult(a.sel)
		}
		if a.reveal {
			if a.sel >= 0 {
				a.searchList.ScrollIntoView(a.sel)
			}
			a.reveal = false
		}
	})
}

// searchRow shows a track name, and the folder it sits in when it is not at
// the top level of the library: two files of the same name tell apart so.
func (a *app) searchRow(c *ui.Context, i int) {
	p := a.results[i]
	ui.Column(c).Padding(6, 14).Gap(1).Children(func() {
		ui.Text(c, trackTitle(p)).FontSize(13)
		if dir := a.trackFolder(p); dir != "" {
			ui.Text(c, dir).FontSize(11).TextColor(c.Theme().TextMuted)
		}
	})
}

// emptyText says why the list shows no rows.
func (a *app) emptyText() string {
	L := current()
	switch {
	case a.p.Folder == "":
		return L.NoFolder
	case a.missing:
		return L.Unavailable
	case a.query != "":
		return L.NoMatch
	}
	return L.NoTracks
}

// trackFolder is the folder a track sits in, relative to the library folder,
// empty for a track at its top level.
func (a *app) trackFolder(p string) string {
	if a.p.Folder == "" {
		return ""
	}
	rel, err := filepath.Rel(a.p.Folder, filepath.Dir(p))
	if err != nil || rel == "." || strings.HasPrefix(rel, "..") {
		return ""
	}
	return rel
}

// refilter rebuilds the rows the field's text matches, and keeps the choice
// inside them.
func (a *app) refilter() {
	q := strings.ToLower(a.query)
	a.results = a.results[:0]
	for _, p := range a.music {
		if q == "" || strings.Contains(strings.ToLower(trackTitle(p)), q) {
			a.results = append(a.results, p)
		}
	}
	switch {
	case len(a.results) == 0:
		a.sel = -1
	case a.sel < 0:
		a.sel = 0
	case a.sel >= len(a.results):
		a.sel = len(a.results) - 1
	}
}

// playResult starts the track of row i, and gives the field the keyboard
// back: a list holding it turns typing into a jump to the first name that
// starts with the letter.
func (a *app) playResult(i int) {
	if i < 0 || i >= len(a.results) {
		return
	}
	if a.playPath(a.results[i]) {
		a.searchInput.Focus()
	}
}

// toggleSearch shows the window, or hides it while it shows: the shortcut
// opens and closes it, as a launcher's does.
func (a *app) toggleSearch() {
	if a.searchWin.IsVisible() {
		a.searchWin.Hide()
		return
	}
	a.openSearch()
}

// openSearch shows the window on every track, with the field ready for
// typing and the list at the track that plays.
func (a *app) openSearch() {
	a.query = ""
	a.sel = 0
	a.refilter()
	a.reveal = true
	a.searchWin.Show()
	a.searchWin.Focus()
	a.searchInput.Focus()
}
