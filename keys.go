package main

import (
	"strconv"
	"strings"

	"github.com/egoist/mygo/ui"
)

// acceleratorFrom turns a captured key press into the accelerator syntax of
// mygo's global shortcuts, and reports false for presses that cannot be one:
// bare modifiers, keys the syntax has no name for, and Shift alone. The key
// itself names the press; the code is the fallback, since a backend that
// reports no code leaves it unknown.
func acceleratorFrom(mods ui.Modifiers, k ui.Key, code ui.KeyCode) (string, bool) {
	name, ok := keyNameFromKey(k)
	if !ok {
		name, ok = keyName(code)
	}
	if !ok {
		return "", false
	}
	return accelFromName(mods, name)
}

// accelFromName joins the modifiers held and a key the syntax knows.
func accelFromName(mods ui.Modifiers, name string) (string, bool) {
	var parts []string
	if mods&ui.Super != 0 {
		parts = append(parts, "Cmd")
	}
	if mods&ui.Ctrl != 0 {
		parts = append(parts, "Ctrl")
	}
	if mods&ui.Alt != 0 {
		parts = append(parts, "Alt")
	}
	// Shift alone would fire on every capital letter, so it needs a modifier
	// that makes the combination a shortcut.
	if len(parts) == 0 {
		return "", false
	}
	if mods&ui.Shift != 0 {
		parts = append(parts, "Shift")
	}
	parts = append(parts, name)
	return strings.Join(parts, "+"), true
}

// keyName names the keys the accelerator syntax knows, by their position on
// the keyboard rather than by what the layout types.
func keyName(code ui.KeyCode) (string, bool) {
	switch code {
	case ui.CodeKeyA:
		return "A", true
	case ui.CodeKeyB:
		return "B", true
	case ui.CodeKeyC:
		return "C", true
	case ui.CodeKeyD:
		return "D", true
	case ui.CodeKeyE:
		return "E", true
	case ui.CodeKeyF:
		return "F", true
	case ui.CodeKeyG:
		return "G", true
	case ui.CodeKeyH:
		return "H", true
	case ui.CodeKeyI:
		return "I", true
	case ui.CodeKeyJ:
		return "J", true
	case ui.CodeKeyK:
		return "K", true
	case ui.CodeKeyL:
		return "L", true
	case ui.CodeKeyM:
		return "M", true
	case ui.CodeKeyN:
		return "N", true
	case ui.CodeKeyO:
		return "O", true
	case ui.CodeKeyP:
		return "P", true
	case ui.CodeKeyQ:
		return "Q", true
	case ui.CodeKeyR:
		return "R", true
	case ui.CodeKeyS:
		return "S", true
	case ui.CodeKeyT:
		return "T", true
	case ui.CodeKeyU:
		return "U", true
	case ui.CodeKeyV:
		return "V", true
	case ui.CodeKeyW:
		return "W", true
	case ui.CodeKeyX:
		return "X", true
	case ui.CodeKeyY:
		return "Y", true
	case ui.CodeKeyZ:
		return "Z", true
	case ui.CodeDigit0:
		return "0", true
	case ui.CodeDigit1:
		return "1", true
	case ui.CodeDigit2:
		return "2", true
	case ui.CodeDigit3:
		return "3", true
	case ui.CodeDigit4:
		return "4", true
	case ui.CodeDigit5:
		return "5", true
	case ui.CodeDigit6:
		return "6", true
	case ui.CodeDigit7:
		return "7", true
	case ui.CodeDigit8:
		return "8", true
	case ui.CodeDigit9:
		return "9", true
	case ui.CodeMinus:
		return "-", true
	case ui.CodeEqual:
		return "=", true
	case ui.CodeBracketLeft:
		return "[", true
	case ui.CodeBracketRight:
		return "]", true
	case ui.CodeBackslash:
		return "\\", true
	case ui.CodeSemicolon:
		return ";", true
	case ui.CodeQuote:
		return "'", true
	case ui.CodeBackquote:
		return "`", true
	case ui.CodeComma:
		return ",", true
	case ui.CodePeriod:
		return ".", true
	case ui.CodeSlash:
		return "/", true
	case ui.CodeSpace:
		return "Space", true
	case ui.CodeEnter:
		return "Enter", true
	case ui.CodeEscape:
		return "Escape", true
	case ui.CodeBackspace:
		return "Backspace", true
	case ui.CodeTab:
		return "Tab", true
	case ui.CodeDelete:
		return "Delete", true
	case ui.CodeInsert:
		return "Insert", true
	case ui.CodeHome:
		return "Home", true
	case ui.CodeEnd:
		return "End", true
	case ui.CodePageUp:
		return "PageUp", true
	case ui.CodePageDown:
		return "PageDown", true
	case ui.CodeArrowUp:
		return "Up", true
	case ui.CodeArrowDown:
		return "Down", true
	case ui.CodeArrowLeft:
		return "Left", true
	case ui.CodeArrowRight:
		return "Right", true
	case ui.CodeF1:
		return "F1", true
	case ui.CodeF2:
		return "F2", true
	case ui.CodeF3:
		return "F3", true
	case ui.CodeF4:
		return "F4", true
	case ui.CodeF5:
		return "F5", true
	case ui.CodeF6:
		return "F6", true
	case ui.CodeF7:
		return "F7", true
	case ui.CodeF8:
		return "F8", true
	case ui.CodeF9:
		return "F9", true
	case ui.CodeF10:
		return "F10", true
	case ui.CodeF11:
		return "F11", true
	case ui.CodeF12:
		return "F12", true
	case ui.CodeF13:
		return "F13", true
	case ui.CodeF14:
		return "F14", true
	case ui.CodeF15:
		return "F15", true
	case ui.CodeF16:
		return "F16", true
	case ui.CodeF17:
		return "F17", true
	case ui.CodeF18:
		return "F18", true
	case ui.CodeF19:
		return "F19", true
	case ui.CodeF20:
		return "F20", true
	}
	return "", false
}

