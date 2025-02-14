package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	textedit "github.com/jamesrd/textedit-go/textedit"
	"os"
)

func main() {
	var fileName string
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}

	p := tea.NewProgram(textedit.InitModelWithFile(fileName))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
