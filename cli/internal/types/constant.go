package types

import "github.com/charmbracelet/lipgloss"

const (
	Padding         = 2
	MaxProcessWidth = 80
)

var (
	HelpStyle           = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render
	CurrentPkgNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("211"))
	DoneStyle           = lipgloss.NewStyle().Margin(1, 2)
	CheckMark           = lipgloss.NewStyle().Foreground(lipgloss.Color("42")).SetString("✓")
)