// displayAccel writes an accelerator the way macOS shows keys: ⌥P, ⌘⇧K.
func displayAccel(acc string) string {
	if acc == "" {
		return ""
	}
	glyphs := map[string]string{
		"Cmd": "⌘", "Command": "⌘", "Super": "⌘", "Meta": "⌘", "Win": "⌘",
		"Ctrl": "⌃", "Control": "⌃",
		"Alt": "⌥", "Option": "⌥",
		"Shift": "⇧",
	}
	keys := map[string]string{
		"Enter": "↩", "Return": "↩", "Tab": "⇥", "Escape": "⎋", "Esc": "⎋",
		"Space": "Space", "Backspace": "⌫", "Delete": "⌦", "Insert": "⌅",
		"Up": "↑", "Down": "↓", "Left": "←", "Right": "→",
		"Home": "↖", "End": "↘", "PageUp": "⇞", "PageDown": "⇟",
	}
	var out strings.Builder
	for _, part := range strings.Split(acc, "+") {
		if g, ok := glyphs[part]; ok {
			out.WriteString(g)
			continue
		}
		if k, ok := keys[part]; ok {
			out.WriteString(k)
			continue
		}
		out.WriteString(strings.ToUpper(part))
	}
	return out.String()
}

// keyNameFromKey names the keys the accelerator syntax knows, from the key a
// keyboard reported, whatever layout it types.
func keyNameFromKey(k ui.Key) (string, bool) {
	switch k {
	case ui.KeyEnter:
		return "Enter", true
	case ui.KeyEscape:
		return "Escape", true
	case ui.KeyBackspace:
		return "Backspace", true
	case ui.KeyTab:
		return "Tab", true
	case ui.KeySpace:
		return "Space", true
	case ui.KeyDelete:
		return "Delete", true
	case ui.KeyInsert:
		return "Insert", true
	case ui.KeyHome:
		return "Home", true
	case ui.KeyEnd:
		return "End", true
	case ui.KeyPageUp:
		return "PageUp", true
	case ui.KeyPageDown:
		return "PageDown", true
	case ui.KeyLeft:
		return "Left", true
	case ui.KeyRight:
		return "Right", true
	case ui.KeyUp:
		return "Up", true
	case ui.KeyDown:
		return "Down", true
	case ui.KeyMinus:
		return "-", true
	case ui.KeyEqual:
		return "=", true
	case ui.KeyComma:
		return ",", true
	case ui.KeyPeriod:
		return ".", true
	case ui.KeySlash:
		return "/", true
	case ui.KeySemicolon:
		return ";", true
	case ui.KeyQuote:
		return "'", true
	case ui.KeyBracketLeft:
		return "[", true
	case ui.KeyBracketRight:
		return "]", true
	case ui.KeyBackslash:
		return "\\", true
	case ui.KeyBackquote:
		return "`", true
	case ui.Key0:
		return "0", true
	case ui.Key1:
		return "1", true
	case ui.Key2:
		return "2", true
	case ui.Key3:
		return "3", true
	case ui.Key4:
		return "4", true
	case ui.Key5:
		return "5", true
	case ui.Key6:
		return "6", true
	case ui.Key7:
		return "7", true
	case ui.Key8:
		return "8", true
	case ui.Key9:
		return "9", true
	}
	if k >= ui.KeyA && k <= ui.KeyZ {
		return string(rune('A' + (k - ui.KeyA))), true
	}
	if k >= ui.KeyF1 && k <= ui.KeyF12 {
		return "F" + strconv.Itoa(int(k-ui.KeyF1)+1), true
	}
	return "", false
}
