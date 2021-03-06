package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/zyedidia/clipboard"
	"github.com/zyedidia/tcell"
)

var bindings map[tcell.Key]func(*View) bool

// InitBindings initializes the keybindings for micro
func InitBindings() {
	bindings = make(map[tcell.Key]func(*View) bool)

	actions := map[string]func(*View) bool{
		"CursorUp":            (*View).CursorUp,
		"CursorDown":          (*View).CursorDown,
		"CursorLeft":          (*View).CursorLeft,
		"CursorRight":         (*View).CursorRight,
		"CursorStart":         (*View).CursorStart,
		"CursorEnd":           (*View).CursorEnd,
		"SelectToStart":       (*View).SelectToStart,
		"SelectToEnd":         (*View).SelectToEnd,
		"SelectLeft":          (*View).SelectLeft,
		"SelectRight":         (*View).SelectRight,
		"WordRight":           (*View).WordRight,
		"WordLeft":            (*View).WordLeft,
		"SelectWordRight":     (*View).SelectWordRight,
		"SelectWordLeft":      (*View).SelectWordLeft,
		"SelectToStartOfLine": (*View).SelectToStartOfLine,
		"SelectToEndOfLine":   (*View).SelectToEndOfLine,
		"InsertEnter":         (*View).InsertEnter,
		"InsertSpace":         (*View).InsertSpace,
		"Backspace":           (*View).Backspace,
		"Delete":              (*View).Delete,
		"InsertTab":           (*View).InsertTab,
		"Save":                (*View).Save,
		"Find":                (*View).Find,
		"FindNext":            (*View).FindNext,
		"FindPrevious":        (*View).FindPrevious,
		"Undo":                (*View).Undo,
		"Redo":                (*View).Redo,
		"Copy":                (*View).Copy,
		"Cut":                 (*View).Cut,
		"CutLine":             (*View).CutLine,
		"Paste":               (*View).Paste,
		"SelectAll":           (*View).SelectAll,
		"OpenFile":            (*View).OpenFile,
		"Start":               (*View).Start,
		"End":                 (*View).End,
		"PageUp":              (*View).PageUp,
		"PageDown":            (*View).PageDown,
		"HalfPageUp":          (*View).HalfPageUp,
		"HalfPageDown":        (*View).HalfPageDown,
		"StartOfLine":         (*View).StartOfLine,
		"EndOfLine":           (*View).EndOfLine,
		"ToggleRuler":         (*View).ToggleRuler,
	}

	keys := map[string]tcell.Key{
		"Up":             tcell.KeyUp,
		"Down":           tcell.KeyDown,
		"Right":          tcell.KeyRight,
		"Left":           tcell.KeyLeft,
		"AltUp":          tcell.KeyAltUp,
		"AltDown":        tcell.KeyAltDown,
		"AltLeft":        tcell.KeyAltLeft,
		"AltRight":       tcell.KeyAltRight,
		"CtrlUp":         tcell.KeyCtrlUp,
		"CtrlDown":       tcell.KeyCtrlDown,
		"CtrlLeft":       tcell.KeyCtrlLeft,
		"CtrlRight":      tcell.KeyCtrlRight,
		"ShiftUp":        tcell.KeyShiftUp,
		"ShiftDown":      tcell.KeyShiftDown,
		"ShiftLeft":      tcell.KeyShiftLeft,
		"ShiftRight":     tcell.KeyShiftRight,
		"AltShiftUp":     tcell.KeyAltShiftUp,
		"AltShiftDown":   tcell.KeyAltShiftDown,
		"AltShiftLeft":   tcell.KeyAltShiftLeft,
		"AltShiftRight":  tcell.KeyAltShiftRight,
		"CtrlShiftUp":    tcell.KeyCtrlShiftUp,
		"CtrlShiftDown":  tcell.KeyCtrlShiftDown,
		"CtrlShiftLeft":  tcell.KeyCtrlShiftLeft,
		"CtrlShiftRight": tcell.KeyCtrlShiftRight,
		"UpLeft":         tcell.KeyUpLeft,
		"UpRight":        tcell.KeyUpRight,
		"DownLeft":       tcell.KeyDownLeft,
		"DownRight":      tcell.KeyDownRight,
		"Center":         tcell.KeyCenter,
		"PgUp":           tcell.KeyPgUp,
		"PgDn":           tcell.KeyPgDn,
		"Home":           tcell.KeyHome,
		"End":            tcell.KeyEnd,
		"Insert":         tcell.KeyInsert,
		"Delete":         tcell.KeyDelete,
		"Help":           tcell.KeyHelp,
		"Exit":           tcell.KeyExit,
		"Clear":          tcell.KeyClear,
		"Cancel":         tcell.KeyCancel,
		"Print":          tcell.KeyPrint,
		"Pause":          tcell.KeyPause,
		"Backtab":        tcell.KeyBacktab,
		"F1":             tcell.KeyF1,
		"F2":             tcell.KeyF2,
		"F3":             tcell.KeyF3,
		"F4":             tcell.KeyF4,
		"F5":             tcell.KeyF5,
		"F6":             tcell.KeyF6,
		"F7":             tcell.KeyF7,
		"F8":             tcell.KeyF8,
		"F9":             tcell.KeyF9,
		"F10":            tcell.KeyF10,
		"F11":            tcell.KeyF11,
		"F12":            tcell.KeyF12,
		"F13":            tcell.KeyF13,
		"F14":            tcell.KeyF14,
		"F15":            tcell.KeyF15,
		"F16":            tcell.KeyF16,
		"F17":            tcell.KeyF17,
		"F18":            tcell.KeyF18,
		"F19":            tcell.KeyF19,
		"F20":            tcell.KeyF20,
		"F21":            tcell.KeyF21,
		"F22":            tcell.KeyF22,
		"F23":            tcell.KeyF23,
		"F24":            tcell.KeyF24,
		"F25":            tcell.KeyF25,
		"F26":            tcell.KeyF26,
		"F27":            tcell.KeyF27,
		"F28":            tcell.KeyF28,
		"F29":            tcell.KeyF29,
		"F30":            tcell.KeyF30,
		"F31":            tcell.KeyF31,
		"F32":            tcell.KeyF32,
		"F33":            tcell.KeyF33,
		"F34":            tcell.KeyF34,
		"F35":            tcell.KeyF35,
		"F36":            tcell.KeyF36,
		"F37":            tcell.KeyF37,
		"F38":            tcell.KeyF38,
		"F39":            tcell.KeyF39,
		"F40":            tcell.KeyF40,
		"F41":            tcell.KeyF41,
		"F42":            tcell.KeyF42,
		"F43":            tcell.KeyF43,
		"F44":            tcell.KeyF44,
		"F45":            tcell.KeyF45,
		"F46":            tcell.KeyF46,
		"F47":            tcell.KeyF47,
		"F48":            tcell.KeyF48,
		"F49":            tcell.KeyF49,
		"F50":            tcell.KeyF50,
		"F51":            tcell.KeyF51,
		"F52":            tcell.KeyF52,
		"F53":            tcell.KeyF53,
		"F54":            tcell.KeyF54,
		"F55":            tcell.KeyF55,
		"F56":            tcell.KeyF56,
		"F57":            tcell.KeyF57,
		"F58":            tcell.KeyF58,
		"F59":            tcell.KeyF59,
		"F60":            tcell.KeyF60,
		"F61":            tcell.KeyF61,
		"F62":            tcell.KeyF62,
		"F63":            tcell.KeyF63,
		"F64":            tcell.KeyF64,
		"CtrlSpace":      tcell.KeyCtrlSpace,
		"CtrlA":          tcell.KeyCtrlA,
		"CtrlB":          tcell.KeyCtrlB,
		"CtrlC":          tcell.KeyCtrlC,
		"CtrlD":          tcell.KeyCtrlD,
		"CtrlE":          tcell.KeyCtrlE,
		"CtrlF":          tcell.KeyCtrlF,
		"CtrlG":          tcell.KeyCtrlG,
		"CtrlH":          tcell.KeyCtrlH,
		"CtrlI":          tcell.KeyCtrlI,
		"CtrlJ":          tcell.KeyCtrlJ,
		"CtrlK":          tcell.KeyCtrlK,
		"CtrlL":          tcell.KeyCtrlL,
		"CtrlM":          tcell.KeyCtrlM,
		"CtrlN":          tcell.KeyCtrlN,
		"CtrlO":          tcell.KeyCtrlO,
		"CtrlP":          tcell.KeyCtrlP,
		"CtrlQ":          tcell.KeyCtrlQ,
		"CtrlR":          tcell.KeyCtrlR,
		"CtrlS":          tcell.KeyCtrlS,
		"CtrlT":          tcell.KeyCtrlT,
		"CtrlU":          tcell.KeyCtrlU,
		"CtrlV":          tcell.KeyCtrlV,
		"CtrlW":          tcell.KeyCtrlW,
		"CtrlX":          tcell.KeyCtrlX,
		"CtrlY":          tcell.KeyCtrlY,
		"CtrlZ":          tcell.KeyCtrlZ,
		"CtrlLeftSq":     tcell.KeyCtrlLeftSq,
		"CtrlBackslash":  tcell.KeyCtrlBackslash,
		"CtrlRightSq":    tcell.KeyCtrlRightSq,
		"CtrlCarat":      tcell.KeyCtrlCarat,
		"CtrlUnderscore": tcell.KeyCtrlUnderscore,
		"Backspace":      tcell.KeyBackspace,
		"Tab":            tcell.KeyTab,
		"Esc":            tcell.KeyEsc,
		"Escape":         tcell.KeyEscape,
		"Enter":          tcell.KeyEnter,
		"Space":          tcell.KeySpace,
		"Backspace2":     tcell.KeyBackspace2,
	}

	var parsed map[string]string
	defaults := DefaultBindings()

	filename := configDir + "/bindings.json"
	if _, e := os.Stat(filename); e == nil {
		input, err := ioutil.ReadFile(filename)
		if err != nil {
			TermMessage("Error reading bindings.json file: " + err.Error())
			return
		}

		err = json.Unmarshal(input, &parsed)
		if err != nil {
			TermMessage("Error reading bindings.json:", err.Error())
		}
	}

	for k, v := range defaults {
		bindings[keys[k]] = actions[v]
	}
	for k, v := range parsed {
		bindings[keys[k]] = actions[v]
	}
}

