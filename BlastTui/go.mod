module github.com/Lemillion24/blast

go 1.22

require (
	// TUI — Charm ecosystem
	github.com/charmbracelet/bubbletea v0.27.1
	github.com/charmbracelet/bubbles v0.20.0
	github.com/charmbracelet/lipgloss v1.0.0

	// CLI flags
	github.com/spf13/cobra v1.8.1

	// Config (YAML/TOML)
	github.com/spf13/viper v1.19.0

	// Réseau — capture de paquets
	github.com/google/gopacket v1.1.19

	// Système — /proc parsing
	github.com/prometheus/procfs v0.15.1
	github.com/shirou/gopsutil/v4 v4.24.11

	// Notifications système (notify-send wrapper)
	github.com/gen2brain/beeep v0.0.0-20240516210008-9c006672e7f4

	// Logging structuré
	github.com/rs/zerolog v1.33.0

	// Export JSON/CSV
	github.com/gocarina/gocsv v0.0.0-20231116093920-b87c2d0e983a

	// Daemon / PID file
	github.com/sevlyar/go-daemon v0.1.6

	// YARA — NOTE: nécessite libyara-dev installée sur le système
	// Installation: sudo apt install libyara-dev  /  sudo pacman -S yara
	// Décommenter quand libyara est disponible :
	// github.com/hillu/go-yara/v4 v4.3.2
)
