package tests

import (
	"bytes"
	"gysmo/src/pkg"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"testing"
)

func mockExecCommand(command string, args ...string) *exec.Cmd {
	cmd := exec.Command("echo")
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}
	cmd.Stdout.(*bytes.Buffer).WriteString(`Screen 0: minimum 8 x 8, current 2560 x 1440, maximum 32767 x 32767
HDMI-1 connected primary 2560x1440+0+0 (normal left inverted right x axis y axis) 597mm x 336mm
   2560x1440     59.95*+
   1920x1080     60.00    59.94
   1680x1050     59.88
   1280x1024     75.02    60.02
   1440x900      59.90
   1280x800      59.91
   1280x720      60.00    59.94
   1024x768      75.03    70.07    60.00
   800x600       75.00    60.32
   720x480       60.00    59.94
   640x480       75.00    72.81    66.67    60.00    59.94
`)
	return cmd
}

func init() {
	// Override the actual functions with the mock functions
	pkg.ExecCommand = mockExecCommand
}

// Test GetCPUInfo function
func TestGetCPUInfo(t *testing.T) {
	mockData := `processor   : 0
vendor_id   : GenuineIntel
cpu family  : 6
model       : 158
model name  : Intel(R) Core(TM) i7-8550U CPU @ 1.80GHz
stepping    : 10
microcode   : 0xea
cpu MHz     : 1992.000
cache size  : 8192 KB
physical id : 0
siblings    : 8
core id     : 0
cpu cores   : 4
apicid      : 0
initial apicid  : 0
fpu     : yes
fpu_exception   : yes
cpuid level : 22
wp      : yes
flags       : fpu vme de pse tsc msr pae mce cx8 apic sep mtrr pge mca cmov pat pse36 clflush dts acpi mmx fxsr sse sse2 ss ht tm pbe syscall nx pdpe1gb rdtscp lm constant_tsc art arch_perfmon pebs bts rep_good nopl xtopology nonstop_tsc cpuid aperfmperf pni pclmulqdq dtes64 monitor ds_cpl vmx smx est tm2 ssse3 sdbg fma cx16 xtpr pdcm pcid sse4_1 sse4_2 x2apic movbe popcnt tsc_deadline_timer aes xsave avx f16c rdrand lahf_lm abm 3dnowprefetch cpuid_fault epb invpcid_single pti ssbd ibrs ibpb stibp tpr_shadow vnmi flexpriority ept vpid ept_ad fsgsbase tsc_adjust bmi1 avx2 smep bmi2 erms invpcid mpx rdseed adx smap clflushopt intel_pt xsaveopt xsavec xgetbv1 xsaves dtherm ida arat pln pts hwp hwp_notify hwp_act_window hwp_epp md_clear flush_l1d
bugs        : spectre_v1 spectre_v2 spec_store_bypass l1tf mds swapgs itlb_multihit
bogomips    : 3984.00
clflush size    : 64
cache_alignment : 64
address sizes   : 39 bits physical, 48 bits virtual
power management:`
	pkg.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", mockData)
	}
	cpuInfo := pkg.GetCPUInfo()

	valid := false
	possibleVendors := []string{"Intel", "AMD"}
	for _, vendors := range possibleVendors {
		if strings.Contains(cpuInfo, vendors) {
			valid = true
			break
		}
	}
	if !valid {
		t.Errorf("Expected RAM info to be one of %v, got %s GB", possibleVendors, cpuInfo)
	}
}

// Test GetRAMUsage function with tolerance range
func TestGetRAMUsage(t *testing.T) {
	ramUsage := pkg.GetRAMUsage()
	expectedMin := 10.0 // Minimum expected value
	expectedMax := 90.0 // Maximum expected value

	ramUsageValue, err := strconv.ParseFloat(strings.TrimSuffix(ramUsage, "%"), 64)
	if err != nil {
		t.Fatalf("Failed to parse RAM usage: %v", err)
	}

	if ramUsageValue < expectedMin || ramUsageValue > expectedMax {
		t.Errorf("Expected RAM usage to be between %.2f%% and %.2f%%, got %.2f%%", expectedMin, expectedMax, ramUsageValue)
	}
}

// Test GetRAMInfo function with possible sizes
func TestGetRAMInfo(t *testing.T) {
	ramInfo := pkg.GetRAMInfo()
	possibleSizes := []uint64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536}

	ramInfoValue, err := strconv.ParseUint(strings.TrimSuffix(ramInfo, " GB"), 10, 64)
	if err != nil {
		t.Fatalf("Failed to parse RAM info: %v", err)
	}

	valid := false
	for _, size := range possibleSizes {
		if ramInfoValue == size {
			valid = true
			break
		}
	}

	if !valid {
		t.Errorf("Expected RAM info to be one of %v, got %d GB", possibleSizes, ramInfoValue)
	}
}