// DefaultBindings returns a map containing micro's default keybindings
func DefaultBindings() map[string]string {
	return map[string]string{
		"Up":             "CursorUp",
		"Down":           "CursorDown",
		"Right":          "CursorRight",
		"Left":           "CursorLeft",
		"ShiftLeft":      "SelectLeft",
		"ShiftRight":     "SelectRight",
		"AltLeft":        "WordLeft",
		"AltRight":       "WordRight",
		"AltShiftRight":  "SelectWordRight",
		"AltShiftLeft":   "SelectWordLeft",
		"CtrlLeft":       "StartOfLine",
		"CtrlRight":      "EndOfLine",
		"CtrlShiftLeft":  "SelectToStartOfLine",
		"CtrlShiftRight": "SelectToEndOfLine",
		"CtrlUp":         "CursorStart",
		"CtrlDown":       "CursorEnd",
		"CtrlShiftUp":    "SelectToStart",
		"CtrlShiftDown":  "SelectToEnd",
		"Enter":          "InsertEnter",
		"Space":          "InsertSpace",
		"Backspace":      "Backspace",
		"Backspace2":     "Backspace",
		"Tab":            "InsertTab",
		"CtrlO":          "OpenFile",
		"CtrlS":          "Save",
		"CtrlF":          "Find",
		"CtrlN":          "FindNext",
		"CtrlP":          "FindPrevious",
		"CtrlZ":          "Undo",
		"CtrlY":          "Redo",
		"CtrlC":          "Copy",
		"CtrlX":          "Cut",
		"CtrlK":          "CutLine",
		"CtrlV":          "Paste",
		"CtrlA":          "SelectAll",
		"Home":           "Start",
		"End":            "End",
		"PgUp":           "PageUp",
		"PgDn":           "PageDown",
		"CtrlU":          "HalfPageUp",
		"CtrlD":          "HalfPageDown",
		"CtrlR":          "ToggleRuler",
		"Delete":         "Delete",
	}
}

