package main

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	info       []string
	targetText string
	typedText  string
}

func initModel() model {
	return model{
		info:       []string{"PunchKeys", "Press q to quit"},
		targetText: "the quick brown fox jumps over the lazy dog",
		typedText:  "",
	}
}
func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+b":
			return m, tea.Quit
		case "enter", "tab", "esc":
			return m, nil
		case "space":
			m.typedText += " "
		case "backspace":
			if len(m.typedText) > 0 {
				m.typedText = m.typedText[:len(m.typedText)-1]
			} else {
				return m, nil
			}
		default:
			if len(m.typedText) < len(m.targetText) {
				m.typedText += msg.String()
			}
		}
		if checker(m) {
			return m, tea.Quit
		}
	}
	return m, nil
}
func (m model) View() tea.View {
	s := strings.Join(m.info, "\n")
	s += "\n\n" + m.targetText
	s += "\n" + rendertext(m.targetText, m.typedText)
	return tea.NewView(s)
}
func rendertext(targetText string, typedText string) string {
	correct := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	incorrect := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	var result strings.Builder
	for i, ch := range typedText {
		if i < len(targetText) && ch == rune(targetText[i]) {
			result.WriteString(correct.Render(string(ch)))
		} else {
			result.WriteString(incorrect.Render(string(ch)))
		}
	}
	return result.String()
}
func checker(m model) bool {
	if len(m.typedText) == len(m.targetText) {
		return true
	}
	return false
}
func main() {
	tea.NewProgram(initModel()).Run()
}
