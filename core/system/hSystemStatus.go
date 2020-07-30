package system

import (
	"net/http"
	"time"

	"github.com/seknox/trasa/utils"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
)

type sysStatus struct {
	HostStat *host.InfoStat         `json:"hostStatus"`
	MemStat  *mem.VirtualMemoryStat `json:"memStatus"`
	DiskStat *disk.UsageStat        `json:"diskStatus"`
	CPUStat  cpustat                `json:"cpuStat"`
}

type cpustat struct {
	CpuCount int       `json:"cpuCount"`
	CPUStat  []float64 `json:"cpuStat"`
}

func SystemStatus(w http.ResponseWriter, r *http.Request) {

	var systemStat sysStatus
	systemStat.HostStat = hostStatus()
	systemStat.MemStat = memStatus()
	systemStat.DiskStat = diskStatus()
	count, perc := cpuStatus()
	systemStat.CPUStat.CpuCount = count
	systemStat.CPUStat.CPUStat = perc

	utils.TrasaResponse(w, 200, "success", "status fetched", "SystemStatus", systemStat)

}

func diskStatus() *disk.UsageStat {
	diskStat, err := disk.Usage("/")
	if err != nil {
		logrus.Error(err)
	}

	return diskStat
}

func memStatus() *mem.VirtualMemoryStat {
	vmem, err := mem.VirtualMemory()
	if err != nil {
		logrus.Error(err)
	}
	return vmem
}

func hostStatus() *host.InfoStat {
	hostd, err := host.Info()
	if err != nil {
		logrus.Error(err)
	}
	return hostd
}

func cpuStatus() (int, []float64) {
	count, err := cpu.Counts(true)
	if err != nil {
		logrus.Error(err)
	}
	cpus, err := cpu.Percent(time.Second*1, false)
	if err != nil {
		logrus.Error(err)
	}
	return count, cpus
}