// CursorUp moves the cursor up
func (v *View) CursorUp() bool {
	v.cursor.ResetSelection()
	v.cursor.Up()
	return true
}

// CursorDown moves the cursor down
func (v *View) CursorDown() bool {
	v.cursor.ResetSelection()
	v.cursor.Down()
	return true
}

// CursorLeft moves the cursor left
func (v *View) CursorLeft() bool {
	if v.cursor.HasSelection() {
		v.cursor.SetLoc(v.cursor.curSelection[0])
		v.cursor.ResetSelection()
	} else {
		v.cursor.Left()
	}
	return true
}

// CursorRight moves the cursor right
func (v *View) CursorRight() bool {
	if v.cursor.HasSelection() {
		v.cursor.SetLoc(v.cursor.curSelection[1] - 1)
		v.cursor.ResetSelection()
	} else {
		v.cursor.Right()
	}
	return true
}

// WordRight moves the cursor one word to the right
func (v *View) WordRight() bool {
	v.cursor.WordRight()
	return true
}

// WordLeft moves the cursor one word to the left
func (v *View) WordLeft() bool {
	v.cursor.WordLeft()
	return true
}

// SelectLeft selects the character to the left of the cursor
func (v *View) SelectLeft() bool {
	loc := v.cursor.Loc()
	if !v.cursor.HasSelection() {
		v.cursor.origSelection[0] = loc
	}
	v.cursor.SelectTo(loc - 1)
	v.cursor.Left()
	return true
}

