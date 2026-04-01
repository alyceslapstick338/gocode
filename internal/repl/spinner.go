package repl

import (
	"fmt"
	"io"
	"sync"
	"time"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

const spinnerColor = "\033[38;5;39m" // blue
const spinnerReset = "\033[0m"
const clearLine = "\r\033[K"

// Spinner shows an animated loading indicator.
type Spinner struct {
	w       io.Writer
	message string
	stop    chan struct{}
	done    chan struct{}
	mu      sync.Mutex
}

// NewSpinner creates a spinner that writes to w.
func NewSpinner(w io.Writer, message string) *Spinner {
	return &Spinner{
		w:       w,
		message: message,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
}

// Start begins the spinner animation in a goroutine.
func (s *Spinner) Start() {
	go func() {
		defer close(s.done)
		i := 0
		ticker := time.NewTicker(80 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-s.stop:
				fmt.Fprint(s.w, clearLine)
				return
			case <-ticker.C:
				s.mu.Lock()
				frame := spinnerFrames[i%len(spinnerFrames)]
				fmt.Fprintf(s.w, "%s%s%s %s", clearLine, spinnerColor, frame, s.message+spinnerReset)
				s.mu.Unlock()
				i++
			}
		}
	}()
}

// Stop halts the spinner and clears the line.
func (s *Spinner) Stop() {
	close(s.stop)
	<-s.done
}

// UpdateMessage changes the spinner text while running.
func (s *Spinner) UpdateMessage(msg string) {
	s.mu.Lock()
	s.message = msg
	s.mu.Unlock()
}
