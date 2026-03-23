// Package monitor collecte les métriques système via /proc et gopsutil.
// Il abstrait les différences Linux/Windows pour permettre le portage futur.
package monitor

import (
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/process"
)

// SystemStats regroupe toutes les métriques collectées en un instant T.
type SystemStats struct {
	CPUPercent   float64
	MemUsed      uint64 // bytes
	MemTotal     uint64 // bytes
	DiskUsed     uint64 // bytes
	DiskTotal    uint64 // bytes
	TopProcesses []ProcessInfo
}

// ProcessInfo décrit un processus en cours d'exécution.
type ProcessInfo struct {
	PID        int32
	Name       string
	CPUPercent float64
	MemRSS     uint64 // bytes
	Username   string
	Status     string
}

// Collect effectue une collecte complète des métriques système.
// C'est la fonction centrale appelée à chaque tick du TUI.
func Collect() (SystemStats, error) {
	stats := SystemStats{}

	// CPU global (moyenne sur tous les cœurs)
	cpuPcts, err := cpu.Percent(0, false)
	if err == nil && len(cpuPcts) > 0 {
		stats.CPUPercent = cpuPcts[0]
	}

	// Mémoire RAM
	vmStat, err := mem.VirtualMemory()
	if err == nil {
		stats.MemUsed = vmStat.Used
		stats.MemTotal = vmStat.Total
	}

	// Disque (partition racine)
	diskStat, err := disk.Usage("/")
	if err == nil {
		stats.DiskUsed = diskStat.Used
		stats.DiskTotal = diskStat.Total
	}

	// Top processus par CPU
	stats.TopProcesses, _ = topProcesses(15)

	return stats, nil
}

// topProcesses retourne les N processus les plus gourmands en CPU.
func topProcesses(n int) ([]ProcessInfo, error) {
	procs, err := process.Processes()
	if err != nil {
		return nil, err
	}

	var infos []ProcessInfo
	for _, p := range procs {
		name, _ := p.Name()
		cpuPct, _ := p.CPUPercent()
		memInfo, _ := p.MemoryInfo()
		user, _ := p.Username()
		status, _ := p.Status()

		rss := uint64(0)
		if memInfo != nil {
			rss = memInfo.RSS
		}

		infos = append(infos, ProcessInfo{
			PID:        p.Pid,
			Name:       name,
			CPUPercent: cpuPct,
			MemRSS:     rss,
			Username:   user,
			Status:     status[0],
		})
	}

	// Tri par CPU décroissant (bubble sort simple pour lisibilité)
	for i := 0; i < len(infos)-1; i++ {
		for j := i + 1; j < len(infos); j++ {
			if infos[j].CPUPercent > infos[i].CPUPercent {
				infos[i], infos[j] = infos[j], infos[i]
			}
		}
	}

	if n > len(infos) {
		n = len(infos)
	}
	return infos[:n], nil
}

// KillProcess envoie SIGKILL au processus identifié par son PID.
// Toujours appelé après confirmation explicite de l'utilisateur dans le TUI.
func KillProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return p.Kill()
}

// SuspendProcess envoie SIGSTOP pour mettre en pause (quarantaine) un processus.
func SuspendProcess(pid int32) error {
	p, err := process.NewProcess(pid)
	if err != nil {
		return err
	}
	return p.Suspend()
}
