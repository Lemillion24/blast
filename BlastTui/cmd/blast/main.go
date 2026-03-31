package main

import (
	"fmt"
	"os"

	"github.com/Lemillion24/blast/internal/tui"
	"github.com/spf13/cobra"
)

var (
	daemonMode bool
	configPath string
	logLevel   string
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "blast",
	Short: "BLAST — Security monitoring, audit & forensic TUI",
	Long: `BLAST est un outil de monitoring système orienté sécurité.
Il combine surveillance temps réel, analyse YARA, capture réseau
et forensic dans une interface TUI (Terminal User Interface).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if daemonMode {
			return runDaemon()
		}
		return runTUI()
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&daemonMode, "daemon", "d", false, "Lancer BLAST en mode daemon (arrière-plan)")
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config/blast.yaml", "Chemin vers le fichier de configuration")
	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "l", "info", "Niveau de log (debug, info, warn, error)")

	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(scanCmd)
	rootCmd.AddCommand(stopCmd)
}

// runTUI lance l'interface TUI interactive.
func runTUI() error {
	return tui.Start()
}

// runDaemon lance BLAST en mode daemon (service arrière-plan).
func runDaemon() error {
	fmt.Println("Lancement de BLAST en mode daemon...")
	// TODO Phase 4 : appeler internal/daemon.Start()
	return nil
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Affiche la version de BLAST",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("BLAST v0.1.0-alpha")
	},
}

var scanCmd = &cobra.Command{
	Use:   "scan [chemin]",
	Short: "Lancer un scan YARA sur un fichier ou répertoire",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("Scan YARA de : %s\n", args[0])
		// TODO Phase 3 : appeler internal/security.ScanPath(args[0])
		return nil
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Arrêter le daemon BLAST en cours d'exécution",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Arrêt du daemon BLAST...")
		// TODO Phase 4 : envoyer signal au daemon via PID file
		return nil
	},
}