// SelectRight selects the character to the right of the cursor
func (v *View) SelectRight() bool {
	loc := v.cursor.Loc()
	if !v.cursor.HasSelection() {
		v.cursor.origSelection[0] = loc
	}
	v.cursor.SelectTo(loc + 1)
	v.cursor.Right()
	return true
}

// SelectWordRight selects the word to the right of the cursor
func (v *View) SelectWordRight() bool {
	loc := v.cursor.Loc()
	if !v.cursor.HasSelection() {
		v.cursor.origSelection[0] = loc
	}
	v.cursor.WordRight()
	v.cursor.SelectTo(v.cursor.Loc())
	return true
}

// SelectWordLeft selects the word to the left of the cursor
func (v *View) SelectWordLeft() bool {
	loc := v.cursor.Loc()
	if !v.cursor.HasSelection() {
		v.cursor.origSelection[0] = loc
	}
	v.cursor.WordLeft()
	v.cursor.SelectTo(v.cursor.Loc())
	return true
}

// StartOfLine moves the cursor to the start of the line
func (v *View) StartOfLine() bool {
	v.cursor.Start()
	return true
}

// EndOfLine moves the cursor to the end of the line
func (v *View) EndOfLine() bool {
	v.cursor.End()
	return true
}

// SelectToStartOfLine selects to the start of the current line
func (v *View) SelectToStartOfLine() bool {
	loc := v.cursor.Loc()
	if !v.cursor.HasSelection() {
		v.cursor.origSelection[0] = loc
	}
	v.cursor.Start()
	v.cursor.SelectTo(v.cursor.Loc())
	return true
}

// SelectToEndOfLine selects to the end of the current line
func (v *View) SelectToEndOfLine() bool {
	loc := v.cursor.Loc()
	if !v.cursor.HasSelection() {
		v.cursor.origSelection[0] = loc
	}
	v.cursor.End()
	v.cursor.SelectTo(v.cursor.Loc())
	return true
}

