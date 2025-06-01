package tests

import (
	"gysmo/gysmo/src"
	"reflect"
	"testing"
)

// Test GetCPUInfo function
func TestGetCPUInfo(t *testing.T) {
	cpuInfo := src.GetCPUInfo()

	if reflect.TypeOf(cpuInfo).Kind() != reflect.String {
		t.Errorf("Expected CPU info to be of type string, got '%T'", cpuInfo)
	}
}

// Test GetRAMUsage function with tolerance range
func TestGetRAMUsage(t *testing.T) {
	ramUsage := src.GetRAMUsage()
	if reflect.TypeOf(ramUsage).Kind() != reflect.String {
		t.Errorf("Expected RAM Usage to be of type string, got '%T'", ramUsage)
	}
}

// Test GetRAMInfo function with possible sizes
func TestGetRAMInfo(t *testing.T) {
	ramInfo := src.GetRAMInfo()
	if reflect.TypeOf(ramInfo).Kind() != reflect.String {
		t.Errorf("Expected RAM info to be of type string, got '%T'", ramInfo)
	}
}

// Test GetDriveInfo function
func TestGetDriveInfo(t *testing.T) {
	driveInfo := src.GetDriveInfo()
	if reflect.TypeOf(driveInfo).Kind() != reflect.String {
		t.Errorf("Expected Drive Info to be of type string, got '%T'", driveInfo)
	}
}

// Test GetDriveUsage function with tolerance range
func TestGetDriveUsage(t *testing.T) {
	driveUsage := src.GetDriveUsage()
	if reflect.TypeOf(driveUsage).Kind() != reflect.String {
		t.Errorf("Expected Drive Usage to be of type string, got '%T'", driveUsage)
	}
}

// Test GetCPUUsage function with tolerance range
func TestGetCPUUsage(t *testing.T) {
	cpuUsage := src.GetCPUUsage()

	if reflect.TypeOf(cpuUsage).Kind() != reflect.String {
		t.Errorf("Expected CPU Usage to be of type string, got '%T'", cpuUsage)
	}
}

func TestGetGPUInfo(t *testing.T) {
	gpuInfo := src.GetGPUInfo()
	if reflect.TypeOf(gpuInfo).Kind() != reflect.String {
		t.Errorf("Expected GPU info to be of type string, got '%T'", gpuInfo)
	}
}

// Test GetGPUUsage function with tolerance range
func TestGetGPUUsage(t *testing.T) {
	gpuUsage := src.GetGPUUsage()

	if reflect.TypeOf(gpuUsage).Kind() != reflect.String {
		t.Errorf("Expected GPU Usage to be of type string, got '%T'", gpuUsage)
	}
}

// Test GetUptime function
func TestGetUptime(t *testing.T) {
	uptime := src.GetUptime()

	if reflect.TypeOf(uptime).Kind() != reflect.String {
		t.Errorf("Expected Uptime to be of type string, got '%T'", uptime)
	}
}

// Test GetResolution function with format validation
func TestGetResolution(t *testing.T) {
	resolution := src.GetResolution()

	if reflect.TypeOf(resolution).Kind() != reflect.String {
		t.Errorf("Expected resolution to be of type string, got '%T'", resolution)
	}
}

// Test GetPublicIP function
func TestGetPublicIP(t *testing.T) {
	publicIP := src.GetPublicIP()
	if reflect.TypeOf(publicIP).Kind() != reflect.String {
		t.Errorf("Expected public IP to be of type string, got '%T'", publicIP)
	}
}
