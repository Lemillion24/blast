# BLAST — Security Monitoring TUI

> **B**ehavior & **L**ive **A**udit **S**ecurity **T**ool

Outil de monitoring système orienté sécurité, forensic et surveillance réseau.
Interface TUI (Terminal User Interface) construite avec [Bubbletea](https://github.com/charmbracelet/bubbletea).

```
┌─────────────────────────────────────────────────────────────┐
│  Monitoring    Réseau     Sécurité    Forensic    Logs      │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  CPU  [████████░░░░░░░░░░░░░░░░░░░░] 27.3/100 %           │
│  RAM  [████████████████░░░░░░░░░░░░] 8.2/16.0 Go           │
│  Disk [████████░░░░░░░░░░░░░░░░░░░░] 120.5/512.0 Go        │
│                                                             │
│  PID     NOM                  CPU%     MEM(Mo)             │
│  1234    firefox              12.4%    412.0               │
│  5678    code                  8.1%    280.3               │
│  ...                                                        │
│                                                             │
├─────────────────────────────────────────────────────────────┤
│  BLAST v0.1.0  ⚠ 2 alerte(s)    [Tab] Changer | [q] Quitter│
└─────────────────────────────────────────────────────────────┘
```

---

## Prérequis

```bash
# Arch Linux
sudo pacman -S go yara libpcap

# Ubuntu/Debian
sudo apt install golang-go libyara-dev libpcap-dev

# Vérifier
go version   # >= 1.22
```

## Installation & lancement

```bash
git clone https://github.com/Lemillion24/blast.git
cd blast/BlastTui
go mod tidy
go run ./cmd/blast/

# Avec privilèges réseau (lecture /proc complète)
sudo go run ./cmd/blast/

# Mode daemon
go run ./cmd/blast/ --daemon

# Scan YARA d'un répertoire
go run ./cmd/blast/ scan /home/user/Downloads

# Voir la version
go run ./cmd/blast/ version

# Compiler un binaire
go build -o blast ./cmd/blast/
./blast
```

---

## Architecture

```
BlastTui/
├── cmd/blast/            # Point d'entrée + CLI (cobra)
│   └── main.go           #   → lance tui.Start() ou le daemon
├── internal/
│   ├── tui/              # Interface utilisateur (bubbletea)
│   │   ├── model.go      #   AppModel racine, 5 onglets, dispatch TickMsg → RefreshMsg
│   │   ├── components/   #   Un fichier par panneau (chacun implémente tea.Model)
│   │   │   ├── monitor.go
│   │   │   ├── network.go
│   │   │   ├── security.go
│   │   │   ├── forensic.go
│   │   │   └── logs.go
│   │   └── styles/       #   Palette de couleurs et styles lipgloss
│   ├── monitor/          # Métriques système (/proc, gopsutil)
│   ├── network/          # Connexions réseau PID↔socket↔DNS via /proc/net/tcp
│   ├── security/         # YARA + règles comportementales YAML
│   ├── forensic/         # Timeline + export JSON/CSV
│   └── alerts/           # Notifications multi-canaux (TUI + notify-send + log)
├── config/
│   └── blast.yaml        # Configuration (rechargement à chaud prévu Phase 4)
├── rules/
│   ├── yara/             # Règles YARA (.yar)
│   └── custom/           # Règles comportementales (.yaml)
├── logs/                 # Logs BLAST et alertes
└── exports/              # Exports forensic JSON/CSV
```

---

## Plan de développement par phases

### Phase 1 — Fondations (semaines 1-2)
**Objectif :** TUI qui tourne, métriques système affichées

- [x] Structure du projet et `go.mod` (module `github.com/Lemillion24/blast`)
- [x] CLI avec `cobra` (flags `--daemon`, `--config`, `scan`, `version`, `stop`)
- [x] Modèle TUI principal Bubbletea (`AppModel`, 5 onglets)
- [x] Styles Lipgloss (palette BLAST vert néon / rouge / orange)
- [x] Module `monitor` : CPU, RAM, Disk, top processus (gopsutil)
- [x] Connecter `monitor` au panneau TUI via `RefreshMsg` (rafraîchissement chaque seconde)
- [x] Connecter `tui.Start()` dans `main.go` (le TUI se lance au démarrage)
- [x] Tous les panneaux implémentent `tea.Model` (`Init`, `Update`, `View`)
- [x] Types de messages typés : `StatsMsg`, `ConnectionsMsg`, `RulesLoadedMsg`, `ScanResultMsg`
- [x] Protection panic sur `p.Status()` (slice vide)
- [x] `go build ./...` + `go vet ./...` passent sans erreur

**Compétences Go acquises :** `tea.Model`, channels, goroutines basics

---

### Phase 2 — Réseau sans capture (semaines 3-4)
**Objectif :** Voir quel processus parle à qui, sans root

- [x] Parser `/proc/net/tcp` et `/proc/net/tcp6`
- [x] Mapper inode → PID via `/proc/[pid]/fd` (buildInodePIDMap)
- [x] Résolution DNS inverse (net.LookupAddr)
- [x] Panneau réseau TUI avec rafraîchissement automatique
- [ ] Compléter `extractInode()` (stub actuel, corrélation PID↔socket incomplète)
- [ ] Cache DNS avec TTL configurable
- [ ] Filtre texte libre dans le panneau réseau
- [ ] Alertes si processus système (ex: `systemd`) ouvre connexion externe

**Compétences Go acquises :** parsing de fichiers, maps, net.Lookup

---

### Phase 3 — Sécurité YARA + règles YAML (semaines 5-7)
**Objectif :** Détection de malware et comportements suspects

- [x] Panneau sécurité TUI avec affichage règles et alertes
- [x] Chargement des règles YAML comportementales (`LoadRulesCmd` → `RulesLoadedMsg`)
- [x] Scan rapide des zones à risque (`QuickScanCmd` → `ScanResultMsg`)
- [x] Scan stub comportemental (détection scripts .sh exécutables)
- [ ] Installer `libyara-dev` et décommenter `go-yara` dans `go.mod`
- [ ] Intégrer `go-yara` dans `internal/security` (scan réel)
- [ ] Moteur de règles YAML : parser les conditions et les évaluer
- [ ] Corrélation réseau + process (ex: bash + connexion sortante = alerte)
- [ ] Kill/Suspend processus avec confirmation TUI
- [ ] Ajouter règles YARA communautaires (Malware Bazaar)

**Compétences Go acquises :** CGO, filepath.Walk, goroutines de scan

---

### Phase 4 — Capture réseau profonde + Daemon (semaines 8-10)
**Objectif :** Voir le payload réseau, fonctionner en service

- [ ] Intégrer `gopacket` + `libpcap` pour capture réelle
- [ ] Dégrader gracieusement si pas de `CAP_NET_RAW`
- [ ] DPI basique : détecter HTTP, DNS, TLS (sans déchiffrer)
- [ ] Alertes sur trafic anormal (volume, destination blacklistée)
- [ ] Mode daemon avec PID file et rechargement config à chaud
- [ ] Commande `blast stop` pour arrêter le daemon

**Compétences Go acquises :** syscalls, signals, CGO avancé

---

### Phase 5 — Forensic & polish (semaines 11-12)
**Objectif :** Export, documentation, présentation portfolio

- [x] Panneau forensic TUI avec timeline d'événements
- [x] Export JSON et CSV horodatés (`ExportJSONCmd`, `ExportCSVCmd`)
- [x] Panneau logs temps réel avec viewport scrollable (bubbles/viewport)
- [x] Système d'alertes multi-canaux (TUI + notify-send + log fichier)
- [ ] Rapport HTML optionnel (avec template Go)
- [ ] Support Windows (abstraction `pkg/sysinfo`)
- [ ] Tests unitaires (monitor, network parser, rules engine)
- [ ] README complet + captures d'écran

---

## Raccourcis clavier

| Touche | Action |
|--------|--------|
| `Tab` / `Shift+Tab` | Changer d'onglet |
| `1` à `5` | Aller directement à l'onglet |
| `s` | Scan YARA rapide (onglet Sécurité) |
| `e` | Export JSON (onglet Forensic) |
| `E` | Export CSV (onglet Forensic) |
| `k` | Kill processus sélectionné (confirmation requise) |
| `q` / `Ctrl+C` | Quitter |

---

## Dépendances clés

| Bibliothèque | Usage | CGO ? |
|---|---|---|
| `charmbracelet/bubbletea` | Framework TUI | Non |
| `charmbracelet/lipgloss` | Styles terminaux | Non |
| `charmbracelet/bubbles` | Composants (viewport, list) | Non |
| `spf13/cobra` | CLI | Non |
| `shirou/gopsutil` | Métriques système cross-platform | Non |
| `google/gopacket` | Capture réseau | **Oui** (libpcap) |
| `hillu/go-yara` | Scan YARA | **Oui** (libyara) |
| `gen2brain/beeep` | Notifications bureau | Non |
| `rs/zerolog` | Logging structuré | Non |
| `gocarina/gocsv` | Export CSV | Non |
