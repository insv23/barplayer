package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/egoist/mygo"
)

// Playback orders.
const (
	orderSequential = "sequential"
	orderShuffle    = "shuffle"
	orderRepeatOne  = "repeatOne"
)

// prefs is what the app remembers between launches. The language is not
// here: mygo keeps the app's language in the system, through App.SetLocale.
type prefs struct {
	Folder     string `json:"folder"`
	Order      string `json:"order"`
	Subfolders bool   `json:"subfolders"`
	ShowTitle  bool   `json:"showTitle"`
	Locale     string `json:"locale"`
	PlayPause  string `json:"playPause"`
	Previous   string `json:"previous"`
	Next       string `json:"next"`
	Search     string `json:"search"`
}

func defaultPrefs() prefs {
	return prefs{
		Order:     orderSequential,
		PlayPause: "Alt+P",
		Previous:  "Alt+[",
		Next:      "Alt+]",
		Search:    "Alt+\\",
	}
}

// Shortcuts may be cleared, and an empty string is a value of its own, so
// loading reads them as pointers: absent keeps the default, "" means none.
type prefsFile struct {
	Folder     string  `json:"folder"`
	Order      string  `json:"order"`
	Subfolders bool    `json:"subfolders"`
	ShowTitle  bool    `json:"showTitle"`
	Locale     string  `json:"locale"`
	PlayPause  *string `json:"playPause"`
	Previous   *string `json:"previous"`
	Next       *string `json:"next"`
	Search     *string `json:"search"`
}

func prefsPath() string {
	dir, err := mygo.App.Path(mygo.PathUserData)
	if err != nil {
		return ""
	}
	return filepath.Join(dir, "prefs.json")
}

func loadPrefs() prefs {
	p := defaultPrefs()
	path := prefsPath()
	if path == "" {
		return p
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return p
	}
	var f prefsFile
	if err := json.Unmarshal(data, &f); err != nil {
		log.Printf("prefs: %v", err)
		return p
	}
	p.Folder, p.Subfolders, p.ShowTitle, p.Locale = f.Folder, f.Subfolders, f.ShowTitle, f.Locale
	if f.Order != "" {
		p.Order = f.Order
	}
	if f.PlayPause != nil {
		p.PlayPause = *f.PlayPause
	}
	if f.Previous != nil {
		p.Previous = *f.Previous
	}
	if f.Next != nil {
		p.Next = *f.Next
	}
	if f.Search != nil {
		p.Search = *f.Search
	}
	return p
}

func (a *app) savePrefs() {
	path := prefsPath()
	if path == "" {
		return
	}
	data, err := json.MarshalIndent(a.p, "", "  ")
	if err != nil {
		log.Printf("prefs: %v", err)
		return
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		log.Printf("prefs: %v", err)
		return
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		log.Printf("prefs: %v", err)
	}
}
