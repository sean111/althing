package formatting

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
)

var HeaderStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#aed2f3"))
var LineStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF"))
var MemberNameStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
var ResponseStyle = lipgloss.NewStyle().Faint(true)
var ErrorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Italic(true)

func Line() {
	width, _, err := term.GetSize(0)
	if err != nil {
		width = 100
	}
	var character = "\u2500"
	line := LineStyle.Render(strings.Repeat(character, width))
	fmt.Println(line)
}
