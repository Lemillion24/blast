package components

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/nossyrrah/blast/internal/monitor"
	"github.com/nossyrrah/blast/internal/tui/styles"
)

// MonitorPanel affiche les métriques système : CPU, RAM, Disk, processus.
type MonitorPanel struct {
	stats monitor.SystemStats
}

func NewMonitorPanel() MonitorPanel {
	return MonitorPanel{}
}

func (m MonitorPanel) Init() tea.Cmd {
	return fetchStatsCmd()
}

// StatsMsg transporte les métriques fraîches vers le panneau.
type StatsMsg monitor.SystemStats

// fetchStatsCmd collecte les métriques système en arrière-plan.
func fetchStatsCmd() tea.Cmd {
	return func() tea.Msg {
		stats, err := monitor.Collect()
		if err != nil {
			return nil
		}
		return StatsMsg(stats)
	}
}

func (m MonitorPanel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case StatsMsg:
		m.stats = monitor.SystemStats(msg)
		return m, nil
	// Rafraîchir à chaque tick
	case interface{ isTickMsg() }:
		return m, fetchStatsCmd()
	}
	return m, nil
}

func (m MonitorPanel) View() string {
	title := styles.PanelTitle.Render("Monitoring Système")

	cpu := renderBar("CPU", m.stats.CPUPercent, 100, "%")
	ram := renderBar("RAM", float64(m.stats.MemUsed), float64(m.stats.MemTotal), "Go")
	disk := renderBar("Disk", float64(m.stats.DiskUsed), float64(m.stats.DiskTotal), "Go")

	procs := renderProcessTable(m.stats.TopProcesses)

	return styles.PanelBorderActive.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			cpu, ram, disk,
			"",
			styles.PanelTitle.Render("Top Processus"),
			procs,
		),
	)
}

// renderBar génère une barre de progression textuelle colorée.
func renderBar(label string, value, max float64, unit string) string {
	if max == 0 {
		max = 1
	}
	percent := value / max
	barWidth := 30
	filled := int(percent * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}

	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	color := styles.ColorPrimary
	if percent > 0.85 {
		color = styles.ColorDanger
	} else if percent > 0.65 {
		color = styles.ColorWarning
	}

	barStyled := lipgloss.NewStyle().Foreground(color).Render(bar)
	valueStr := fmt.Sprintf("%.1f/%-.1f %s", value, max, unit)

	return fmt.Sprintf("%-6s [%s] %s", label, barStyled, styles.TextMuted.Render(valueStr))
}

// renderProcessTable construit le tableau des processus.
func renderProcessTable(procs []monitor.ProcessInfo) string {
	if len(procs) == 0 {
		return styles.TextMuted.Render("  Aucun processus...")
	}

	header := fmt.Sprintf("%-8s %-20s %6s %8s", "PID", "NOM", "CPU%", "MEM(Mo)")
	header = styles.TextMuted.Render(header)

	var rows []string
	rows = append(rows, header)

	for _, p := range procs {
		cpuColor := lipgloss.NewStyle()
		if p.CPUPercent > 80 {
			cpuColor = lipgloss.NewStyle().Foreground(styles.ColorDanger)
		}
		row := fmt.Sprintf("%-8d %-20s %s %8.1f",
			p.PID,
			truncate(p.Name, 20),
			cpuColor.Render(fmt.Sprintf("%5.1f%%", p.CPUPercent)),
			float64(p.MemRSS)/1024/1024,
		)
		rows = append(rows, row)
	}

	return strings.Join(rows, "\n")
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}
