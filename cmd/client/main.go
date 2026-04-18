package main

import (
	"fmt"
	"math"
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
	startTime     time.Time
	timeLimit     time.Duration
	testStarted   bool
	timeUp        bool
	remaining     time.Duration
	wpm           float64
	accuracy      float64
	showStats     bool
}

type cursorBlink struct{}

type timer struct{}

func blink() tea.Cmd {
	return tea.Tick(time.Duration(500)*time.Millisecond, func(t time.Time) tea.Msg {
		return cursorBlink{}
	})
}

func timerCalc() tea.Cmd {
	return tea.Tick(100*time.Millisecond, func(t time.Time) tea.Msg {
		return timer{}
	})
}

func initModel() model {
	return model{
		info:          []string{"PunchKeys", "Press ctrl+b to quit"},
		targetText:    "the quick brown fox jumps over the lazy dog",
		typedText:     "",
		visibleCursor: true,
		timeLimit:     15 * time.Second,
		remaining:     15 * time.Second,
	}
}
func (m model) Init() tea.Cmd {
	return blink()
}
func calcStats(m *model) {
	elapsed := time.Since(m.startTime).Seconds()
	if elapsed <= 0 {
		elapsed = 1
	}

	// WPM: (characters / 5) / minutes
	minutes := elapsed / 60.0
	words := float64(len(m.typedText)) / 5.0
	m.wpm = math.Round(words / minutes)

	// Accuracy: correct characters / total typed * 100
	correct := 0
	for i := 0; i < len(m.typedText) && i < len(m.targetText); i++ {
		if m.typedText[i] == m.targetText[i] {
			correct++
		}
	}
	if len(m.typedText) > 0 {
		m.accuracy = math.Round(float64(correct) / float64(len(m.typedText)) * 100)
	} else {
		m.accuracy = 0
	}

	m.timeUp = true
	m.showStats = true
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		// Results screen key handling
		if m.showStats {
			switch msg.String() {
			case "r":
				newModel := initModel()
				return newModel, blink()
			case "q", "ctrl+b":
				return m, tea.Quit
			}
			return m, nil
		}

		if m.timeUp {
			return m, nil
		}

		if !m.testStarted {
			m.testStarted = true
			m.startTime = time.Now()
		}

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
			m.remaining = m.timeLimit - time.Since(m.startTime)
			if m.remaining < 0 {
				m.remaining = 0
			}
			calcStats(&m)
			return m, nil
		}
		if m.testStarted && m.remaining == m.timeLimit {
			return m, timerCalc()
		}
	case timer:
		if m.testStarted && !m.timeUp {
			elapsed := time.Since(m.startTime)
			m.remaining = m.timeLimit - elapsed
			if m.remaining <= 0 {
				m.remaining = 0
				calcStats(&m)
				return m, nil
			}
			return m, timerCalc()
		}

	case cursorBlink:
		if !m.timeUp {
			m.visibleCursor = !m.visibleCursor
			return m, blink()
		}
	}

	return m, nil
}
func (m model) View() tea.View {
	if m.showStats {
		s := "\n  Time's Up!\n\n"
		s += fmt.Sprintf("  WPM:        %.0f\n", m.wpm)
		s += fmt.Sprintf("  Accuracy:   %.0f%%\n", m.accuracy)
		s += fmt.Sprintf("  Characters: %d\n", len(m.typedText))
		s += "\n  Press r to retry\n"
		s += "  Press q to quit\n"
		return tea.NewView(s)
	}

	s := strings.Join(m.info, "\n")
	s += fmt.Sprintf("\nTime remaining: %d", int(m.remaining.Round(time.Second).Seconds()))
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
		if i == curpos {
			if visible {
				result.WriteString("|")
			} else {
				result.WriteString(" ")
			}
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
	if curpos == maxlen {
		if visible {
			result.WriteString("|")
		} else {
			result.WriteString(" ")
		}
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
