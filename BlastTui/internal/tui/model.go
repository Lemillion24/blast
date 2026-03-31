// Package tui gère l'interface utilisateur terminal de BLAST via Bubbletea.
// Architecture : un modèle racine (AppModel) qui contient N sous-modèles,
// un par panneau. Les messages (tea.Msg) circulent du root vers les panneaux.
package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Lemillion24/blast/internal/tui/components"
	"github.com/Lemillion24/blast/internal/tui/styles"
)

// Tab représente un onglet/panneau du TUI.
type Tab int

const (
	TabMonitor  Tab = iota // CPU, RAM, Disk, processus
	TabNetwork             // trafic réseau, PID↔socket
	TabSecurity            // YARA, règles, alertes
	TabForensic            // timeline, exports
	TabLogs                // logs temps réel
	tabCount
)

var tabNames = []string{
	"  Monitoring ",
	"  Réseau    ",
	" 󰒃 Sécurité  ",
	"  Forensic  ",
	"  Logs      ",
}

// TickMsg est envoyé périodiquement pour rafraîchir les données.
type TickMsg time.Time

// AppModel est le modèle racine de l'application Bubbletea.
type AppModel struct {
	activeTab Tab
	width     int
	height    int

	// Sous-composants (un par panneau)
	monitorPanel  components.MonitorPanel
	networkPanel  components.NetworkPanel
	securityPanel components.SecurityPanel
	forensicPanel components.ForensicPanel
	logsPanel     components.LogsPanel

	// Compteur d'alertes pour la barre de statut
	alertCount int
	ready      bool
}

// New crée et initialise le modèle principal.
func New() AppModel {
	return AppModel{
		activeTab:     TabMonitor,
		monitorPanel:  components.NewMonitorPanel(),
		networkPanel:  components.NewNetworkPanel(),
		securityPanel: components.NewSecurityPanel(),
		forensicPanel: components.NewForensicPanel(),
		logsPanel:     components.NewLogsPanel(),
	}
}

// Start lance la boucle TUI Bubbletea.
func Start() error {
	p := tea.NewProgram(
		New(),
		tea.WithAltScreen(),       // plein écran sans polluer le terminal
		tea.WithMouseCellMotion(), // support souris optionnel
	)
	_, err := p.Run()
	return err
}

// tickCmd envoie un TickMsg toutes les secondes.
func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// Init est appelé une fois au démarrage par Bubbletea.
func (m AppModel) Init() tea.Cmd {
	return tea.Batch(
		tickCmd(),
		m.monitorPanel.Init(),
		m.networkPanel.Init(),
		m.securityPanel.Init(),
	)
}

// Update gère tous les messages entrants (clavier, tick, données).
func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit

		// Navigation entre onglets avec Tab / Shift+Tab
		case "tab":
			m.activeTab = (m.activeTab + 1) % tabCount
		case "shift+tab":
			m.activeTab = (m.activeTab - 1 + tabCount) % tabCount

		// Raccourcis directs par numéro
		case "1":
			m.activeTab = TabMonitor
		case "2":
			m.activeTab = TabNetwork
		case "3":
			m.activeTab = TabSecurity
		case "4":
			m.activeTab = TabForensic
		case "5":
			m.activeTab = TabLogs
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true

	case TickMsg:
		cmds = append(cmds, tickCmd())
		// Envoyer un RefreshMsg au panneau actif pour rafraîchir ses données
		refresh := components.RefreshMsg{}
		switch m.activeTab {
		case TabMonitor:
			updated, cmd := m.monitorPanel.Update(refresh)
			m.monitorPanel = updated.(components.MonitorPanel)
			cmds = append(cmds, cmd)
		case TabNetwork:
			updated, cmd := m.networkPanel.Update(refresh)
			m.networkPanel = updated.(components.NetworkPanel)
			cmds = append(cmds, cmd)
		}
		return m, tea.Batch(cmds...)
	}

	// Propager le message au panneau actif
	switch m.activeTab {
	case TabMonitor:
		updated, cmd := m.monitorPanel.Update(msg)
		m.monitorPanel = updated.(components.MonitorPanel)
		cmds = append(cmds, cmd)
	case TabNetwork:
		updated, cmd := m.networkPanel.Update(msg)
		m.networkPanel = updated.(components.NetworkPanel)
		cmds = append(cmds, cmd)
	case TabSecurity:
		updated, cmd := m.securityPanel.Update(msg)
		m.securityPanel = updated.(components.SecurityPanel)
		cmds = append(cmds, cmd)
	case TabForensic:
		updated, cmd := m.forensicPanel.Update(msg)
		m.forensicPanel = updated.(components.ForensicPanel)
		cmds = append(cmds, cmd)
	case TabLogs:
		updated, cmd := m.logsPanel.Update(msg)
		m.logsPanel = updated.(components.LogsPanel)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View construit le rendu complet de l'écran.
func (m AppModel) View() string {
	if !m.ready {
		return "\n  Initialisation de BLAST..."
	}

	header := m.renderTabs()
	content := m.renderActivePanel()
	statusBar := m.renderStatusBar()

	// Hauteur disponible pour le contenu (total - header - statusbar)
	contentHeight := m.height - lipgloss.Height(header) - lipgloss.Height(statusBar)
	_ = contentHeight // sera utilisé pour contraindre les panneaux

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		statusBar,
	)
}

// renderTabs construit la barre d'onglets.
func (m AppModel) renderTabs() string {
	var tabs []string
	for i, name := range tabNames {
		tab := Tab(i)
		if tab == m.activeTab {
			tabs = append(tabs, styles.PanelTitle.Render(name))
		} else {
			tabs = append(tabs, styles.TextMuted.Render(name))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
}

// renderActivePanel délègue le rendu au panneau actif.
func (m AppModel) renderActivePanel() string {
	switch m.activeTab {
	case TabMonitor:
		return m.monitorPanel.View()
	case TabNetwork:
		return m.networkPanel.View()
	case TabSecurity:
		return m.securityPanel.View()
	case TabForensic:
		return m.forensicPanel.View()
	case TabLogs:
		return m.logsPanel.View()
	}
	return ""
}

// renderStatusBar construit la barre de statut en bas.
func (m AppModel) renderStatusBar() string {
	left := " BLAST v0.1.0 "
	right := " [Tab] Changer | [q] Quitter "

	if m.alertCount > 0 {
		left += styles.BadgeDanger.Render(
			"  " + string(rune('0'+m.alertCount)) + " alerte(s)",
		)
	}

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 0 {
		gap = 0
	}

	return styles.StatusBar.Width(m.width).Render(
		left + lipgloss.NewStyle().Width(gap).Render("") + right,
	)
}