// CursorStart moves the cursor to the start of the buffer
func (v *View) CursorStart() bool {
	v.cursor.x = 0
	v.cursor.y = 0
	return true
}

// CursorEnd moves the cursor to the end of the buffer
func (v *View) CursorEnd() bool {
	v.cursor.SetLoc(len(v.buf.text))
	return true
}

// SelectToStart selects the text from the cursor to the start of the buffer
func (v *View) SelectToStart() bool {
	loc := v.cursor.Loc()
	if !v.cursor.HasSelection() {
		v.cursor.origSelection[0] = loc
	}
	v.CursorStart()
	v.cursor.SelectTo(0)
	return true
}

// SelectToEnd selects the text from the cursor to the end of the buffer
func (v *View) SelectToEnd() bool {
	loc := v.cursor.Loc()
	if !v.cursor.HasSelection() {
		v.cursor.origSelection[0] = loc
	}
	v.CursorEnd()
	v.cursor.SelectTo(len(v.buf.text))
	return true
}

// InsertSpace inserts a space
func (v *View) InsertSpace() bool {
	// Insert a space
	if v.cursor.HasSelection() {
		v.cursor.DeleteSelection()
		v.cursor.ResetSelection()
	}
	v.eh.Insert(v.cursor.Loc(), " ")
	v.cursor.Right()
	return true
}

// InsertEnter inserts a newline plus possible some whitespace if autoindent is on
func (v *View) InsertEnter() bool {
	// Insert a newline
	if v.cursor.HasSelection() {
		v.cursor.DeleteSelection()
		v.cursor.ResetSelection()
	}

	v.eh.Insert(v.cursor.Loc(), "\n")
	ws := GetLeadingWhitespace(v.buf.lines[v.cursor.y])
	v.cursor.Right()

	if settings.AutoIndent {
		v.eh.Insert(v.cursor.Loc(), ws)
		for i := 0; i < len(ws); i++ {
			v.cursor.Right()
		}
	}
	v.cursor.lastVisualX = v.cursor.GetVisualX()
	return true
}

// Backspace deletes the previous character
func (v *View) Backspace() bool {
	// Delete a character
	if v.cursor.HasSelection() {
		v.cursor.DeleteSelection()
		v.cursor.ResetSelection()
	} else if v.cursor.Loc() > 0 {
		// We have to do something a bit hacky here because we want to
		// delete the line by first moving left and then deleting backwards
		// but the undo redo would place the cursor in the wrong place
		// So instead we move left, save the position, move back, delete
		// and restore the position

		// If the user is using spaces instead of tabs and they are deleting
		// whitespace at the start of the line, we should delete as if its a
		// tab (tabSize number of spaces)
		lineStart := v.buf.lines[v.cursor.y][:v.cursor.x]
		if settings.TabsToSpaces && IsSpaces(lineStart) && len(lineStart) != 0 && len(lineStart)%settings.TabSize == 0 {
			loc := v.cursor.Loc()
			v.cursor.SetLoc(loc - settings.TabSize)
			cx, cy := v.cursor.x, v.cursor.y
			v.cursor.SetLoc(loc)
			v.eh.Remove(loc-settings.TabSize, loc)
			v.cursor.x, v.cursor.y = cx, cy
		} else {
			v.cursor.Left()
			cx, cy := v.cursor.x, v.cursor.y
			v.cursor.Right()
			loc := v.cursor.Loc()
			v.eh.Remove(loc-1, loc)
			v.cursor.x, v.cursor.y = cx, cy
		}
	}
	v.cursor.lastVisualX = v.cursor.GetVisualX()
	return true
}

// Delete deletes the next character
func (v *View) Delete() bool {
	if v.cursor.HasSelection() {
		v.cursor.DeleteSelection()
		v.cursor.ResetSelection()
	} else {
		loc := v.cursor.Loc()
		if loc < len(v.buf.text) {
			v.eh.Remove(loc, loc+1)
		}
	}
	return true
}

