package pkg

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const defaultConfigValue = "Replace This"

// OSRelease structure
type OSRelease struct {
	ANSI_COLOR        string
	BUG_REPORT_URL    string
	BUILD_ID          string
	CPE_NAME          string
	DEFAULT_HOSTNAME  string
	DOCUMENTATION_URL string
	HOME_URL          string
	ID                string
	ID_LIKE           string
	IMAGE_ID          string
	IMAGE_VERSION     string
	LOGO              string
	NAME              string
	PRETTY_NAME       string
	SUPPORT_URL       string
	VARIANT           string
	VARIANT_ID        string
	VENDOR_NAME       string
	VENDOR_URL        string
	VERSION           string
	VERSION_CODENAME  string
	VERSION_ID        string
}

var (
	ReadFile    = os.ReadFile
	ExecCommand = exec.Command
	OpenFile    = os.Open
	CurrentUser = user.Current
	Hostname    = os.Hostname
	LookupEnv   = os.LookupEnv
	ReadDir     = os.ReadDir
	ReadAll     = io.ReadAll
	HttpGet     = http.Get
)

// GetOsRelease reads the OS release information from a reader
func GetOsRelease(reader io.Reader) OSRelease {
	osRelease := OSRelease{}
	data := make(map[string]string)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := strings.Trim(parts[1], `"`)

		switch key {
		case "ANSI_COLOR":
			osRelease.ANSI_COLOR = value
			data["os_ansi_color"] = value
		case "BUG_REPORT_URL":
			osRelease.BUG_REPORT_URL = value
			data["os_bug_report_url"] = value
		case "BUILD_ID":
			osRelease.BUILD_ID = value
			data["os_build_id"] = value
		case "CPE_NAME":
			osRelease.CPE_NAME = value
			data["os_cpe_name"] = value
		case "DEFAULT_HOSTNAME":
			osRelease.DEFAULT_HOSTNAME = value
			data["os_default_hostname"] = value
		case "DOCUMENTATION_URL":
			osRelease.DOCUMENTATION_URL = value
			data["os_documentation_url"] = value
		case "HOME_URL":
			osRelease.HOME_URL = value
			data["os_home_url"] = value
		case "ID":
			osRelease.ID = value
			data["os_id"] = value
		case "ID_LIKE":
			osRelease.ID_LIKE = value
			data["os_id_like"] = value
		case "IMAGE_ID":
			osRelease.IMAGE_ID = value
			data["os_image_id"] = value
		case "IMAGE_VERSION":
			osRelease.IMAGE_VERSION = value
			data["os_image_version"] = value
		case "LOGO":
			osRelease.LOGO = value
			data["os_logo"] = value
		case "NAME":
			osRelease.NAME = value
			data["os_name"] = value
		case "PRETTY_NAME":
			osRelease.PRETTY_NAME = value
			data["os_pretty_name"] = value
		case "SUPPORT_URL":
			osRelease.SUPPORT_URL = value
			data["os_support_url"] = value
		case "VARIANT":
			osRelease.VARIANT = value
			data["os_variant"] = value
		case "VARIANT_ID":
			osRelease.VARIANT_ID = value
			data["os_variant_id"] = value
		case "VENDOR_NAME":
			osRelease.VENDOR_NAME = value
			data["os_vendor_name"] = value
		case "VENDOR_URL":
			osRelease.VENDOR_URL = value
			data["os_vendor_url"] = value
		case "VERSION":
			osRelease.VERSION = value
			data["os_version"] = value
		case "VERSION_CODENAME":
			osRelease.VERSION_CODENAME = value
			data["os_version_codename"] = value
		case "VERSION_ID":
			osRelease.VERSION_ID = value
			data["os_version_id"] = value
		}
	}

	// Save the data to datafile.json
	SaveDataToFile(data, "data/datafile.json")

	return osRelease
}

// GetDesktopManager returns the name of the desktop manager
func GetDesktopManager() string {
	envVars := []string{
		"XDG_CURRENT_DESKTOP",
		"DESKTOP_SESSION",
		"GDMSESSION",
		"XDG_SESSION_DESKTOP",
	}
	value := GetEnvVar(envVars)
	if value != defaultConfigValue {
		SaveDataToFile(map[string]string{"dm": value}, "data/datafile.json")
		return value
	}

	systemFiles := map[string]string{
		"/etc/xdg/autostart/gnome-session-properties.desktop": "GNOME",
		"/etc/xdg/autostart/kdeinit5.desktop":                 "KDE",
		"/etc/xdg/autostart/lxsession.desktop":                "LXDE",
		"/etc/xdg/autostart/xfce4-session.desktop":            "XFCE",
	}
	for file, dm := range systemFiles {
		if _, err := os.Stat(file); err == nil {
			SaveDataToFile(map[string]string{"dm": dm}, "data/datafile.json")
			return dm
		}
	}

	processes := map[string]string{
		"gnome-session": "GNOME",
		"ksmserver":     "KDE",
		"lxsession":     "LXDE",
		"xfce4-session": "XFCE",
	}
	SaveDataToFile(map[string]string{"dm": GetRunningProcess(processes)}, "data/datafile.json")
	return GetRunningProcess(processes)
}

