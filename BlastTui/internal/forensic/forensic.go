// Package forensic enregistre les événements de sécurité dans une timeline
// et permet leur export en JSON ou CSV pour analyse externe.
package forensic

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gocarina/gocsv"
)

// Event représente un événement enregistré dans la timeline forensic.
type Event struct {
	Timestamp time.Time `json:"timestamp" csv:"timestamp"`
	Type      string    `json:"type"      csv:"type"`   // "ALERT" | "SCAN" | "NET" | "INFO"
	Source    string    `json:"source"    csv:"source"` // module émetteur
	Detail    string    `json:"detail"    csv:"detail"`
	Severity  string    `json:"severity"  csv:"severity"`
	PID       int       `json:"pid,omitempty" csv:"pid"`
}

// EventMsg est le message Bubbletea transportant un nouvel événement.
type EventMsg Event

// ExportDoneMsg notifie le TUI que l'export est terminé (contient le chemin du fichier).
type ExportDoneMsg string

// Recorder maintient la timeline en mémoire et gère les exports.
type Recorder struct {
	events  []Event
	logDir  string
}

// NewRecorder crée un Recorder avec le répertoire d'export donné.
func NewRecorder(logDir string) *Recorder {
	return &Recorder{logDir: logDir}
}

// Record ajoute un événement à la timeline.
func (r *Recorder) Record(e Event) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}
	r.events = append(r.events, e)
}

// ExportJSON écrit la timeline complète en JSON dans exports/.
func (r *Recorder) ExportJSON(events []Event) (string, error) {
	if err := os.MkdirAll(r.logDir, 0o750); err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%s/blast_timeline_%s.json",
		r.logDir, time.Now().Format("20060102_150405"))

	data, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		return "", err
	}
	return filename, os.WriteFile(filename, data, 0o640)
}

// ExportCSV écrit la timeline complète en CSV dans exports/.
func (r *Recorder) ExportCSV(events []Event) (string, error) {
	if err := os.MkdirAll(r.logDir, 0o750); err != nil {
		return "", err
	}
	filename := fmt.Sprintf("%s/blast_timeline_%s.csv",
		r.logDir, time.Now().Format("20060102_150405"))

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o640)
	if err != nil {
		return "", err
	}
	defer f.Close()

	return filename, gocsv.MarshalFile(&events, f)
}

// ExportJSONCmd est la commande Bubbletea pour l'export JSON (non-bloquant).
func ExportJSONCmd(events []Event) tea.Cmd {
	return func() tea.Msg {
		r := NewRecorder("exports")
		path, err := r.ExportJSON(events)
		if err != nil {
			return ExportDoneMsg("ERREUR: " + err.Error())
		}
		return ExportDoneMsg(path)
	}
}

// ExportCSVCmd est la commande Bubbletea pour l'export CSV (non-bloquant).
func ExportCSVCmd(events []Event) tea.Cmd {
	return func() tea.Msg {
		r := NewRecorder("exports")
		path, err := r.ExportCSV(events)
		if err != nil {
			return ExportDoneMsg("ERREUR: " + err.Error())
		}
		return ExportDoneMsg(path)
	}
}
