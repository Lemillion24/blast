package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Lemillion24/blast/internal/network"
	"github.com/Lemillion24/blast/internal/tui/styles"
)

// NetworkPanel affiche les connexions actives avec PID ↔ socket ↔ DNS.
type NetworkPanel struct {
	connections []network.Connection
	filter      string // filtre texte libre (TODO Phase 2)
}

func NewNetworkPanel() NetworkPanel {
	return NetworkPanel{}
}

func (n NetworkPanel) Init() tea.Cmd {
	return fetchConnectionsCmd()
}

// ConnectionsMsg transporte la liste des connexions fraîches.
type ConnectionsMsg []network.Connection

func fetchConnectionsCmd() tea.Cmd {
	return func() tea.Msg {
		conns, err := network.ListConnections()
		if err != nil {
			return nil
		}
		return ConnectionsMsg(conns)
	}
}

func (n NetworkPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case ConnectionsMsg:
		n.connections = []network.Connection(msg)
		return n, nil
	}
	return n, nil
}

func (n NetworkPanel) View() string {
	title := styles.PanelTitle.Render("Surveillance Réseau")
	subtitle := styles.TextMuted.Render("PID ↔ Socket ↔ Destination DNS")

	header := fmt.Sprintf("%-8s %-20s %-22s %-22s %-10s %-30s",
		"PID", "PROCESSUS", "LOCAL", "DISTANT", "ÉTAT", "HOSTNAME")
	header = styles.TextMuted.Render(header)

	var rows []string
	rows = append(rows, header)
	rows = append(rows, strings.Repeat("─", 120))

	for _, c := range n.connections {
		state := c.State
		stateStyle := lipgloss.NewStyle()
		switch state {
		case "ESTABLISHED":
			stateStyle = lipgloss.NewStyle().Foreground(styles.ColorPrimary)
		case "TIME_WAIT", "CLOSE_WAIT":
			stateStyle = lipgloss.NewStyle().Foreground(styles.ColorWarning)
		case "SYN_SENT":
			stateStyle = lipgloss.NewStyle().Foreground(styles.ColorDanger)
		}

		row := fmt.Sprintf("%-8d %-20s %-22s %-22s %-10s %-30s",
			c.PID,
			truncate(c.ProcessName, 20),
			c.LocalAddr,
			c.RemoteAddr,
			stateStyle.Render(truncate(state, 10)),
			truncate(c.Hostname, 30),
		)
		rows = append(rows, row)
	}

	if len(n.connections) == 0 {
		rows = append(rows, styles.TextMuted.Render("  Aucune connexion active..."))
	}

	footer := styles.TextMuted.Render(
		fmt.Sprintf("\n  %d connexion(s) active(s)", len(n.connections)),
	)

	return styles.PanelBorder.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			title, subtitle, "",
			strings.Join(rows, "\n"),
			footer,
		),
	)
}
