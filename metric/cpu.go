package metric

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/process"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type cpuStatsData struct {
	Id                 string `json:"id"`
	TransactionId      string `json:"transactionId"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTimestamp      int64  `json:"statTimestamp"`

	// ProcessCPUPercent is the process usage of the total CPU time
	ProcessCPUPercent float64 `json:"procPercent"`

	// SystemCPUPercent is the system usage of the total CPU time
	SystemCPUPercent float64 `json:"sysPercent"`
}

func prepareCPUStatsData(metric *metric) cpuStatsData {
	return cpuStatsData{
		Id:                 plugin.GenerateNewId(),
		TransactionId:      plugin.TransactionId,
		ApplicationName:    plugin.ApplicationName,
		ApplicationId:      plugin.ApplicationId,
		ApplicationVersion: plugin.ApplicationVersion,
		ApplicationProfile: plugin.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           cpuStat,
		StatTimestamp:      metric.statTimestamp,
		ProcessCPUPercent:  metric.processCpuPercent,
		SystemCPUPercent:   metric.systemCpuPercent,
	}
}

func getCPUUsagePercentage(p *process.Process) (float64, float64, error) {
	sysUsage, err := cpu.Percent(0, false)
	if err != nil {
		return 0, 0, err
	}

	processUsage, err := p.Percent(0)
	if err != nil {
		return 0, 0, err
	}

	return processUsage, sysUsage[0], nil
}
