package main

import (
	"math/rand"
	"os"
	"path/filepath"
	"strings"

	"github.com/egoist/mygo"
	"github.com/egoist/mygo/audio"

	// The formats the app plays. Each decodes a file type for package audio.
	_ "github.com/egoist/mygo/audio/aac"
	_ "github.com/egoist/mygo/audio/flac"
	_ "github.com/egoist/mygo/audio/mp3"
	_ "github.com/egoist/mygo/audio/mp4"
	_ "github.com/egoist/mygo/audio/vorbis"
)

func (a *app) currentPath() string {
	if a.cur >= 0 && a.cur < len(a.music) {
		return a.music[a.cur]
	}
	return ""
}

func trackTitle(path string) string {
	base := filepath.Base(path)
	return strings.TrimSuffix(base, filepath.Ext(base))
}

// toggle plays the selected track, or pauses and resumes the one playing.
func (a *app) toggle() {
	if len(a.music) == 0 {
		return
	}
	if a.player == nil {
		path := a.currentPath()
		if path == "" {
			path = a.music[0]
		}
		a.playPath(path)
		return
	}
	if a.playing {
		a.player.Pause()
		a.playing = false
	} else {
		a.player.Play()
		a.playing = true
	}
	mygo.NowPlaying.Sync()
	a.refresh()
}

func (a *app) next() {
	if len(a.music) == 0 {
		return
	}
	if a.player == nil {
		// Nothing plays: move the selection the menu shows.
		if path, _ := a.nextTrack(); path != "" {
			if i, ok := a.index[path]; ok {
				a.cur = i
				a.refresh()
			}
		}
		return
	}
	a.advance(false)
}

// previous steps back through what played, which is what makes a shuffle
// reversible: it returns to the track heard before, not to a new one.
func (a *app) previous() {
	if len(a.music) == 0 {
		return
	}
	if a.player == nil {
		if a.hpos > 0 {
			a.hpos--
			if i, ok := a.index[a.history[a.hpos]]; ok {
				a.cur = i
			}
		} else if a.cur > 0 {
			a.cur--
		}
		a.refresh()
		return
	}
	for i := a.hpos - 1; i >= 0; i-- {
		if _, ok := a.index[a.history[i]]; !ok {
			continue
		}
		a.hpos = i
		if a.playKnown(a.history[i]) {
			return
		}
	}
	// Nothing played before this track: begin it again.
	a.player.Seek(0)
	mygo.NowPlaying.Sync()
	a.refresh()
}

// dislike moves the track that plays into the library's dislike folder and
// takes it off the list. Playback follows what the user asked for: a track
// that was heard gives way to the next one, and a paused one leaves the list
// without starting anything.
func (a *app) dislike() {
	track := a.currentPath()
	if track == "" {
		return
	}
	audible := a.playing
	if _, err := moveToDislike(a.p.Folder, track); err != nil {
		mygo.Dialog.Error(current().DislikeFailed, err.Error())
		return
	}
	a.stopPlayer()
	a.loadLibrary() // the track is gone: the list moves on to the one after it
	if !audible {
		return
	}
	if path := a.currentPath(); path != "" && !a.playPath(path) {
		a.advance(false)
	}
}

// advance moves to the next track. auto is true when the current track
// ended, which is the only case that Repeat One repeats.
func (a *app) advance(auto bool) {
	if len(a.music) == 0 {
		return
	}
	if auto && a.p.Order == orderRepeatOne && a.player != nil {
		a.player.Play()
		a.playing = true
		mygo.NowPlaying.Sync()
		a.refresh()
		return
	}
	// Tracks may vanish between the scan and the click; skip the ones that
	// no longer open.
	for i := 0; i < 8; i++ {
		path, known := a.nextTrack()
		if path == "" {
			return
		}
		if known {
			if a.playKnown(path) {
				return
			}
		} else if a.playPath(path) {
			return
		}
		a.dropPath(path)
	}
}

// nextTrack names the next track, and tells whether it comes from the
// history the app already played, as after Previous.
func (a *app) nextTrack() (string, bool) {
	if len(a.music) == 0 {
		return "", false
	}
	if a.p.Order == orderShuffle {
		// After going back, a shuffle plays forward through what played.
		for i := a.hpos + 1; i < len(a.history); i++ {
			if _, ok := a.index[a.history[i]]; ok {
				a.hpos = i
				return a.history[i], true
			}
		}
		return a.fromQueue(), false
	}
	if a.cur < 0 {
		return a.music[0], false
	}
	return a.music[(a.cur+1)%len(a.music)], false
}

// fromQueue takes the next track of the shuffled round, shuffling again once
// the round is spent and never repeating the track that just played.
func (a *app) fromQueue() string {
	for a.qpos < len(a.queue) {
		p := a.queue[a.qpos]
		a.qpos++
		if _, ok := a.index[p]; ok {
			return p
		}
	}
	cur := a.currentPath()
	a.buildQueue()
	if len(a.queue) > 1 && a.queue[0] == cur {
		a.queue[0], a.queue[1] = a.queue[1], a.queue[0]
	}
	if a.qpos < len(a.queue) {
		p := a.queue[a.qpos]
		a.qpos++
		return p
	}
	return ""
}

func (a *app) buildQueue() {
	a.queue = append([]string(nil), a.music...)
	rand.Shuffle(len(a.queue), func(i, j int) {
		a.queue[i], a.queue[j] = a.queue[j], a.queue[i]
	})
	a.qpos = 0
}

func (a *app) dropPath(path string) {
	delete(a.index, path)
	kept := a.music[:0]
	for _, p := range a.music {
		if p != path {
			kept = append(kept, p)
		}
	}
	a.music = kept
	for i, p := range a.music {
		a.index[p] = i
	}
	if a.cur >= len(a.music) {
		a.cur = len(a.music) - 1
	}
}

// playPath opens a track, starts it, and records it as played.
func (a *app) playPath(path string) bool {
	if !a.start(path) {
		return false
	}
	a.pushHistory(path)
	return true
}

// playKnown starts a track the history holds, as going back to it does, and
// leaves the tracks after it in place so that next returns through them.
func (a *app) playKnown(path string) bool {
	if !a.start(path) {
		return false
	}
	if i, ok := a.index[path]; ok {
		a.cur = i
	}
	return true
}

func (a *app) start(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	player, err := audio.OpenStream(f)
	if err != nil {
		f.Close()
		return false
	}

	a.stopPlayer()

	a.player = player
	a.playing = true
	player.OnEnd(func() {
		a.win.Update(func() { a.advance(true) })
	})

	mygo.NowPlaying.Set(mygo.NowPlayingInfo{
		Title:    trackTitle(path),
		Duration: player.Duration(),
		Commands: mygo.MediaPlay | mygo.MediaPause | mygo.MediaNext | mygo.MediaPrevious,
	})
	// Control takes play, pause, stop and seek, and shows the player's
	// state; next and previous still arrive at the OnCommand above.
	a.release = mygo.NowPlaying.Control(player)
	player.Play()
	mygo.NowPlaying.Sync()

	if i, ok := a.index[path]; ok {
		a.cur = i
	}
	a.refresh()
	return true
}

func (a *app) stopPlayer() {
	if a.release != nil {
		a.release()
		a.release = nil
	}
	if a.player != nil {
		a.player.Close()
		a.player = nil
	}
	a.playing = false
}

func (a *app) pushHistory(path string) {
	if a.hpos >= 0 && a.hpos+1 < len(a.history) {
		a.history = a.history[:a.hpos+1]
	}
	a.history = append(a.history, path)
	a.hpos = len(a.history) - 1
}