// InsertTab inserts a tab or spaces
func (v *View) InsertTab() bool {
	// Insert a tab
	if v.cursor.HasSelection() {
		v.cursor.DeleteSelection()
		v.cursor.ResetSelection()
	}
	if settings.TabsToSpaces {
		v.eh.Insert(v.cursor.Loc(), Spaces(settings.TabSize))
		for i := 0; i < settings.TabSize; i++ {
			v.cursor.Right()
		}
	} else {
		v.eh.Insert(v.cursor.Loc(), "\t")
		v.cursor.Right()
	}
	return true
}

// Save the buffer to disk
func (v *View) Save() bool {
	// If this is an empty buffer, ask for a filename
	if v.buf.path == "" {
		filename, canceled := messenger.Prompt("Filename: ")
		if !canceled {
			v.buf.path = filename
			v.buf.name = filename
		} else {
			return true
		}
	}
	err := v.buf.Save()
	if err != nil {
		messenger.Error(err.Error())
	} else {
		messenger.Message("Saved " + v.buf.path)
		switch v.buf.filetype {
		case "Go":
			v.GoSave()
		}
	}
	return true
}

// GoSave saves the current file (must be a go file) and runs goimports or gofmt
// depending on the user's configuration
func (v *View) GoSave() {
	if settings.GoImports == true {
		messenger.Message("Running goimports...")
		err := goimports(v.buf.path)
		if err != nil {
			messenger.Error(err)
		} else {
			messenger.Message("Saved " + v.buf.path)
		}
		v.reOpen()
	} else if settings.GoFmt == true {
		messenger.Message("Running gofmt...")
		err := gofmt(v.buf.path)
		if err != nil {
			messenger.Error(err)
		} else {
			messenger.Message("Saved " + v.buf.path)
		}
		v.reOpen()
		return
	}

	return
}

// Find opens a prompt and searches forward for the input
func (v *View) Find() bool {
	if v.cursor.HasSelection() {
		searchStart = v.cursor.curSelection[1]
	} else {
		searchStart = ToCharPos(v.cursor.x, v.cursor.y, v.buf)
	}
	BeginSearch()
	return true
}

// FindNext searches forwards for the last used search term
func (v *View) FindNext() bool {
	if v.cursor.HasSelection() {
		searchStart = v.cursor.curSelection[1]
	} else {
		searchStart = ToCharPos(v.cursor.x, v.cursor.y, v.buf)
	}
	messenger.Message("Find: " + lastSearch)
	Search(lastSearch, v, true)
	return true
}

// FindPrevious searches backwards for the last used search term
func (v *View) FindPrevious() bool {
	if v.cursor.HasSelection() {
		searchStart = v.cursor.curSelection[0]
	} else {
		searchStart = ToCharPos(v.cursor.x, v.cursor.y, v.buf)
	}
	messenger.Message("Find: " + lastSearch)
	Search(lastSearch, v, false)
	return true
}

// Undo undoes the last action
func (v *View) Undo() bool {
	v.eh.Undo()
	return true
}

// Redo redoes the last action
func (v *View) Redo() bool {
	v.eh.Redo()
	return true
}

// Copy the selection to the system clipboard
func (v *View) Copy() bool {
	if v.cursor.HasSelection() {
		clipboard.WriteAll(v.cursor.GetSelection())
		v.freshClip = true
	}
	return true
}

// CutLine cuts the current line to the clipboard
func (v *View) CutLine() bool {
	v.cursor.SelectLine()
	if v.freshClip == true {

		if v.cursor.HasSelection() {
			if clip, err := clipboard.ReadAll(); err != nil {
				messenger.Error(err)
			} else {
				clipboard.WriteAll(clip + v.cursor.GetSelection())
			}
		}
	} else if time.Since(v.lastCutTime)/time.Second > 10*time.Second || v.freshClip == false {
		v.Copy()
	}
	v.freshClip = true
	v.lastCutTime = time.Now()
	v.cursor.DeleteSelection()
	v.cursor.ResetSelection()
	return true
}