// Test GetDriveInfo function
func TestGetDriveInfo(t *testing.T) {
	// Mock exec.Command to return a predefined output
	pkg.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "sda  931.5G  /")
	}

	driveInfo := pkg.GetDriveInfo()
	expected := "sda, 931.5G, /"
	if driveInfo != expected {
		t.Errorf("Expected drive info to be '%s', got '%s'", expected, driveInfo)
	}
}

// Test GetDriveUsage function with tolerance range
func TestGetDriveUsage(t *testing.T) {
	// Mock exec.Command to return a predefined output
	pkg.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "Filesystem      Size  Used Avail Use% Mounted on\n/dev/sda1       931G  100G  831G  11% /")
	}

	driveUsage := pkg.GetDriveUsage()
	expectedMin := 0.0   // Minimum expected value
	expectedMax := 100.0 // Maximum expected value

	driveUsageValue, err := strconv.ParseFloat(strings.TrimSuffix(driveUsage, "%"), 64)
	if err != nil {
		t.Fatalf("Failed to parse drive usage: %v", err)
	}

	if driveUsageValue < expectedMin || driveUsageValue > expectedMax {
		t.Errorf("Expected drive usage to be between %.2f%% and %.2f%%, got %.2f%%", expectedMin, expectedMax, driveUsageValue)
	}
}

// Test GetCPUUsage function with tolerance range
func TestGetCPUUsage(t *testing.T) {
	// Mock exec.Command to return a predefined output
	pkg.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "cpu  3357 0 4313 1362393 0 0 0 0 0 0")
	}

	cpuUsage := pkg.GetCPUUsage()
	expectedMin := 0.0   // Minimum expected value
	expectedMax := 100.0 // Maximum expected value

	cpuUsageValue, err := strconv.ParseFloat(strings.TrimSuffix(cpuUsage, "%"), 64)
	if err != nil {
		t.Fatalf("Failed to parse CPU usage: %v", err)
	}

	if cpuUsageValue < expectedMin || cpuUsageValue > expectedMax {
		t.Errorf("Expected CPU usage to be between %.2f%% and %.2f%%, got %.2f%%", expectedMin, expectedMax, cpuUsageValue)
	}
}

// Test GetGPUInfo function
func TestGetGPUInfo(t *testing.T) {
	// Mock exec.Command to return a predefined output
	pkg.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "NVIDIA GeForce GTX 1050")
	}

	gpuInfo := pkg.GetGPUInfo()
	expected := "NVIDIA GeForce GTX 1050"
	if gpuInfo != expected {
		t.Errorf("Expected GPU info to be '%s', got '%s'", expected, gpuInfo)
	}
}

// Test GetGPUUsage function with tolerance range
func TestGetGPUUsage(t *testing.T) {
	// Mock exec.Command to return a predefined output
	pkg.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", "50")
	}

	gpuUsage := pkg.GetGPUUsage()
	expectedMin := 0.0   // Minimum expected value
	expectedMax := 100.0 // Maximum expected value

	gpuUsageValue, err := strconv.ParseFloat(strings.TrimSuffix(gpuUsage, "%"), 64)
	if err != nil {
		t.Fatalf("Failed to parse GPU usage: %v", err)
	}

	if gpuUsageValue < expectedMin || gpuUsageValue > expectedMax {
		t.Errorf("Expected GPU usage to be between %.2f%% and %.2f%%, got %.2f%%", expectedMin, expectedMax, gpuUsageValue)
	}
}

// Test GetUptime function
func TestGetUptime(t *testing.T) {
	// Mock exec.Command to return a predefined output
	pkg.ExecCommand = func(name string, arg ...string) *exec.Cmd {
		return exec.Command("echo", " 10:10:10 up 10 days,  1:01,  1 user,  load average: 0.00, 0.01, 0.05")
	}

	uptime := pkg.GetUptime()
	expected := "10 days"
	if uptime != expected {
		t.Errorf("Expected uptime to be '%s', got '%s'", expected, uptime)
	}
}

// Test GetResolution function with format validation
func TestGetResolution(t *testing.T) {
	resolution := pkg.GetResolution()
	matched, err := regexp.MatchString(`^\d+x\d+$`, resolution)
	if err != nil {
		t.Fatalf("Failed to compile regex: %v", err)
	}
	if !matched {
		t.Errorf("Expected resolution to match format 'yyyyyxyyyyy', got '%s'", resolution)
	}
}

// Test GetPublicIP function
func TestGetPublicIP(t *testing.T) {
	// Mock http.Get to return a predefined response
	pkg.HttpGet = func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(strings.NewReader("8.8.8.8")),
		}, nil
	}

	publicIP := pkg.GetPublicIP()
	expected := "8.8.8.8"
	if publicIP != expected {
		t.Errorf("Expected public IP to be '%s', got '%s'", expected, publicIP)
	}
}
