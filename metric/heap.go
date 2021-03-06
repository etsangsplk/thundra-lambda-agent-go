package metric

import (
	"fmt"
	"runtime"

	"github.com/thundra-io/thundra-lambda-agent-go/plugin"
)

type heapStatsData struct {
	Id                 string `json:"id"`
	TransactionId      string `json:"transactionId"`
	ApplicationName    string `json:"applicationName"`
	ApplicationId      string `json:"applicationId"`
	ApplicationVersion string `json:"applicationVersion"`
	ApplicationProfile string `json:"applicationProfile"`
	ApplicationType    string `json:"applicationType"`
	StatName           string `json:"statName"`
	StatTimestamp      int64  `json:"statTimestamp"`

	// HeapAlloc is bytes of allocated heap objects.
	//
	// "Allocated" heap objects include all reachable objects, as
	// well as unreachable objects that the garbage collector has
	// not yet freed.
	HeapAlloc uint64 `json:"heapAlloc"`

	// HeapSys estimates the largest size the heap has had.
	HeapSys uint64 `json:"heapSys"`

	// HeapInuse is bytes in in-use spans.

	// In-use spans have at least one object in them. These spans
	// can only be used for other objects of roughly the same
	// size.
	HeapInuse uint64 `json:"heapInuse"`

	// HeapObjects is the number of allocated heap objects.
	HeapObjects uint64 `json:"heapObjects"`

	// MemoryPercent returns how many percent of the total RAM this process uses
	MemoryPercent float32 `json:"memoryPercent"`
}

func prepareHeapStatsData(metric *metric, memStats *runtime.MemStats) heapStatsData {
	mp, err := proc.MemoryPercent()
	if err != nil {
		fmt.Println(err)
	}

	return heapStatsData{
		Id:                 plugin.GenerateNewId(),
		TransactionId:      plugin.TransactionId,
		ApplicationName:    plugin.ApplicationName,
		ApplicationId:      plugin.ApplicationId,
		ApplicationVersion: plugin.ApplicationVersion,
		ApplicationProfile: plugin.ApplicationProfile,
		ApplicationType:    plugin.ApplicationType,
		StatName:           heapStat,
		StatTimestamp:      metric.span.statTimestamp,

		HeapAlloc:     memStats.HeapAlloc,
		HeapSys:       memStats.HeapSys,
		HeapInuse:     memStats.HeapInuse,
		HeapObjects:   memStats.HeapObjects,
		MemoryPercent: mp,
	}
}
