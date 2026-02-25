package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	Purple    = lipgloss.Color("#7C3AED")
	Cyan      = lipgloss.Color("#06B6D4")
	Green     = lipgloss.Color("#10B981")
	Red       = lipgloss.Color("#EF4444")
	Yellow    = lipgloss.Color("#F59E0B")
	Gray      = lipgloss.Color("#6B7280")
	DarkGray  = lipgloss.Color("#374151")
	White     = lipgloss.Color("#F9FAFB")
	BgDark    = lipgloss.Color("#111827")
	Pink      = lipgloss.Color("#EC4899")

	// Styles
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Cyan).
			Background(lipgloss.Color("#1E1B4B")).
			Padding(0, 2).
			MarginBottom(1)

	SubtitleStyle = lipgloss.NewStyle().
			Foreground(Purple).
			Bold(true).
			MarginBottom(1)

	InstalledStyle = lipgloss.NewStyle().
			Foreground(Green).
			Strikethrough(true)

	VersionStyle = lipgloss.NewStyle().
			Foreground(Gray).
			Italic(true)

	NotInstalledStyle = lipgloss.NewStyle().
				Foreground(Yellow)

	SelectedStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)

	CursorStyle = lipgloss.NewStyle().
			Foreground(Pink).
			Bold(true)

	CheckedStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	UncheckedStyle = lipgloss.NewStyle().
			Foreground(DarkGray)

	HelpStyle = lipgloss.NewStyle().
			Foreground(Gray).
			MarginTop(1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(Purple).
			Padding(1, 2).
			MarginBottom(1)

	ActiveTabStyle = lipgloss.NewStyle().
			Foreground(White).
			Background(Purple).
			Padding(0, 2).
			Bold(true)

	InactiveTabStyle = lipgloss.NewStyle().
				Foreground(Gray).
				Background(lipgloss.Color("#1F2937")).
				Padding(0, 2)

	ProgressStyle = lipgloss.NewStyle().
			Foreground(Cyan)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(Red).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(Green).
			Bold(true)

	InputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(Cyan).
			Padding(0, 1)

	LabelStyle = lipgloss.NewStyle().
			Foreground(White).
			Bold(true)

	DimStyle = lipgloss.NewStyle().
			Foreground(Gray)

	LogoStyle = lipgloss.NewStyle().
			Foreground(Cyan).
			Bold(true)
)
