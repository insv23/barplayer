package main

import (
	"strings"

	"github.com/egoist/mygo"
)

// text holds every string the app says, in one language. The two tables
// below are the whole of the app's translation: mygo resolves the language
// (App.Locale) and persists the choice, and the tables say the words.
type text struct {
	Play, Pause, PlayPause          string
	Next, Previous                  string
	Order, Sequential, Shuffle, One string
	Library, NoFolder, Unavailable  string
	Include, Choose, Reveal         string
	Search, NoMatch, NoTracks       string
	Settings, Quit                  string
	NothingPlaying                  string
	Language, FollowSystem          string
	Shortcuts, General              string
	StartAtLogin, ShowTitle         string
	Record, Unset, PressKeys, Hint  string
	Conflict, Tracks, Clear         string
	ResetDefault                    string
}

var english = text{
	Play: "Play", Pause: "Pause", PlayPause: "Play / Pause",
	Next: "Next", Previous: "Previous",
	Order: "Play Order", Sequential: "In Order", Shuffle: "Shuffle", One: "Repeat One",
	Library: "Library", NoFolder: "No folder chosen", Unavailable: "Folder unavailable",
	Include: "Include Subfolders", Choose: "Choose Folder…", Reveal: "Reveal in Finder",
	Search: "Search", NoMatch: "No track matches", NoTracks: "No tracks in this folder",
	Settings: "Settings", Quit: "Quit",
	NothingPlaying: "Not playing",
	Language:       "Language", FollowSystem: "System",
	Shortcuts: "Keyboard Shortcuts", General: "General",
	StartAtLogin: "Start at login", ShowTitle: "Show track name in the menu bar",
	Record: "Record shortcut", Unset: "Not set", PressKeys: "Press keys…",
	Hint:     "Click a shortcut to change it · Esc cancels · ⌫ clears",
	Conflict: "That shortcut is used by another row.",
	Tracks:   "%d tracks",
	Clear:    "Clear",

	ResetDefault: "Restore default",
}

var chinese = text{
	Play: "播放", Pause: "暂停", PlayPause: "播放 / 暂停",
	Next: "下一首", Previous: "上一首",
	Order: "播放顺序", Sequential: "顺序", Shuffle: "随机", One: "单曲循环",
	Library: "曲库", NoFolder: "未选择文件夹", Unavailable: "文件夹不可用",
	Include: "包含子文件夹", Choose: "选择文件夹…", Reveal: "在 Finder 中显示",
	Search: "搜索", NoMatch: "没有匹配的曲目", NoTracks: "文件夹里没有曲目",
	Settings: "设置", Quit: "退出",
	NothingPlaying: "未在播放",
	Language:       "语言", FollowSystem: "跟随系统",
	Shortcuts: "快捷键", General: "通用",
	StartAtLogin: "登录时启动", ShowTitle: "在菜单栏显示曲目名",
	Record: "录制快捷键", Unset: "未设置", PressKeys: "按下按键…",
	Hint:     "点按快捷键可重新设置 · Esc 取消 · ⌫ 清除",
	Conflict: "这个快捷键已被另一项使用。",
	Tracks:   "%d 首",
	Clear:    "清除",

	ResetDefault: "恢复默认",
}

// current returns the table for the language MyGo resolved for the app.
func current() text {
	if strings.HasPrefix(mygo.App.Locale(), "zh") {
		return chinese
	}
	return english
}

// localeIndex maps the saved choice to the segmented control: 0 follows the
// system, 1 is Chinese, 2 is English.
func localeIndex(tag string) int {
	switch {
	case strings.HasPrefix(tag, "zh"):
		return 1
	case tag != "":
		return 2
	}
	return 0
}
