// Package security gère la détection de menaces :
// - Scan YARA des fichiers (via go-yara/CGO)
// - Règles comportementales custom (YAML)
// - Alertes et scoring de risque
package security

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v3"
)

// Alert représente une menace détectée.
type Alert struct {
	Timestamp time.Time
	Severity  string // "critical" | "high" | "medium" | "low"
	RuleName  string
	Target    string // fichier ou PID concerné
	Detail    string
	Source    string // "yara" | "behavior" | "network"
}

// Rule représente une règle de détection chargée.
type Rule struct {
	Name        string
	Source      string // "yara" | "yaml"
	Description string
	Severity    string
	Enabled     bool
}

// BehaviorRule est une règle comportementale définie en YAML.
// Elle décrit des patterns suspects (ex: process qui écoute sur un port > 1024
// et envoie plus de 10 Mo en moins d'une minute).
type BehaviorRule struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description"`
	Severity    string            `yaml:"severity"`
	Conditions  map[string]string `yaml:"conditions"`
	Action      string            `yaml:"action"` // "alert" | "kill" | "suspend"
}

// Scanner orchestre les différents moteurs de détection.
type Scanner struct {
	rules         []Rule
	behaviorRules []BehaviorRule
	rulesDir      string
	// yaraScanner  *yara.Scanner  // décommenter quand CGO yara disponible
}

// NewScanner crée un scanner avec le répertoire de règles donné.
func NewScanner(rulesDir string) *Scanner {
	return &Scanner{rulesDir: rulesDir}
}

// LoadRules charge les règles YARA et YAML depuis le répertoire configuré.
func (s *Scanner) LoadRules() error {
	s.rules = nil

	// Charger les règles comportementales YAML
	yamlDir := filepath.Join(s.rulesDir, "custom")
	entries, err := os.ReadDir(yamlDir)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("lecture répertoire règles: %w", err)
	}

	for _, entry := range entries {
		if filepath.Ext(entry.Name()) != ".yaml" && filepath.Ext(entry.Name()) != ".yml" {
			continue
		}
		path := filepath.Join(yamlDir, entry.Name())
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var br BehaviorRule
		if err := yaml.Unmarshal(data, &br); err != nil {
			continue
		}
		s.behaviorRules = append(s.behaviorRules, br)
		s.rules = append(s.rules, Rule{
			Name:        br.Name,
			Source:      "yaml",
			Description: br.Description,
			Severity:    br.Severity,
			Enabled:     true,
		})
	}

	// TODO Phase 3 : charger les règles YARA via go-yara
	// yaraDir := filepath.Join(s.rulesDir, "yara")
	// compiler, _ := yara.NewCompiler()
	// ... compiler.AddFile(f) pour chaque .yar ...
	// s.yaraScanner, _ = yara.NewScanner(compiler.GetRules())

	return nil
}

// ScanFile analyse un fichier avec toutes les règles disponibles.
func (s *Scanner) ScanFile(path string) ([]Alert, error) {
	var alerts []Alert

	// TODO Phase 3 : scan YARA réel
	// matches, err := s.yaraScanner.ScanFile(path)
	// for _, m := range matches { alerts = append(alerts, yaraMatchToAlert(m, path)) }

	// Stub : vérifier l'extension comme exemple comportemental
	if filepath.Ext(path) == ".sh" {
		info, err := os.Stat(path)
		if err == nil && info.Mode()&0o111 != 0 {
			alerts = append(alerts, Alert{
				Timestamp: time.Now(),
				Severity:  "medium",
				RuleName:  "executable_script",
				Target:    path,
				Detail:    "Script shell exécutable détecté",
				Source:    "behavior",
			})
		}
	}

	return alerts, nil
}

// RulesLoadedMsg transporte les règles chargées depuis le disque.
type RulesLoadedMsg []Rule

// ScanResultMsg transporte les alertes issues d'un scan.
type ScanResultMsg []Alert

// LoadRulesCmd est la commande Bubbletea pour charger les règles.
func LoadRulesCmd() tea.Cmd {
	return func() tea.Msg {
		scanner := NewScanner("rules")
		if err := scanner.LoadRules(); err != nil {
			return RulesLoadedMsg(nil)
		}
		return RulesLoadedMsg(scanner.rules)
	}
}

// QuickScanCmd lance un scan rapide sur /tmp et /dev/shm (zones à risque).
func QuickScanCmd() tea.Cmd {
	return func() tea.Msg {
		scanner := NewScanner("rules")
		_ = scanner.LoadRules()

		var allAlerts []Alert
		targets := []string{"/tmp", "/dev/shm", "/var/tmp"}

		for _, target := range targets {
			_ = filepath.Walk(target, func(path string, info os.FileInfo, err error) error {
				if err != nil || info.IsDir() {
					return nil
				}
				alerts, _ := scanner.ScanFile(path)
				allAlerts = append(allAlerts, alerts...)
				return nil
			})
		}

		return ScanResultMsg(allAlerts)
	}
}