func GetUsername() string {
	user, err := user.Current()
	if err != nil {
		SaveDataToFile(map[string]string{"user": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}
	value := strings.TrimRight(string(user.Username), "\x00")
	SaveDataToFile(map[string]string{"user": value}, "data/datafile.json")
	return value
}

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		SaveDataToFile(map[string]string{"hostname": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}
	value := strings.TrimRight(string(hostname), "\x00")
	SaveDataToFile(map[string]string{"hostname": value}, "data/datafile.json")
	return value
}

func GetKernelVersion() string {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		SaveDataToFile(map[string]string{"kernel": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}
	value := charsToString(uname.Release)
	SaveDataToFile(map[string]string{"kernel": value}, "data/datafile.json")
	return value
}

func GetShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		SaveDataToFile(map[string]string{"shell": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}
	SaveDataToFile(map[string]string{"shell": shell}, "data/datafile.json")
	return shell
}

func GetTerminal() string {
	envVars := []string{
		"TERM_PROGRAM",
		"COLORTERM",
		"TERM",
	}
	value := GetEnvVar(envVars)
	if value != defaultConfigValue {
		SaveDataToFile(map[string]string{"term": value}, "data/datafile.json")
		return value
	}

	processes := map[string]string{
		"gnome-terminal": "GNOME Terminal",
		"konsole":        "Konsole",
		"xfce4-terminal": "XFCE Terminal",
		"lxterminal":     "LXTerminal",
		"alacritty":      "Alacritty",
		"kitty":          "Kitty",
		"terminator":     "Terminator",
		"st":             "st",
		"urxvt":          "URxvt",
		"xterm":          "xterm",
		"tmux":           "tmux",
		"screen":         "screen",
		"ghostty":        "Ghostty",
	}
	SaveDataToFile(map[string]string{"term": GetRunningProcess(processes)}, "data/datafile.json")
	return GetRunningProcess(processes)
}

// Cache the contents of /proc/meminfo to avoid repeated reads
var memInfoCache string
var memInfoOnce sync.Once
var cpuInfoCache string
var cpuInfoOnce sync.Once

func readMemInfo() string {
	memInfoOnce.Do(func() {
		data, err := os.ReadFile("/proc/meminfo")
		if err != nil {
			memInfoCache = ""
		} else {
			memInfoCache = string(data)
		}
	})
	return memInfoCache
}

func readCPUInfo() string {
	cpuInfoOnce.Do(func() {
		data, err := os.ReadFile("/proc/cpuinfo")
		if err != nil {
			cpuInfoCache = ""
		} else {
			cpuInfoCache = string(data)
		}
	})
	SaveDataToFile(map[string]string{"cpu": cpuInfoCache}, "data/datafile.json")
	return cpuInfoCache
}

func GetCPUInfo() string {
	cpuInfo := readCPUInfo()
	scanner := bufio.NewScanner(strings.NewReader(cpuInfo))

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "model name") {
			fields := strings.Split(line, ":")
			if len(fields) > 1 {
				value := strings.TrimSpace(fields[1])
				SaveDataToFile(map[string]string{"cpu": value}, "data/datafile.json")
				return value
			}
		}
	}
	SaveDataToFile(map[string]string{"cpu": defaultConfigValue}, "data/datafile.json")
	return defaultConfigValue
}

func GetRAMUsage() string {
	memInfo := readMemInfo()
	if memInfo == "" {
		SaveDataToFile(map[string]string{"ram %": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}

	scanner := bufio.NewScanner(strings.NewReader(memInfo))
	memTotal := uint64(0)
	memAvailable := uint64(0)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			memTotal, _ = strconv.ParseUint(fields[1], 10, 64)
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			memAvailable, _ = strconv.ParseUint(fields[1], 10, 64)
		}
	}

	if memTotal == 0 {
		SaveDataToFile(map[string]string{"ram %": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}

	memUsed := memTotal - memAvailable
	ramUsage := (float64(memUsed) / float64(memTotal)) * 100.0

	value := fmt.Sprintf("%.2f%%", ramUsage)
	SaveDataToFile(map[string]string{"ram %": value}, "data/datafile.json")
	return value
}

func GetRAMInfo() string {
	memInfo := readMemInfo()
	if memInfo == "" {
		SaveDataToFile(map[string]string{"ram": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}

	scanner := bufio.NewScanner(strings.NewReader(memInfo))
	memTotal := uint64(0)

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			memTotal, _ = strconv.ParseUint(fields[1], 10, 64)
			break
		}
	}

	if memTotal == 0 {
		SaveDataToFile(map[string]string{"ram": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}

	// Convert kB to MB
	memTotalMB := memTotal / 1024

	possibleSizes := []uint64{1, 2, 4, 8, 16, 32, 64, 128, 256, 512, 1024, 2048, 4096, 8192, 16384, 32768, 65536}

	closestSize := possibleSizes[0]
	for _, size := range possibleSizes {
		if Abs(int64(size)-int64(memTotalMB)) < Abs(int64(closestSize)-int64(memTotalMB)) {
			closestSize = size
		}
	}

	if closestSize >= 1024 {
		value := fmt.Sprintf("%d GB", closestSize/1024)
		SaveDataToFile(map[string]string{"ram": value}, "data/datafile.json")
		return value
	}
	value := fmt.Sprintf("%d MB", closestSize)
	SaveDataToFile(map[string]string{"ram": value}, "data/datafile.json")
	return value
}

func GetDriveInfo() string {
	cmd := ExecCommand("lsblk", "-o", "NAME,SIZE,MOUNTPOINT")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"drive": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[2] == "/" {
			name := strings.TrimPrefix(fields[0], "├─")
			name = strings.TrimPrefix(name, "└─")
			value := fmt.Sprintf("%s, %s, %s", name, fields[1], fields[2])
			SaveDataToFile(map[string]string{"drive": value}, "data/datafile.json")
			return value
		}
	}

	SaveDataToFile(map[string]string{"drive": defaultConfigValue}, "data/datafile.json")
	return defaultConfigValue
}

func GetDriveUsage() string {
	cmd := ExecCommand("df", "-h", "/")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"drive %": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		fields := strings.Fields(lines[1])
		if len(fields) >= 5 {
			value := fields[4]
			SaveDataToFile(map[string]string{"drive %": value}, "data/datafile.json")
			return value
		}
	}

	SaveDataToFile(map[string]string{"drive %": defaultConfigValue}, "data/datafile.json")
	return defaultConfigValue
}

func GetCPUUsage() string {
	idle0, total0 := GetCPUSample()
	time.Sleep(1 * time.Second)
	idle1, total1 := GetCPUSample()

	idleTicks := float64(idle1 - idle0)
	totalTicks := float64(total1 - total0)

	cpuUsage := (1.0 - idleTicks/totalTicks) * 100.0

	value := fmt.Sprintf("%.2f%%", cpuUsage)
	SaveDataToFile(map[string]string{"cpu %": value}, "data/datafile.json")
	return value
}

func GetCPUSample() (uint64, uint64) {
	data, err := OpenFile("/proc/stat")
	if err != nil {
		return 0, 0
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)
			idle, _ := strconv.ParseUint(fields[4], 10, 64)
			total := uint64(0)
			for _, field := range fields[1:] {
				value, _ := strconv.ParseUint(field, 10, 64)
				total += value
			}
			return idle, total
		}
	}

	return 0, 0
}

func GetGPUInfo() string {
	if isCommandAvailable("nvidia-smi") {
		return GetNvidiaGPUInfo()
	} else if isCommandAvailable("rocm-smi") {
		return GetAmdGPUInfo()
	} else if isCommandAvailable("intel_gpu_top") {
		return GetIntelGPUInfo()
	}
	SaveDataToFile(map[string]string{"gpu": defaultConfigValue}, "data/datafile.json")
	return defaultConfigValue
}

func GetGPUUsage() string {
	if isCommandAvailable("nvidia-smi") {
		return GetNvidiaGPUUsage()
	} else if isCommandAvailable("rocm-smi") {
		return GetAmdGPUUsage()
	} else if isCommandAvailable("intel_gpu_top") {
		return GetIntelGPUUsage()
	}
	SaveDataToFile(map[string]string{"gpu %": defaultConfigValue}, "data/datafile.json")
	return defaultConfigValue
}

func GetNvidiaGPUInfo() string {
	cmd := ExecCommand("nvidia-smi", "--query-gpu=name", "--format=csv,noheader")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"gpu": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}
	value := strings.TrimSpace(string(output))
	SaveDataToFile(map[string]string{"gpu": value}, "data/datafile.json")
	return value
}

