# Barplayer — working notes for agents

macOS menu bar music player. Go + MyGo native UI (no webview, no cgo).
Product docs are in `README.md` (English) and `README_zh.md`.

## Layout

| Path | What |
| --- | --- |
| `main.go` | app state, startup wiring, tray icon, quit cleanup |
| `player.go` | playback, shuffle queue, history, next/previous, repeat one |
| `library.go` | scanning, FSEvents watching, carrying playback across library changes |
| `menu.go` | tray menu and its updates, menu bar title, application menu |
| `settings.go` | settings window, shortcut recorder, language, login item, folder picker |
| `prefs.go` | preferences (an empty shortcut is a value, decoded through pointers) |
| `i18n.go` | the two string tables |
| `keys.go` | key naming (`ui.Key` first, `ui.KeyCode` fallback), accelerators, macOS symbols |
| `mygo.json` | name, identifier, version, build output |
| `assets/` | README screenshots only (`menu-en.png`, `menu-zh.png`) |
| `TESTING.md` | manual test checklist; this project has no automated tests |

MyGo's own docs are read through `~/.agents/skills/mygo/` (SKILL.md, reference.md,
ui-patterns.md). They cover the APIs this app leans on: tray, menus, Now Playing,
global shortcuts, `WatchFiles`, the shortcut recorder pattern and scroll bars.

## Develop

```sh
go tool mygo dev      # run with live reload
```

The dev build calls `SetName("Barplayer")`, so it shares
`~/Library/Application Support/Barplayer/` with the packaged app: preferences,
window state and language carry across. Delete `prefs.json` to reset.

## Build

```sh
go tool mygo build
```

Produces, for this machine's architecture only (**arm64**; do not build
universal, Apple Silicon is the only target):

```
build/darwin-arm64/Barplayer.app
build/darwin-arm64/barplayer-0.1.0-arm64.dmg
```

The version in the file names comes from `mygo.json`.

Signing is **ad hoc** (`identity "-"`), there is no Developer ID and no
notarization. `codesign -dv` reports `Signature=adhoc`. Consequence for anyone
who downloads a release: macOS quarantines the app, and the install steps must
include

```sh
xattr -dr com.apple.quarantine /Applications/Barplayer.app
```

Never claim the app is signed or notarized in docs or release notes.

## Install locally

```sh
rm -rf /Applications/Barplayer.app
cp -R build/darwin-arm64/Barplayer.app /Applications/
```

`/Applications` is writable by the admin group here, so no `sudo` is needed.

## Release

1. Bump `version` in `mygo.json`.
2. Add the version to `CHANGELOG.md` and `CHANGELOG_zh.md` (Keep a Changelog
   headings, English file first, then translate).
3. `go tool mygo build`.
4. Write the checksum next to the artifacts:

   ```sh
   cd build/darwin-arm64
   shasum -a 256 barplayer-0.1.0-arm64.dmg > barplayer-0.1.0-arm64.dmg.sha256
   ```

5. Create the release with both files attached, tag matching `mygo.json`
   with a `v` prefix:

   ```sh
   gh release create v0.1.0 \
     build/darwin-arm64/barplayer-0.1.0-arm64.dmg \
     build/darwin-arm64/barplayer-0.1.0-arm64.dmg.sha256 \
     --title "Barplayer 0.1.0" --notes-file /tmp/notes.md
   ```

   The notes are the release's CHANGELOG section plus the download steps and the
   `xattr` command.

6. Check that the links in both READMEs point at the new release.

## House rules

- `build/` and `.mygo/` stay out of git (already in `.gitignore`).
- Docs are bilingual: `README.md` and `CHANGELOG.md` are the English originals,
  the `_zh` files are the Chinese versions. Keep them in step.
- Screenshots live in `assets/`, menu only. No settings-window screenshot.
- User-facing strings go through `i18n.go`; never hardcode a label in a view.
- Preferences are JSON with pointer decoding for shortcuts: `""` means cleared
  and must survive a reload.
- The app has no Dock icon (`ActivationPolicyAccessory`) and no main window; the
  settings window is created hidden at startup because a tray icon alone does
  not keep the app alive.
