package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nossyrrah/blast/internal/security"
	"github.com/nossyrrah/blast/internal/tui/styles"
	"github.com/charmbracelet/lipgloss"
)

// SecurityPanel affiche les alertes YARA, règles actives et processus suspects.
type SecurityPanel struct {
	alerts    []security.Alert
	rules     []security.Rule
	scanState string // "idle" | "scanning" | "done"
}

func NewSecurityPanel() SecurityPanel {
	return SecurityPanel{scanState: "idle"}
}

func (s SecurityPanel) Init() tea.Cmd {
	return security.LoadRulesCmd()
}

// AlertMsg transporte une nouvelle alerte de sécurité.
type AlertMsg security.Alert

// RulesLoadedMsg transporte les règles chargées depuis le disque.
type RulesLoadedMsg []security.Rule

func (s SecurityPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case AlertMsg:
		s.alerts = append([]security.Alert{security.Alert(msg)}, s.alerts...)
		// Garder les 100 dernières alertes
		if len(s.alerts) > 100 {
			s.alerts = s.alerts[:100]
		}
	case RulesLoadedMsg:
		s.rules = []security.Rule(msg)
	case tea.KeyMsg:
		switch msg.String() {
		case "s":
			s.scanState = "scanning"
			return s, security.QuickScanCmd()
		}
	}
	return s, nil
}

func (s SecurityPanel) View() string {
	title := styles.PanelTitle.Render("Sécurité & YARA")

	// Bloc règles actives
	rulesTitle := styles.TextMuted.Render(fmt.Sprintf("Règles chargées : %d", len(s.rules)))
	var rulesList []string
	for i, r := range s.rules {
		if i >= 5 {
			rulesList = append(rulesList, styles.TextMuted.Render(fmt.Sprintf("  ... et %d autre(s)", len(s.rules)-5)))
			break
		}
		rulesList = append(rulesList, fmt.Sprintf("  ✓ %s [%s]", r.Name, r.Source))
	}

	// Bloc alertes
	alertsTitle := styles.AlertCritical.Render(fmt.Sprintf("Alertes (%d)", len(s.alerts)))
	var alertRows []string
	for i, a := range s.alerts {
		if i >= 10 {
			break
		}
		severity := severityBadge(a.Severity)
		alertRows = append(alertRows, fmt.Sprintf("  %s %s → %s", severity, a.RuleName, truncate(a.Target, 40)))
	}
	if len(s.alerts) == 0 {
		alertRows = append(alertRows, styles.TextMuted.Render("  Aucune alerte"))
	}

	hint := styles.TextMuted.Render("[s] Scan rapide  [r] Recharger règles  [k] Kill processus sélectionné")

	return styles.PanelBorder.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			title, "",
			rulesTitle,
			strings.Join(rulesList, "\n"),
			"",
			alertsTitle,
			strings.Join(alertRows, "\n"),
			"",
			hint,
		),
	)
}

func severityBadge(level string) string {
	switch strings.ToLower(level) {
	case "critical":
		return styles.BadgeDanger.Render("CRIT")
	case "high":
		return lipgloss.NewStyle().
			Background(styles.ColorWarning).
			Foreground(lipgloss.Color("#000000")).
			Padding(0, 1).Render("HIGH")
	default:
		return lipgloss.NewStyle().
			Background(styles.ColorMuted).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(0, 1).Render("INFO")
	}
}
