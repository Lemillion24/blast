package components

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Lemillion24/blast/internal/forensic"
	"github.com/Lemillion24/blast/internal/tui/styles"
)

// ForensicPanel affiche la timeline des événements et permet l'export.
type ForensicPanel struct {
	events    []forensic.Event
	exporting bool
	lastExport string
}

func NewForensicPanel() ForensicPanel {
	return ForensicPanel{}
}

func (f ForensicPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case forensic.EventMsg:
		f.events = append([]forensic.Event{forensic.Event(msg)}, f.events...)
		if len(f.events) > 500 {
			f.events = f.events[:500]
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "e":
			f.exporting = true
			return f, forensic.ExportJSONCmd(f.events)
		case "E":
			f.exporting = true
			return f, forensic.ExportCSVCmd(f.events)
		}
	case forensic.ExportDoneMsg:
		f.exporting = false
		f.lastExport = string(msg)
	}
	return f, nil
}

func (f ForensicPanel) View() string {
	title := styles.PanelTitle.Render("Forensic & Timeline")

	var rows []string
	header := fmt.Sprintf("%-20s %-12s %-15s %s", "HORODATAGE", "TYPE", "SOURCE", "DÉTAIL")
	rows = append(rows, styles.TextMuted.Render(header))
	rows = append(rows, strings.Repeat("─", 90))

	for i, e := range f.events {
		if i >= 20 {
			rows = append(rows, styles.TextMuted.Render(fmt.Sprintf("  ... %d événement(s) plus anciens", len(f.events)-20)))
			break
		}
		ts := e.Timestamp.Format(time.RFC3339)
		row := fmt.Sprintf("%-20s %-12s %-15s %s",
			ts,
			eventTypeBadge(e.Type),
			truncate(e.Source, 15),
			truncate(e.Detail, 50),
		)
		rows = append(rows, row)
	}

	if len(f.events) == 0 {
		rows = append(rows, styles.TextMuted.Render("  Aucun événement enregistré"))
	}

	exportStatus := ""
	if f.exporting {
		exportStatus = styles.AlertWarning.Render("  Export en cours...")
	} else if f.lastExport != "" {
		exportStatus = styles.TextMuted.Render("  Exporté : " + f.lastExport)
	}

	hint := styles.TextMuted.Render("[e] Export JSON  [E] Export CSV  [c] Effacer")

	return styles.PanelBorder.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			title, "",
			strings.Join(rows, "\n"),
			"",
			exportStatus,
			hint,
		),
	)
}

func eventTypeBadge(t string) string {
	switch t {
	case "ALERT":
		return styles.BadgeDanger.Render("ALERT")
	case "SCAN":
		return lipgloss.NewStyle().
			Background(lipgloss.Color("#1A6B4A")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).Render("SCAN ")
	case "NET":
		return lipgloss.NewStyle().
			Background(lipgloss.Color("#1A3A6B")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).Render("NET  ")
	default:
		return styles.TextMuted.Render("INFO ")
	}
}
