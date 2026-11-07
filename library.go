package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/egoist/mygo"
)

// audioExt lists the types the app plays, matching the format packages
// imported in player.go.
var audioExt = map[string]bool{
	".mp3": true, ".m4a": true, ".aac": true,
	".flac": true, ".ogg": true, ".oga": true, ".wav": true,
}

func hasAudioExt(name string) bool {
	return audioExt[strings.ToLower(filepath.Ext(name))]
}

// dislikeDir is where Dislike puts a track: inside the library folder, out of
// the scan, and visible in Finder, so a track comes back by dragging it out.
const dislikeDir = "_dislike"

// ignoredEntry reports the names the library leaves out: the hidden entries
// macOS and other apps keep, and the dislike folder.
func ignoredEntry(name string) bool {
	return strings.HasPrefix(name, ".") || name == dislikeDir
}

// moveToDislike moves a track into the dislike folder of the library, naming
// it apart from the tracks already there. The path is on the same volume, so
// a track the player has open moves without complaint.
func moveToDislike(folder, track string) (string, error) {
	dir := filepath.Join(folder, dislikeDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	base, ext := filepath.Base(track), filepath.Ext(track)
	target := filepath.Join(dir, base)
	for n := 2; ; n++ {
		if _, err := os.Lstat(target); os.IsNotExist(err) {
			break
		}
		if n > 999 {
			return "", fmt.Errorf("%s: too many tracks of that name", base)
		}
		target = filepath.Join(dir, fmt.Sprintf("%s (%d)%s", strings.TrimSuffix(base, ext), n, ext))
	}
	if err := os.Rename(track, target); err != nil {
		return "", err
	}
	return target, nil
}

// scanLibrary lists the audio files of folder, sorted by path so that a
// folder's tracks stay together. With sub, it walks the folders inside.
func scanLibrary(folder string, sub bool) ([]string, error) {
	var files []string
	if sub {
		err := filepath.WalkDir(folder, func(p string, d fs.DirEntry, err error) error {
			if err != nil {
				return nil // unreadable folders are left out, not fatal
			}
			if p == folder {
				return nil
			}
			if ignoredEntry(d.Name()) {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}
			if !d.IsDir() && hasAudioExt(d.Name()) {
				files = append(files, p)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		entries, err := os.ReadDir(folder)
		if err != nil {
			return nil, err
		}
		for _, e := range entries {
			if e.IsDir() || ignoredEntry(e.Name()) || !hasAudioExt(e.Name()) {
				continue
			}
			files = append(files, filepath.Join(folder, e.Name()))
		}
	}
	sort.Strings(files)
	return files, nil
}

func (a *app) loadLibrary() {
	if a.p.Folder == "" {
		a.setLibrary(nil)
		return
	}
	files, err := scanLibrary(a.p.Folder, a.p.Subfolders)
	if err != nil {
		a.missing = true
		a.setLibrary(nil)
		return
	}
	a.missing = false
	a.setLibrary(files)
}

// setLibrary replaces the track list, keeping the track that plays: the
// same path, or the place that path took in the sorted list.
func (a *app) setLibrary(files []string) {
	curPath := a.currentPath()
	playing := a.player != nil

	a.music = files
	a.index = make(map[string]int, len(files))
	for i, p := range files {
		a.index[p] = i
	}
	a.buildQueue()
	a.refilter()

	if len(files) == 0 {
		a.stopPlayer()
		a.cur = -1
		a.refresh()
		return
	}
	if i, ok := a.index[curPath]; ok {
		a.cur = i
		a.refresh()
		return
	}
	pos := sort.SearchStrings(files, curPath)
	if pos >= len(files) {
		pos = len(files) - 1
	}
	a.cur = pos
	if playing {
		a.stopPlayer()
		if !a.playPath(files[pos]) {
			a.advance(false)
		}
	}
	a.refresh()
}

// watchFolder follows the folder, so that added and removed tracks show up
// without a refresh: the watch reports hints, and each batch rereads it.
func (a *app) watchFolder() {
	a.stopWatch()
	if a.p.Folder == "" {
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	a.cancelWatch = cancel

	w, err := mygo.WatchFiles(ctx, a.p.Folder, mygo.FileWatchOptions{
		Recursive: a.p.Subfolders,
		Ignore: func(p string, isDir bool) bool {
			return ignoredEntry(path.Base(p))
		},
	})
	if err != nil {
		a.missing = true
		a.refresh()
		return
	}
	a.missing = false
	go func() {
		defer w.Close()
		for {
			if _, err := w.Next(ctx); err != nil {
				// io.EOF while the watch is still wanted means the folder
				// went away.
				if ctx.Err() == nil {
					a.win.Update(a.onFolderGone)
				}
				return
			}
			a.win.Update(a.loadLibrary)
		}
	}()
}

func (a *app) stopWatch() {
	if a.cancelWatch != nil {
		a.cancelWatch()
		a.cancelWatch = nil
	}
}

func (a *app) onFolderGone() {
	a.missing = true
	a.setLibrary(nil)
}
