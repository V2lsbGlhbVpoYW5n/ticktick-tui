package tui

import "github.com/charmbracelet/lipgloss"

type Color = lipgloss.Color

const (
	BLACK          Color = lipgloss.Color("0")  // Terminal color index 0 (black)
	RED            Color = lipgloss.Color("1")  // Terminal color index 1 (red)
	GREEN          Color = lipgloss.Color("2")  // Terminal color index 2 (green)
	YELLOW         Color = lipgloss.Color("3")  // Terminal color index 3 (yellow)
	BLUE           Color = lipgloss.Color("4")  // Terminal color index 4 (blue)
	MAGENTA        Color = lipgloss.Color("5")  // Terminal color index 5 (magenta)
	CYAN           Color = lipgloss.Color("6")  // Terminal color index 6 (cyan)
	LIGHT_GRAY     Color = lipgloss.Color("7")  // Terminal color index 7 (light gray)
	DARK_GRAY      Color = lipgloss.Color("8")  // Terminal color index 8 (dark gray)
	BRIGHT_RED     Color = lipgloss.Color("9")  // Terminal color index 9 (bright red)
	BRIGHT_GREEN   Color = lipgloss.Color("10") // Terminal color index 10 (bright green)
	BRIGHT_YELLOW  Color = lipgloss.Color("11") // Terminal color index 11 (bright yellow)
	BRIGHT_BLUE    Color = lipgloss.Color("12") // Terminal color index 12 (bright blue)
	BRIGHT_MAGENTA Color = lipgloss.Color("13") // Terminal color index 13 (bright magenta)
	BRIGHT_CYAN    Color = lipgloss.Color("14") // Terminal color index 14 (bright cyan)
	WHITE          Color = lipgloss.Color("15") // Terminal color index 15 (white)
)

var (
	// Status bar styles
	statusBarStyle = lipgloss.NewStyle().
			Foreground(WHITE).
			Height(1)

	statusLeftStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1)

	statusModeStyle = lipgloss.NewStyle().
			Background(GREEN).
			Foreground(BLACK).
			Bold(true).
			Padding(0, 1)

	statusMessageStyle = lipgloss.NewStyle().
				Foreground(WHITE).
				Padding(0, 1)

	statusErrorStyle = lipgloss.NewStyle().
				Background(RED).
				Foreground(WHITE).
				Bold(true).
				Padding(0, 1)

	statusRightStyle = lipgloss.NewStyle().
				Background(LIGHT_GRAY).
				Foreground(BLACK).
				Padding(0, 1)

	// Spinner style
	spinnerStyle = lipgloss.NewStyle().
			Foreground(CYAN)

	// Form styles
	formStyle = lipgloss.NewStyle().
			Padding(2, 4)

	formTitleStyle = lipgloss.NewStyle().
			Foreground(BRIGHT_BLUE).
			Bold(true).
			Margin(0, 0, 1, 0)

	formFocusedStyle = lipgloss.NewStyle().
				Background(LIGHT_GRAY).
				Foreground(BLACK)

	formBlurredStyle = lipgloss.NewStyle().
				Foreground(DARK_GRAY)

	formFocusedPlaceHolderStyle = lipgloss.NewStyle().
					Background(LIGHT_GRAY).
					Foreground(BRIGHT_BLUE)

	formBlurredPlaceHolderStyle = lipgloss.NewStyle().
					Foreground(DARK_GRAY)

	authURLStyle = lipgloss.NewStyle().
			Foreground(CYAN).
			Margin(1, 0)

	// List styles
	listStyle = lipgloss.NewStyle().
			Padding(1, 0)

	listNormalTitleStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(WHITE)).
				Padding(0, 0, 0, 2)

	listNormalDescStyle = listNormalTitleStyle.Copy().
				Foreground(lipgloss.Color(DARK_GRAY))

	listSelectedTitleStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder(), false, false, false, true).
				BorderForeground(lipgloss.Color(BRIGHT_BLUE)).
				Foreground(lipgloss.Color(BRIGHT_BLUE)).
				Padding(0, 0, 0, 1)

	listSelectedDescStyle = listSelectedTitleStyle.Copy().
				Foreground(lipgloss.Color(BLUE))

	// Message styles
	messageStyle = lipgloss.NewStyle().
			Foreground(GREEN).
			Bold(true).
			Padding(0, 1).
			Margin(0, 0, 1, 0)

	messageErrorStyle = lipgloss.NewStyle().
				Foreground(BRIGHT_RED).
				Bold(true).
				Padding(0, 1).
				Margin(0, 0, 1, 0)

	// Help styles
	helpStyle = lipgloss.NewStyle().
			Foreground(DARK_GRAY).
			Padding(0, 1).
			Height(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(DARK_GRAY).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(DARK_GRAY)

	// Paginator styles
	paginatorActive   = lipgloss.NewStyle().Foreground(lipgloss.Color(WHITE)).Render("◈ ")
	paginatorInactive = lipgloss.NewStyle().Foreground(lipgloss.Color(DARK_GRAY)).Render("◇ ")

	// Priority styles
	priorityNone   = lipgloss.NewStyle().Foreground(LIGHT_GRAY).Render("None")
	priorityLow    = lipgloss.NewStyle().Foreground(BRIGHT_BLUE).Render("Low")
	priorityMedium = lipgloss.NewStyle().Foreground(BRIGHT_YELLOW).Render("Medium")
	priorityHigh   = lipgloss.NewStyle().Foreground(BRIGHT_RED).Render("High")
)
