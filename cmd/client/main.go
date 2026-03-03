package main

import (
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type model struct {
	info          []string
	targetText    string
	typedText     string
	visibleCursor bool
}

type cursorBlink struct{}

func blink() tea.Cmd {
	return tea.Tick(time.Duration(500)*time.Millisecond, func(t time.Time) tea.Msg {
		return cursorBlink{}
	})
}

func initModel() model {
	return model{
		info:          []string{"PunchKeys", "Press ctrl+b to quit"},
		targetText:    "the quick brown fox jumps over the lazy dog",
		typedText:     "",
		visibleCursor: true,
	}
}
func (m model) Init() tea.Cmd {
	return blink()
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
			m.typedText += msg.String()
		}
		if checker(m) {
			return m, tea.Quit
		}
	case cursorBlink:
		m.visibleCursor = !m.visibleCursor
		return m, blink()
	}

	return m, nil
}
func (m model) View() tea.View {
	s := strings.Join(m.info, "\n")
	s += "\n" + rendertext(m.targetText, m.typedText, m.visibleCursor)
	return tea.NewView(s)
}
func rendertext(targetText string, typedText string, visible bool) string {
	correct := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	incorrect := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	extra := lipgloss.NewStyle().Foreground(lipgloss.Color("1")).Faint(true)
	standard := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Faint(true)
	var result strings.Builder
	maxlen := max(len(targetText), len(typedText))
	curpos := len(typedText)
	for i := 0; i < maxlen; i++ {
		if i == curpos && visible {
			result.WriteString("|")
		}
		if i >= len(targetText) && i < len(typedText) {
			result.WriteString(extra.Render(string(typedText[i])))
			continue
		}
		if i < len(targetText) && i < len(typedText) {
			ch := rune(targetText[i])
			if ch == rune(typedText[i]) {
				result.WriteString(correct.Render(string(ch)))
			} else {
				result.WriteString(incorrect.Render(string(ch)))
			}
		} else if i < len(targetText) {
			ch := rune(targetText[i])
			result.WriteString(standard.Render(string(ch)))
		}
	}
	if curpos == maxlen && visible {
		result.WriteString("|")
	}
	return result.String()
}
func checker(m model) bool {
	if m.typedText == m.targetText {
		return true
	}
	return false
}
func main() {
	tea.NewProgram(initModel()).Run()
}
