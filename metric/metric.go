package metric

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type metric struct {
	span *metricSpan

	disableGCStats        bool
	disableHeapStats      bool
	disableGoroutineStats bool
	disableCPUStats       bool
	disableDiskStats      bool
	disableNetStats       bool
}

// metricSpan collects information related to metric plugin per invocation.
type metricSpan struct {
	statTimestamp     int64
	startGCCount      uint32
	endGCCount        uint32
	startPauseTotalNs uint64
	endPauseTotalNs   uint64
	startCPUTimeStat  *cpuTimesStat
	endCPUTimeStat    *cpuTimesStat
	process           *process.Process
	processCpuPercent float64
	systemCpuPercent  float64
	currDiskStat      *process.IOCountersStat
	prevDiskStat      *process.IOCountersStat
	currNetStat       *net.IOCountersStat
	prevNetStat       *net.IOCountersStat
}

func (metric *metric) BeforeExecution(ctx context.Context, request json.RawMessage, wg *sync.WaitGroup) {
	metric.span.statTimestamp = plugin.GetTimestamp()

	if !metric.disableGCStats {
		m := &runtime.MemStats{}
		runtime.ReadMemStats(m)

		metric.span.startGCCount = m.NumGC
		metric.span.startPauseTotalNs = m.PauseTotalNs
	}

	if !metric.disableCPUStats {
		metric.span.startCPUTimeStat = sampleCPUtimesStat()
	}

	wg.Done()
}

func (metric *metric) AfterExecution(ctx context.Context, request json.RawMessage, response interface{}, err interface{}) ([]interface{}, string) {
	mStats := &runtime.MemStats{}
	runtime.ReadMemStats(mStats)

	var stats []interface{}

	if !metric.disableHeapStats {
		h := prepareHeapStatsData(metric, mStats)
		stats = append(stats, h)
	}

	if !metric.disableGCStats {
		metric.span.endGCCount = mStats.NumGC
		metric.span.endPauseTotalNs = mStats.PauseTotalNs

		gc := prepareGCStatsData(metric, mStats)
		stats = append(stats, gc)
	}

	if !metric.disableGoroutineStats {
		g := prepareGoRoutineStatsData(metric)
		stats = append(stats, g)
	}

	if !metric.disableCPUStats {
		metric.span.endCPUTimeStat = sampleCPUtimesStat()

		metric.span.processCpuPercent = getProcessUsagePercent(metric)
		metric.span.systemCpuPercent = getSystemUsagePercent(metric)

		c := prepareCPUStatsData(metric)
		stats = append(stats, c)
	}

	if !metric.disableDiskStats {
		diskStat, err := metric.span.process.IOCounters()
		if err != nil {
			fmt.Println(err)
		} else {
			metric.span.currDiskStat = diskStat
			d := prepareDiskStatsData(metric)
			stats = append(stats, d)
		}
	}

	if !metric.disableNetStats {
		netIOStat, err := net.IOCounters(false)
		if err != nil {
			fmt.Println(err)
		} else {
			metric.span.currNetStat = &netIOStat[all]
			n := prepareNetStatsData(metric)
			stats = append(stats, n)
		}
	}

	return stats, statDataType
}

//OnPanic just collect the metrics and send them as in the AfterExecution
func (metric *metric) OnPanic(ctx context.Context, request json.RawMessage, err interface{}, stackTrace []byte) ([]interface{}, string) {
	return metric.AfterExecution(ctx, request, nil, err)
}
