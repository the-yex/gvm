package main

// A simple program demonstrating the spinner component from the Bubbles
// component library.

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/the-yex/gvm/internal/tui/spinner"
	"os"
)

func main() {
	p := tea.NewProgram(spinner.NewSpinner())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
