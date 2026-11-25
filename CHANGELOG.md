# Changelog

Notable changes to Barplayer. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the version
numbers follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

English · [简体中文](CHANGELOG_zh.md)

## [Unreleased]

## [0.1.2] - 2026-11-25

### Added

- **Volume**: a slider in the settings window that sets the level of what Barplayer plays, from silence to 200% in 10% steps, remembered across launches. A shared system volume cannot give one app its own level: music sits lower than the videos and calls that the system volume is set for, without turning the whole Mac down and back up again. The system volume keys still change everything. Above 100% a master that already peaks near full scale clips, because the mixer limits at full scale and has no limiter.

## [0.1.1] - 2026-11-07

### Added

- A library window under **Search…**: every track of the folder in a list, a search field that filters by name as you type, a click to play, and the row that plays highlighted. It hides rather than closes, so the app stays in the menu bar.
- A global shortcut for the library window (`⌥\` by default), shown in the menu and rebindable in the settings window; it opens the window and hides it again while it shows.
- `⌘W` hides the window frontmost, the settings window included.
- The menu shows each control's global shortcut on the right.
- A reset button left of each shortcut in the settings window, which puts that one back to the default (`⌥P`, `⌥[`, `⌥]`, `⌥\`). It shows while the row differs from the default.
- Hover tooltips on the reset and clear buttons beside a shortcut, saying what each one does.
- **Dislike**: the track row now opens a submenu, and its first item moves the track that plays into `_dislike` at the top of the library folder. The list moves on without it and nothing is deleted: dragging the file back into the library makes it a track again. It sits one level down, where a slip of the mouse stays clear of Play / Pause.
- Dislike keeps the playback the user asked for: a track that was heard gives way to the next one in the list, and a paused one leaves the list without starting anything. A move that fails (read-only folder, no permission) says so in a dialog and leaves playback and the list as they were. A track inside a subfolder lands at the top level too, and a name already taken becomes `track (2).mp3`.
- Reveal in Finder, in the same submenu, goes to the file of the track that plays.
- The scan and the file watch skip `_dislike`, so changes inside it leave the list still, and a track put there does not come back.

## [0.1.0] - 2026-11-07

The first release.

### Added

- A menu bar player with no Dock icon and no main window: choose the music folder from the menu and it plays.
- MP3, FLAC, M4A (AAC-LC and ALAC), AAC, Ogg Vorbis and WAV.
- A library that follows the folder: tracks added or removed appear without a refresh, and deleting the track that plays moves on to the next one.
- **Include Subfolders**, for a folder laid out by artist or album. Off by default.
- Three play orders: sequential, which loops the library; shuffle, one round without repeats; and repeat one, which repeats only when a track ends by itself.
- Previous walks the playback history, so in shuffle it returns to the track just heard instead of drawing a new one.
- The keyboard's media keys, Control Center, and the system's now playing, including next and previous.
- Global shortcuts for play/pause, previous and next (`⌥P`, `⌥[`, `⌥]` by default), rebindable or clearable from the settings window, with a message when a combination is taken.
- A settings window with the language (System, 简体中文, English), the three shortcuts, start at login, and the track name in the menu bar.
- Quit from the menu, which releases the media keys back to other apps.
