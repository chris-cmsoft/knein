package picker

import (
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	contexts []string
	query    string
	cursor   int
	limit    int
	selected string
	err      error
}

// SelectContext opens the interactive context picker and returns the selected
// context. An empty string means the user quit without selecting a context.
func SelectContext(contexts []string, limit int) (string, error) {
	initial := model{
		contexts: contexts,
		limit:    limit,
	}

	final, err := tea.NewProgram(initial).Run()
	if err != nil {
		return "", err
	}

	m, ok := final.(model)
	if !ok {
		return "", nil
	}
	if m.err != nil {
		return "", m.err
	}

	return m.selected, nil
}

// OpenK9s starts k9s for context in the current terminal.
func OpenK9s(context string) error {
	cmd := exec.Command("k9s", "--context="+context)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyEnter:
			if selected := m.hovered(); selected != "" {
				m.selected = selected
				return m, tea.Quit
			}
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.matches())-1 {
				m.cursor++
			}
		case tea.KeyBackspace, tea.KeyCtrlH:
			if len(m.query) > 0 {
				m.query = m.query[:len(m.query)-1]
				m.clampCursor()
			}
		case tea.KeySpace:
			m.query += " "
			m.clampCursor()
		default:
			if msg.Type == tea.KeyRunes {
				m.query += msg.String()
				m.clampCursor()
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true)
	mutedStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	selectedStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	cursorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("212"))

	var b strings.Builder
	b.WriteString(titleStyle.Render("Kubernetes context"))
	b.WriteString("\n")
	b.WriteString("> ")
	b.WriteString(m.query)
	b.WriteString(cursorStyle.Render("|"))
	b.WriteString("\n\n")

	options := m.options()
	if len(options) == 0 {
		b.WriteString(mutedStyle.Render("No matching contexts"))
	} else {
		for i, option := range options {
			cursor := "  "
			render := option.Key
			if i == m.cursor {
				cursor = "> "
				render = selectedStyle.Render(option.Key)
			}
			b.WriteString(cursor)
			b.WriteString(render)
			if i < len(options)-1 {
				b.WriteString("\n")
			}
		}
	}

	b.WriteString("\n\n")
	b.WriteString(mutedStyle.Render("Type to filter; spaces combine terms. Use arrows, Enter, Esc."))
	return b.String()
}

func (m model) matches() []string {
	return FilterContexts(m.contexts, m.query, m.limit)
}

func (m model) options() []huh.Option[string] {
	matches := m.matches()
	options := make([]huh.Option[string], 0, len(matches))
	for _, context := range matches {
		options = append(options, huh.NewOption(context, context))
	}
	return options
}

func (m model) hovered() string {
	matches := m.matches()
	if m.cursor < 0 || m.cursor >= len(matches) {
		return ""
	}
	return matches[m.cursor]
}

func (m *model) clampCursor() {
	matches := m.matches()
	if len(matches) == 0 {
		m.cursor = 0
		return
	}
	if m.cursor >= len(matches) {
		m.cursor = len(matches) - 1
	}
}
