# Changelog

Notable changes to Barplayer. The format follows
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the version
numbers follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

English · [简体中文](CHANGELOG_zh.md)

## [Unreleased]

### Added

- A library window under **Search…**: every track of the folder in a list, a search field that filters by name as you type, a click to play, and the row that plays highlighted. It hides rather than closes, so the app stays in the menu bar.
- A global shortcut for the library window (`⌥\` by default), shown in the menu and rebindable in the settings window; it opens the window and hides it again while it shows.
- `⌘W` hides the window frontmost, the settings window included.
- The menu shows each control's global shortcut on the right.

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
