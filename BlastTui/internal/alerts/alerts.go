// Package alerts centralise l'envoi de notifications multi-canaux :
// panneau TUI, notify-send (bureau), et fichier log persistant.
package alerts

import (
	"fmt"
	"os"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/rs/zerolog"
)

// Level représente la criticité d'une alerte.
type Level int

const (
	LevelInfo Level = iota
	LevelWarning
	LevelHigh
	LevelCritical
)

func (l Level) String() string {
	switch l {
	case LevelWarning:
		return "WARN"
	case LevelHigh:
		return "HIGH"
	case LevelCritical:
		return "CRITICAL"
	default:
		return "INFO"
	}
}

// Manager orchestre l'envoi d'alertes vers tous les canaux configurés.
type Manager struct {
	logger     zerolog.Logger
	logFile    *os.File
	tuiChannel chan<- Notification // channel vers le TUI
}

// Notification est envoyée au TUI via channel.
type Notification struct {
	Level     Level
	Title     string
	Message   string
	Timestamp time.Time
}

// NewManager crée un Manager d'alertes avec fichier log et channel TUI.
func NewManager(logPath string, tuiChan chan<- Notification) (*Manager, error) {
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o640)
	if err != nil {
		return nil, fmt.Errorf("ouverture log alertes: %w", err)
	}

	logger := zerolog.New(f).With().Timestamp().Logger()

	return &Manager{
		logger:     logger,
		logFile:    f,
		tuiChannel: tuiChan,
	}, nil
}

// Send envoie une alerte sur tous les canaux actifs.
func (m *Manager) Send(level Level, title, message string) {
	notif := Notification{
		Level:     level,
		Title:     title,
		Message:   message,
		Timestamp: time.Now(),
	}

	// 1. Log fichier (toujours)
	logEvent := m.logger.With().
		Str("level", level.String()).
		Str("title", title).
		Logger()

	switch level {
	case LevelCritical:
		logEvent.Error().Msg(message)
	case LevelHigh:
		logEvent.Warn().Msg(message)
	default:
		logEvent.Info().Msg(message)
	}

	// 2. Notification bureau (notify-send via beeep)
	if level >= LevelHigh {
		icon := ""
		if level == LevelCritical {
			icon = "dialog-error"
		} else {
			icon = "dialog-warning"
		}
		// Non-bloquant : on ignore l'erreur si le bureau n'est pas disponible
		_ = beeep.Notify("BLAST — "+title, message, icon)
	}

	// 3. Canal TUI (non-bloquant pour éviter de bloquer les goroutines)
	if m.tuiChannel != nil {
		select {
		case m.tuiChannel <- notif:
		default:
			// Canal plein, on drop plutôt que bloquer
		}
	}
}

// Close libère les ressources du Manager.
func (m *Manager) Close() error {
	if m.logFile != nil {
		return m.logFile.Close()
	}
	return nil
}
