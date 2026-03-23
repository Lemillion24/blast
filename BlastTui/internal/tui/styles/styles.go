package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Palette de couleurs BLAST
	ColorPrimary  = lipgloss.Color("#00FF9C") // vert néon
	ColorDanger   = lipgloss.Color("#FF4444") // rouge alerte
	ColorWarning  = lipgloss.Color("#FFB800") // orange avertissement
	ColorMuted    = lipgloss.Color("#555555") // gris discret
	ColorBg       = lipgloss.Color("#0D0D0D") // fond quasi-noir
	ColorBorder   = lipgloss.Color("#1A1A2E") // bleu nuit bordures
	ColorSelected = lipgloss.Color("#16213E") // sélection panneau

	// Bordure standard pour les panneaux
	PanelBorder = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(ColorBorder).
			Padding(0, 1)

	// Panneau actif (sélectionné)
	PanelBorderActive = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(ColorPrimary).
				Padding(0, 1)

	// Titre de panneau
	PanelTitle = lipgloss.NewStyle().
			Foreground(ColorPrimary).
			Bold(true).
			Padding(0, 1)

	// Texte d'alerte critique
	AlertCritical = lipgloss.NewStyle().
			Foreground(ColorDanger).
			Bold(true)

	// Texte d'avertissement
	AlertWarning = lipgloss.NewStyle().
			Foreground(ColorWarning)

	// Texte normal atténué
	TextMuted = lipgloss.NewStyle().
			Foreground(ColorMuted)

	// Barre de statut en bas
	StatusBar = lipgloss.NewStyle().
			Background(ColorBorder).
			Foreground(ColorPrimary).
			Padding(0, 1)

	// Badge pour les valeurs numériques (CPU%, RAM%)
	Badge = lipgloss.NewStyle().
		Foreground(ColorBg).
		Background(ColorPrimary).
		Padding(0, 1).
		Bold(true)

	// Badge danger
	BadgeDanger = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(ColorDanger).
			Padding(0, 1).
			Bold(true)
)
