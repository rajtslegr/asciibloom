// Package animation provides terminal manipulation and rendering for ASCII animations.
package animation

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"

	"golang.org/x/term"
)

// Terminal manages raw terminal mode and provides utilities for terminal I/O.
type Terminal struct {
	width         int
	height        int
	origState     *term.State
	file          *os.File
	interrupt     chan struct{}
	interruptOnce sync.Once
}

// NewTerminal initializes a terminal in raw mode for animation.
// Returns an error if not running in a terminal or if raw mode fails.
func NewTerminal() (*Terminal, error) {
	file := os.Stdout
	if !term.IsTerminal(int(file.Fd())) {
		return nil, errors.New("not a terminal")
	}

	origState, err := term.MakeRaw(int(file.Fd()))
	if err != nil {
		return nil, err
	}

	width, height, err := term.GetSize(int(file.Fd()))
	if err != nil {
		_ = term.Restore(int(file.Fd()), origState)
		return nil, err
	}

	t := &Terminal{
		width:     width,
		height:    height,
		origState: origState,
		file:      file,
		interrupt: make(chan struct{}),
	}

	t.clearScreen()
	t.hideCursor()
	t.watchInterrupts()

	return t, nil
}

// InterruptChan returns a channel that receives a signal when Ctrl+C is pressed.
func (t *Terminal) InterruptChan() <-chan struct{} {
	return t.interrupt
}

func (t *Terminal) watchInterrupts() {
	go func() {
		buf := make([]byte, 1)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil || n == 0 {
				return
			}
			if buf[0] == 3 {
				t.interruptOnce.Do(func() {
					close(t.interrupt)
				})
				return
			}
		}
	}()
}

// Width returns the terminal width in columns.
func (t *Terminal) Width() int { return t.width }

// Height returns the terminal height in rows.
func (t *Terminal) Height() int { return t.height }

// Restore returns the terminal to its original state.
func (t *Terminal) Restore() {
	t.showCursor()
	t.clearScreen()
	t.moveCursor(1, 1)
	if t.origState != nil {
		_ = term.Restore(int(t.file.Fd()), t.origState)
	}
}

// Write writes a string to the terminal.
func (t *Terminal) Write(s string) {
	_, _ = t.file.WriteString(s)
}

// WriteAt writes a string to the terminal at the specified position.
func (t *Terminal) WriteAt(s string, x, y int) {
	t.moveCursor(x, y)
	t.Write(s)
}

func (t *Terminal) clearScreen() {
	t.Write("\033[2J")
}

func (t *Terminal) moveCursor(x, y int) {
	t.Write("\033[" + strconv.Itoa(y) + ";" + strconv.Itoa(x) + "H")
}

func (t *Terminal) hideCursor() {
	t.Write("\033[?25l")
}

func (t *Terminal) showCursor() {
	t.Write("\033[?25h")
}

// StringWidth returns the display width of a string in runes.
func (t *Terminal) StringWidth(s string) int {
	return utf8.RuneCountInString(strings.TrimSpace(s))
}
