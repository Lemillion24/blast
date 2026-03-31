package components

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"

	"github.com/Lemillion24/blast/internal/tui/styles"
)

// LogsPanel affiche les logs structurés en temps réel avec scroll.
type LogsPanel struct {
	viewport viewport.Model
	lines    []string
	ready    bool
}

// LogLineMsg transporte une nouvelle ligne de log.
type LogLineMsg string

func NewLogsPanel() LogsPanel {
	return LogsPanel{}
}

func (l LogsPanel) Init() tea.Cmd {
	return nil
}

func (l LogsPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if !l.ready {
			l.viewport = viewport.New(msg.Width-4, msg.Height-8)
			l.viewport.SetContent("")
			l.ready = true
		} else {
			l.viewport.Width = msg.Width - 4
			l.viewport.Height = msg.Height - 8
		}

	case LogLineMsg:
		l.lines = append(l.lines, colorizeLog(string(msg)))
		if len(l.lines) > 1000 {
			l.lines = l.lines[len(l.lines)-1000:]
		}
		if l.ready {
			l.viewport.SetContent(strings.Join(l.lines, "\n"))
			l.viewport.GotoBottom()
		}
	}

	if l.ready {
		var cmd tea.Cmd
		l.viewport, cmd = l.viewport.Update(msg)
		return l, cmd
	}
	return l, nil
}

func (l LogsPanel) View() string {
	title := styles.PanelTitle.Render("Logs Temps Réel")
	hint := styles.TextMuted.Render("[↑↓ PgUp PgDn] Scroll  [g] Début  [G] Fin")

	var content string
	if l.ready {
		content = l.viewport.View()
	} else {
		content = styles.TextMuted.Render("  En attente de logs...")
	}

	return styles.PanelBorder.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			title, "",
			content,
			"",
			hint,
		),
	)
}

// colorizeLog applique une couleur selon le niveau de log.
func colorizeLog(line string) string {
	switch {
	case strings.Contains(line, "CRITICAL") || strings.Contains(line, "ERROR"):
		return styles.AlertCritical.Render(line)
	case strings.Contains(line, "WARN"):
		return styles.AlertWarning.Render(line)
	case strings.Contains(line, "DEBUG"):
		return styles.TextMuted.Render(line)
	default:
		return line
	}
}
