package pkg

import (
	"fmt"
	"sort"
	"time"
)

// This file was only for debugging and testing optimization.
// I keep it just in case

// Function type for functions that return a string
type StringFunc func() string

// Function type for functions that return two uint64 values
type Uint64Func func() (uint64, uint64)

// MeasureTime is a higher-order function that measures the execution time of a function
func MeasureTime(name string, f StringFunc) (string, time.Duration) {
	start := time.Now()
	result := f()
	duration := time.Since(start)
	return result, duration
}

// MeasureTimeUint64 is a higher-order function that measures the execution time of a function returning two uint64 values
func MeasureTimeUint64(name string, f Uint64Func) (uint64, uint64, time.Duration) {
	start := time.Now()
	idle, total := f()
	duration := time.Since(start)
	return idle, total, duration
}

func MeasureMenuItems() string {
	config := Config{
		Items: []ConfigItem{
			{Keyword: "os_pretty_name"},
			{Keyword: "hostname"},
			{Keyword: "kernel"},
			{Keyword: "shell"},
			{Keyword: "uptime"},
			{Keyword: "cpu"},
			{Keyword: "ram"},
			{Keyword: "drive"},
			{Keyword: "gpu %"},
			{Keyword: "cpu %"},
			{Keyword: "ram %"},
			{Keyword: "drive %"},
			{Keyword: "term"},
			{Keyword: "processes"},
			{Keyword: "wm"},
			{Keyword: "ip"},
			{Keyword: "public ip"},
			{Keyword: "resolution"},
		},
	}
	useDataFile := false // or true, depending on what you want to test
	items := MenuItems(config, useDataFile)
	return fmt.Sprintf("%v", items)
}

// Function to measure and rank the execution time of all functions
func MeasureAllFunctions() {
	type FunctionResult struct {
		Name     string
		Duration time.Duration
	}

	var results []FunctionResult

	// Measure each function and store the results
	_, duration := MeasureTime("GetCPUInfo", GetCPUInfo)
	results = append(results, FunctionResult{"GetCPUInfo", duration})

	_, duration = MeasureTime("GetRAMUsage", GetRAMUsage)
	results = append(results, FunctionResult{"GetRAMUsage", duration})

	_, duration = MeasureTime("GetRAMInfo", GetRAMInfo)
	results = append(results, FunctionResult{"GetRAMInfo", duration})

	_, _, duration = MeasureTimeUint64("GetCPUSample", GetCPUSample)
	results = append(results, FunctionResult{"GetCPUSample", duration})

	_, duration = MeasureTime("GetShell", GetShell)
	results = append(results, FunctionResult{"GetShell", duration})

	_, duration = MeasureTime("GetTerminal", GetTerminal)
	results = append(results, FunctionResult{"GetTerminal", duration})

	_, duration = MeasureTime("GetDriveInfo", GetDriveInfo)
	results = append(results, FunctionResult{"GetDriveInfo", duration})

	_, duration = MeasureTime("GetDriveUsage", GetDriveUsage)
	results = append(results, FunctionResult{"GetDriveUsage", duration})

	_, duration = MeasureTime("GetGPUInfo", GetGPUInfo)
	results = append(results, FunctionResult{"GetGPUInfo", duration})

	_, duration = MeasureTime("GetGPUUsage", GetGPUUsage)
	results = append(results, FunctionResult{"GetGPUUsage", duration})

	_, duration = MeasureTime("GetNvidiaGPUInfo", GetNvidiaGPUInfo)
	results = append(results, FunctionResult{"GetNvidiaGPUInfo", duration})

	_, duration = MeasureTime("GetNvidiaGPUUsage", GetNvidiaGPUUsage)
	results = append(results, FunctionResult{"GetNvidiaGPUUsage", duration})

	_, duration = MeasureTime("GetAmdGPUInfo", GetAmdGPUInfo)
	results = append(results, FunctionResult{"GetAmdGPUInfo", duration})

	_, duration = MeasureTime("GetAmdGPUUsage", GetAmdGPUUsage)
	results = append(results, FunctionResult{"GetAmdGPUUsage", duration})

	_, duration = MeasureTime("GetIntelGPUInfo", GetIntelGPUInfo)
	results = append(results, FunctionResult{"GetIntelGPUInfo", duration})

	_, duration = MeasureTime("GetIntelGPUUsage", GetIntelGPUUsage)
	results = append(results, FunctionResult{"GetIntelGPUUsage", duration})

	_, duration = MeasureTime("GetUptime", GetUptime)
	results = append(results, FunctionResult{"GetUptime", duration})

	_, duration = MeasureTime("GetWM", GetWM)
	results = append(results, FunctionResult{"GetWM", duration})

	_, duration = MeasureTime("GetResolution", GetResolution)
	results = append(results, FunctionResult{"GetResolution", duration})

	_, duration = MeasureTime("GetIP", GetIP)
	results = append(results, FunctionResult{"GetIP", duration})

	_, duration = MeasureTime("GetPublicIP", GetPublicIP)
	results = append(results, FunctionResult{"GetPublicIP", duration})

	// measure menuitems function
	_, duration = MeasureTime("menuitems", MeasureMenuItems)
	results = append(results, FunctionResult{"menuitems", duration})

	// Sort the results by duration
	sort.Slice(results, func(i, j int) bool {
		return results[i].Duration > results[j].Duration
	})

	// Print the sorted results
	for _, result := range results {
		fmt.Printf("%s took %v\n", result.Name, result.Duration)
	}
}