func GetNvidiaGPUUsage() string {
	cmd := ExecCommand("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"gpu %": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}
	value := strings.TrimSpace(string(output)) + "%"
	SaveDataToFile(map[string]string{"gpu %": value}, "data/datafile.json")
	return value
}

func GetAmdGPUInfo() string {
	cmd := ExecCommand("rocm-smi", "--showproductname")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"gpu": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		value := strings.TrimSpace(lines[1])
		SaveDataToFile(map[string]string{"gpu": value}, "data/datafile.json")
		return value
	}
	SaveDataToFile(map[string]string{"gpu": defaultConfigValue}, "data/datafile.json")
	return defaultConfigValue
}

func GetAmdGPUUsage() string {
	cmd := ExecCommand("rocm-smi", "--showuse")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"gpu %": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "GPU use") {
			fields := strings.Fields(line)
			if len(fields) > 2 {
				value := fields[2]
				SaveDataToFile(map[string]string{"gpu %": value}, "data/datafile.json")
				return value
			}
		}
	}
	SaveDataToFile(map[string]string{"gpu %": defaultConfigValue}, "data/datafile.json")
	return defaultConfigValue
}

func GetIntelGPUInfo() string {
	cmd := ExecCommand("lspci", "-nn", "-d", "8086:")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"gpu": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "VGA compatible controller") {
			value := strings.TrimSpace(line)
			SaveDataToFile(map[string]string{"gpu": value}, "data/datafile.json")
			return value
		}
	}

	SaveDataToFile(map[string]string{"gpu": defaultConfigValue}, "data/datafile.json")
	return defaultConfigValue
}