// Cut the selection to the system clipboard
func (v *View) Cut() bool {
	if v.cursor.HasSelection() {
		clipboard.WriteAll(v.cursor.GetSelection())
		v.cursor.DeleteSelection()
		v.cursor.ResetSelection()
		v.freshClip = true
	}
	return true
}

// Paste whatever is in the system clipboard into the buffer
// Delete and paste if the user has a selection
func (v *View) Paste() bool {
	if v.cursor.HasSelection() {
		v.cursor.DeleteSelection()
		v.cursor.ResetSelection()
	}
	clip, _ := clipboard.ReadAll()
	v.eh.Insert(v.cursor.Loc(), clip)
	v.cursor.SetLoc(v.cursor.Loc() + Count(clip))
	v.freshClip = false
	return true
}

// SelectAll selects the entire buffer
func (v *View) SelectAll() bool {
	v.cursor.curSelection[1] = 0
	v.cursor.curSelection[0] = v.buf.Len()
	// Put the cursor at the beginning
	v.cursor.x = 0
	v.cursor.y = 0
	return true
}

// OpenFile opens a new file in the buffer
func (v *View) OpenFile() bool {
	if v.CanClose("Continue? (yes, no, save) ") {
		filename, canceled := messenger.Prompt("File to open: ")
		if canceled {
			return true
		}
		home, _ := homedir.Dir()
		filename = strings.Replace(filename, "~", home, 1)
		file, err := ioutil.ReadFile(filename)

		if err != nil {
			messenger.Error(err.Error())
			return true
		}
		buf := NewBuffer(string(file), filename)
		v.OpenBuffer(buf)
	}
	return true
}

// Start moves the viewport to the start of the buffer
func (v *View) Start() bool {
	v.topline = 0
	return false
}

// End moves the viewport to the end of the buffer
func (v *View) End() bool {
	if v.height > len(v.buf.lines) {
		v.topline = 0
	} else {
		v.topline = len(v.buf.lines) - v.height
	}
	return false
}

// PageUp scrolls the view up a page
func (v *View) PageUp() bool {
	if v.topline > v.height {
		v.ScrollUp(v.height)
	} else {
		v.topline = 0
	}
	return false
}

// PageDown scrolls the view down a page
func (v *View) PageDown() bool {
	if len(v.buf.lines)-(v.topline+v.height) > v.height {
		v.ScrollDown(v.height)
	} else if len(v.buf.lines) >= v.height {
		v.topline = len(v.buf.lines) - v.height
	}
	return false
}

// HalfPageUp scrolls the view up half a page
func (v *View) HalfPageUp() bool {
	if v.topline > v.height/2 {
		v.ScrollUp(v.height / 2)
	} else {
		v.topline = 0
	}
	return false
}

// HalfPageDown scrolls the view down half a page
func (v *View) HalfPageDown() bool {
	if len(v.buf.lines)-(v.topline+v.height) > v.height/2 {
		v.ScrollDown(v.height / 2)
	} else {
		if len(v.buf.lines) >= v.height {
			v.topline = len(v.buf.lines) - v.height
		}
	}
	return false
}

// ToggleRuler turns line numbers off and on
func (v *View) ToggleRuler() bool {
	if settings.Ruler == false {
		settings.Ruler = true
	} else {
		settings.Ruler = false
	}
	return false
}

// None is no action
func None() bool {
	return false
}

// gofmt runs gofmt on a file
func gofmt(file string) error {
	cmd := exec.Command("gofmt", "-w", file)
	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		return errors.New("Check syntax ") //TODO: highlight or display locations
	}
	return nil
}

// goimports runs goimports on a file
func goimports(file string) error {
	cmd := exec.Command("goimports", "-w", file)
	cmd.Start()
	err := cmd.Wait()
	if err != nil {
		return errors.New("Check syntax ") //TODO: highlight or display locations
	}
	return nil
}
