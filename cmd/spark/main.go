package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/charmbracelet/bubbletea"
	"github.com/dpeluche/spark/internal/tui"
)

func main() {
	// FORCE LOGGING FOR DEBUGGING
	f, err := tea.LogToFile("spark_debug.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	// Panic Recovery
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("CRITICAL PANIC:", r)
			fmt.Println("Stack Trace:")
			debug.PrintStack()
			os.Exit(1)
		}
	}()

	p := tea.NewProgram(tui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

	fmt.Println("\n  See you later, Space Cowboy... 🚀")
	fmt.Println("  Spark sequence complete.\n")
}