func GetIntelGPUUsage() string {
	cmd := ExecCommand("sh", "-c", "timeout 1s intel_gpu_top -o - | grep 'Render/3D' | awk '{print $2}'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		SaveDataToFile(map[string]string{"gpu %": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}
	value := strings.TrimSpace(out.String())
	SaveDataToFile(map[string]string{"gpu %": value}, "data/datafile.json")
	return value
}

func GetUptime() string {
	cmd := ExecCommand("uptime")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		SaveDataToFile(map[string]string{"uptime": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}

	uptimeOutput := out.String()
	parts := strings.Split(uptimeOutput, "up ")
	if len(parts) < 2 {
		SaveDataToFile(map[string]string{"uptime": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}

	uptimePart := strings.Split(parts[1], ",")[0]

	value := strings.TrimSpace(uptimePart)

	SaveDataToFile(map[string]string{"uptime": value}, "data/datafile.json")
	return value
}

func GetWM() string {
	envVars := []string{
		"XDG_SESSION_DESKTOP",
		"XDG_CURRENT_DESKTOP",
		"DESKTOP_SESSION",
	}
	value := GetEnvVar(envVars)
	if value != defaultConfigValue {
		SaveDataToFile(map[string]string{"wm": value}, "data/datafile.json")
		return value
	}

	processes := map[string]string{
		"i3":           "i3",
		"bspwm":        "bspwm",
		"openbox":      "Openbox",
		"awesome":      "Awesome",
		"herbstluftwm": "herbstluftwm",
		"fluxbox":      "Fluxbox",
		"icewm":        "IceWM",
		"kwin":         "KWin",
		"mutter":       "Mutter",
		"compiz":       "Compiz",
	}

	SaveDataToFile(map[string]string{"wm": GetRunningProcess(processes)}, "data/datafile.json")
	return GetRunningProcess(processes)
}

func GetResolution() string {
	cmd := exec.Command("xrandr")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		SaveDataToFile(map[string]string{"resolution": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}

	lines := strings.Split(out.String(), "\n")
	for _, line := range lines {
		if strings.Contains(line, "current") {
			fields := strings.Fields(line)
			for i, field := range fields {
				if field == "current" && i+2 < len(fields) {
					resolution := fields[i+1] + "x" + fields[i+3]
					value := strings.TrimSuffix(resolution, ",")
					SaveDataToFile(map[string]string{"resolution": value}, "data/datafile.json")
					return value
				}
			}
		}
	}

	SaveDataToFile(map[string]string{"resolution": defaultConfigValue}, "data/datafile.json")
	return defaultConfigValue
}

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		SaveDataToFile(map[string]string{"ip": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				SaveDataToFile(map[string]string{"ip": string(ipnet.IP.String())}, "data/datafile.json")
				return ipnet.IP.String()
			}
		}
	}

	SaveDataToFile(map[string]string{"ip": defaultConfigValue}, "data/datafile.json")
	return defaultConfigValue
}

func GetRunningProcessesCount() string {
	procDir := "/proc"
	entries, err := os.ReadDir(procDir)
	if err != nil {
		fmt.Println("Error reading /proc directory:", err)
		SaveDataToFile(map[string]string{"processes": "0"}, "data/datafile.json")
		return "0"
	}

	count := 0
	for _, entry := range entries {
		if entry.IsDir() {
			if _, err := strconv.Atoi(entry.Name()); err == nil {
				count++
			}
		}
	}

	value := fmt.Sprintf("%d", count)
	SaveDataToFile(map[string]string{"processes": value}, "data/datafile.json")
	return value
}

func GetPublicIP() string {
	resp, err := HttpGet("https://api.ipify.org?format=text")
	if err != nil {
		SaveDataToFile(map[string]string{"public ip": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}
	defer resp.Body.Close()

	ip, err := ReadAll(resp.Body)
	if err != nil {
		SaveDataToFile(map[string]string{"public ip": defaultConfigValue}, "data/datafile.json")
		return defaultConfigValue
	}

	SaveDataToFile(map[string]string{"public ip": string(ip)}, "data/datafile.json")
	return string(ip)
}

func MenuItems(config Config, usedatafile bool) map[string]string {
	items := make(map[string]string)
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	if usedatafile {
		datafile, err := os.Open("data/datafile.json")
		if err != nil {
			fmt.Println("Error opening data file:", err)
			return nil
		}
		defer datafile.Close()

		decoder := json.NewDecoder(datafile)
		if err := decoder.Decode(&items); err != nil {
			fmt.Println("Error decoding data file:", err)
			return nil
		}
	} else {
		osReleaseFile, err := os.Open("/etc/os-release")
		if err != nil {
			return nil
		}
		defer osReleaseFile.Close()

		osRelease := GetOsRelease(osReleaseFile)

		for _, item := range config.Items {
			wg.Add(1)
			go func(item ConfigItem) {
				defer wg.Done()
				var value string
				switch item.Keyword {
				case "os_ansi_color":
					value = osRelease.ANSI_COLOR
				case "os_pretty_name":
					value = osRelease.PRETTY_NAME
				case "os_bug_report_url":
					value = osRelease.BUG_REPORT_URL
				case "os_build_id":
					value = osRelease.BUILD_ID
				case "os_cpe_name":
					value = osRelease.CPE_NAME
				case "os_default_hostname":
					value = osRelease.DEFAULT_HOSTNAME
				case "os_documentation_url":
					value = osRelease.DOCUMENTATION_URL
				case "os_home_url":
					value = osRelease.HOME_URL
				case "os_id":
					value = osRelease.ID
				case "os_id_like":
					value = osRelease.ID_LIKE
				case "os_image_id":
					value = osRelease.IMAGE_ID
				case "os_image_version":
					value = osRelease.IMAGE_VERSION
				case "os_version":
					value = osRelease.VERSION
				case "os_logo":
					value = osRelease.LOGO
				case "os_name":
					value = osRelease.NAME
				case "os_support_url":
					value = osRelease.SUPPORT_URL
				case "os_variant":
					value = osRelease.VARIANT
				case "os_variant_id":
					value = osRelease.VARIANT_ID
				case "os_vendor_name":
					value = osRelease.VENDOR_NAME
				case "os_vendor_url":
					value = osRelease.VENDOR_URL
				case "os_version_codename":
					value = osRelease.VERSION_CODENAME
				case "os_version_id":
					value = osRelease.VERSION_ID
				case "user":
					value = GetUsername()
				case "hostname":
					value = GetHostname()
				case "kernel":
					value = GetKernelVersion()
				case "shell":
					value = GetShell()
				case "uptime":
					value = GetUptime()
				case "dm":
					value = GetDesktopManager()
				case "gpu":
					value = GetGPUInfo()
				case "cpu":
					value = GetCPUInfo()
				case "ram":
					value = GetRAMInfo()
				case "drive":
					value = GetDriveInfo()
				case "gpu %":
					value = GetGPUUsage()
				case "cpu %":
					value = GetCPUUsage()
				case "ram %":
					value = GetRAMUsage()
				case "drive %":
					value = GetDriveUsage()
				case "term":
					value = GetTerminal()
				case "processes":
					value = GetRunningProcessesCount()
				case "wm":
					value = GetWM()
				case "ip":
					value = GetIP()
				case "public ip":
					value = GetPublicIP()
				case "resolution":
					value = GetResolution()
				}
				mu.Lock()
				items[item.Keyword] = value
				mu.Unlock()
			}(item)
		}
		wg.Wait()
	}

	return items
}
