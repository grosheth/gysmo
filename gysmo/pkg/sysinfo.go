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
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"
)

const defaultConfigValue = "Not Found"

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

func GetOsRelease(reader io.Reader) OSRelease {
	osRelease := OSRelease{}
	data := make(map[string]string)
	scanner := bufio.NewScanner(reader)

	fieldMap := map[string]func(string){
		"ANSI_COLOR":        func(value string) { osRelease.ANSI_COLOR = value; data["os_ansi_color"] = value },
		"BUG_REPORT_URL":    func(value string) { osRelease.BUG_REPORT_URL = value; data["os_bug_report_url"] = value },
		"BUILD_ID":          func(value string) { osRelease.BUILD_ID = value; data["os_build_id"] = value },
		"CPE_NAME":          func(value string) { osRelease.CPE_NAME = value; data["os_cpe_name"] = value },
		"DEFAULT_HOSTNAME":  func(value string) { osRelease.DEFAULT_HOSTNAME = value; data["os_default_hostname"] = value },
		"DOCUMENTATION_URL": func(value string) { osRelease.DOCUMENTATION_URL = value; data["os_documentation_url"] = value },
		"HOME_URL":          func(value string) { osRelease.HOME_URL = value; data["os_home_url"] = value },
		"ID":                func(value string) { osRelease.ID = value; data["os_id"] = value },
		"ID_LIKE":           func(value string) { osRelease.ID_LIKE = value; data["os_id_like"] = value },
		"IMAGE_ID":          func(value string) { osRelease.IMAGE_ID = value; data["os_image_id"] = value },
		"IMAGE_VERSION":     func(value string) { osRelease.IMAGE_VERSION = value; data["os_image_version"] = value },
		"LOGO":              func(value string) { osRelease.LOGO = value; data["os_logo"] = value },
		"NAME":              func(value string) { osRelease.NAME = value; data["os_name"] = value },
		"PRETTY_NAME":       func(value string) { osRelease.PRETTY_NAME = value; data["os_pretty_name"] = value },
		"SUPPORT_URL":       func(value string) { osRelease.SUPPORT_URL = value; data["os_support_url"] = value },
		"VARIANT":           func(value string) { osRelease.VARIANT = value; data["os_variant"] = value },
		"VARIANT_ID":        func(value string) { osRelease.VARIANT_ID = value; data["os_variant_id"] = value },
		"VENDOR_NAME":       func(value string) { osRelease.VENDOR_NAME = value; data["os_vendor_name"] = value },
		"VENDOR_URL":        func(value string) { osRelease.VENDOR_URL = value; data["os_vendor_url"] = value },
		"VERSION":           func(value string) { osRelease.VERSION = value; data["os_version"] = value },
		"VERSION_CODENAME":  func(value string) { osRelease.VERSION_CODENAME = value; data["os_version_codename"] = value },
		"VERSION_ID":        func(value string) { osRelease.VERSION_ID = value; data["os_version_id"] = value },
	}

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := parts[0]
		value := strings.Trim(parts[1], `"`)

		if assignFunc, exists := fieldMap[key]; exists {
			assignFunc(value)
		}
	}

	SaveDataToFile(data)

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
		SaveDataToFile(map[string]string{"dm": value})
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
			SaveDataToFile(map[string]string{"dm": dm})
			return dm
		}
	}

	processes := map[string]string{
		"gnome-session": "GNOME",
		"ksmserver":     "KDE",
		"lxsession":     "LXDE",
		"xfce4-session": "XFCE",
	}
	SaveDataToFile(map[string]string{"dm": GetRunningProcess(processes)})
	return GetRunningProcess(processes)
}

func GetUsername() string {
	user, err := user.Current()
	if err != nil {
		SaveDataToFile(map[string]string{"user": defaultConfigValue})
		return defaultConfigValue
	}
	value := strings.TrimRight(string(user.Username), "\x00")
	SaveDataToFile(map[string]string{"user": value})
	return value
}

func GetHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		SaveDataToFile(map[string]string{"hostname": defaultConfigValue})
		return defaultConfigValue
	}
	value := strings.TrimRight(string(hostname), "\x00")
	SaveDataToFile(map[string]string{"hostname": value})
	return value
}

func GetKernelVersion() string {
	var uname syscall.Utsname
	if err := syscall.Uname(&uname); err != nil {
		SaveDataToFile(map[string]string{"kernel": defaultConfigValue})
		return defaultConfigValue
	}
	value := CharsToString(uname.Release)
	SaveDataToFile(map[string]string{"kernel": value})
	return value
}

func GetShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		SaveDataToFile(map[string]string{"shell": defaultConfigValue})
		return defaultConfigValue
	}
	SaveDataToFile(map[string]string{"shell": shell})
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
		SaveDataToFile(map[string]string{"term": value})
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
	SaveDataToFile(map[string]string{"term": GetRunningProcess(processes)})
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

func GetMotherboardInfo() string {
	data, err := os.ReadFile("/sys/class/dmi/id/board_name")
	if err != nil {
		return defaultConfigValue
	}
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	var motherboardInfo string
	for scanner.Scan() {
		line := scanner.Text()
		motherboardInfo = strings.TrimSpace(line)
	}

	if err := scanner.Err(); err != nil {
		return defaultConfigValue
	}

	return motherboardInfo
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
				SaveDataToFile(map[string]string{"cpu": value})
				return value
			}
		}
	}
	SaveDataToFile(map[string]string{"cpu": defaultConfigValue})
	return defaultConfigValue
}

func GetRAMUsage() string {
	memInfo := readMemInfo()
	if memInfo == "" {
		SaveDataToFile(map[string]string{"ram %": defaultConfigValue})
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
		SaveDataToFile(map[string]string{"ram %": defaultConfigValue})
		return defaultConfigValue
	}

	memUsed := memTotal - memAvailable
	ramUsage := (float64(memUsed) / float64(memTotal)) * 100.0

	value := fmt.Sprintf("%.2f%%", ramUsage)
	SaveDataToFile(map[string]string{"ram %": value})
	return value
}

func GetRAMInfo() string {
	memInfo := readMemInfo()
	if memInfo == "" {
		SaveDataToFile(map[string]string{"ram": defaultConfigValue})
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
		SaveDataToFile(map[string]string{"ram": defaultConfigValue})
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
		SaveDataToFile(map[string]string{"ram": value})
		return value
	}
	value := fmt.Sprintf("%d MB", closestSize)
	SaveDataToFile(map[string]string{"ram": value})
	return value
}

func GetDriveInfo() string {
	cmd := ExecCommand("lsblk", "-o", "NAME,SIZE,MOUNTPOINT")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"drive": defaultConfigValue})
		return defaultConfigValue
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 3 && fields[2] == "/" {
			name := strings.TrimPrefix(fields[0], "├─")
			name = strings.TrimPrefix(name, "└─")
			value := fmt.Sprintf("%s, %s, %s", name, fields[1], fields[2])
			SaveDataToFile(map[string]string{"drive": value})
			return value
		}
	}

	SaveDataToFile(map[string]string{"drive": defaultConfigValue})
	return defaultConfigValue
}

func GetDriveUsage() string {
	cmd := ExecCommand("df", "-h", "/")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"drive %": defaultConfigValue})
		return defaultConfigValue
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		fields := strings.Fields(lines[1])
		if len(fields) >= 5 {
			value := fields[4]
			SaveDataToFile(map[string]string{"drive %": value})
			return value
		}
	}

	SaveDataToFile(map[string]string{"drive %": defaultConfigValue})
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
	SaveDataToFile(map[string]string{"cpu %": value})
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
	if IsCommandAvailable("nvidia-smi") {
		return GetNvidiaGPUInfo()
	} else if IsCommandAvailable("rocm-smi") {
		return GetAmdGPUInfo()
	} else if IsCommandAvailable("intel_gpu_top") {
		return GetIntelGPUInfo()
	}
	SaveDataToFile(map[string]string{"gpu": defaultConfigValue})
	return defaultConfigValue
}

func GetGPUUsage() string {
	if IsCommandAvailable("nvidia-smi") {
		return GetNvidiaGPUUsage()
	} else if IsCommandAvailable("rocm-smi") {
		return GetAmdGPUUsage()
	} else if IsCommandAvailable("intel_gpu_top") {
		return GetIntelGPUUsage()
	}
	SaveDataToFile(map[string]string{"gpu %": defaultConfigValue})
	return defaultConfigValue
}

func GetNvidiaGPUInfo() string {
	cmd := ExecCommand("nvidia-smi", "--query-gpu=name", "--format=csv,noheader")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"gpu": defaultConfigValue})
		return defaultConfigValue
	}
	value := strings.TrimSpace(string(output))
	SaveDataToFile(map[string]string{"gpu": value})
	return value
}

func GetNvidiaGPUUsage() string {
	cmd := ExecCommand("nvidia-smi", "--query-gpu=utilization.gpu", "--format=csv,noheader,nounits")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"gpu %": defaultConfigValue})
		return defaultConfigValue
	}
	value := strings.TrimSpace(string(output)) + "%"
	SaveDataToFile(map[string]string{"gpu %": value})
	return value
}

func GetAmdGPUInfo() string {
	cmd := ExecCommand("rocm-smi", "--showproductname")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"gpu": defaultConfigValue})
		return defaultConfigValue
	}
	lines := strings.Split(string(output), "\n")
	if len(lines) > 1 {
		value := strings.TrimSpace(lines[1])
		SaveDataToFile(map[string]string{"gpu": value})
		return value
	}
	SaveDataToFile(map[string]string{"gpu": defaultConfigValue})
	return defaultConfigValue
}

func GetAmdGPUUsage() string {
	cmd := ExecCommand("rocm-smi", "--showuse")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"gpu %": defaultConfigValue})
		return defaultConfigValue
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "GPU use") {
			fields := strings.Fields(line)
			if len(fields) > 2 {
				value := fields[2]
				SaveDataToFile(map[string]string{"gpu %": value})
				return value
			}
		}
	}
	SaveDataToFile(map[string]string{"gpu %": defaultConfigValue})
	return defaultConfigValue
}

func GetIntelGPUInfo() string {
	cmd := ExecCommand("lspci", "-nn", "-d", "8086:")
	output, err := cmd.Output()
	if err != nil {
		SaveDataToFile(map[string]string{"gpu": defaultConfigValue})
		return defaultConfigValue
	}
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "VGA compatible controller") {
			value := strings.TrimSpace(line)
			SaveDataToFile(map[string]string{"gpu": value})
			return value
		}
	}

	SaveDataToFile(map[string]string{"gpu": defaultConfigValue})
	return defaultConfigValue
}

func GetIntelGPUUsage() string {
	cmd := ExecCommand("sh", "-c", "timeout 1s intel_gpu_top -o - | grep 'Render/3D' | awk '{print $2}'")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		SaveDataToFile(map[string]string{"gpu %": defaultConfigValue})
		return defaultConfigValue
	}
	value := strings.TrimSpace(out.String())
	SaveDataToFile(map[string]string{"gpu %": value})
	return value
}

func GetUptime() string {
	cmd := ExecCommand("uptime")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		SaveDataToFile(map[string]string{"uptime": defaultConfigValue})
		return defaultConfigValue
	}

	uptimeOutput := out.String()
	parts := strings.Split(uptimeOutput, "up ")
	if len(parts) < 2 {
		SaveDataToFile(map[string]string{"uptime": defaultConfigValue})
		return defaultConfigValue
	}

	uptimePart := strings.Split(parts[1], ",")[0]

	value := strings.TrimSpace(uptimePart)

	SaveDataToFile(map[string]string{"uptime": value})
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
		SaveDataToFile(map[string]string{"wm": value})
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

	SaveDataToFile(map[string]string{"wm": GetRunningProcess(processes)})
	return GetRunningProcess(processes)
}

func GetResolution() string {
	cmd := exec.Command("xrandr")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		SaveDataToFile(map[string]string{"resolution": defaultConfigValue})
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
					SaveDataToFile(map[string]string{"resolution": value})
					return value
				}
			}
		}
	}

	SaveDataToFile(map[string]string{"resolution": defaultConfigValue})
	return defaultConfigValue
}

func GetIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		SaveDataToFile(map[string]string{"ip": defaultConfigValue})
		return defaultConfigValue
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				SaveDataToFile(map[string]string{"ip": string(ipnet.IP.String())})
				return ipnet.IP.String()
			}
		}
	}

	SaveDataToFile(map[string]string{"ip": defaultConfigValue})
	return defaultConfigValue
}

func GetRunningProcessesCount() string {
	procDir := "/proc"
	entries, err := os.ReadDir(procDir)
	if err != nil {
		fmt.Println("Error reading /proc directory:", err)
		SaveDataToFile(map[string]string{"processes": "0"})
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
	SaveDataToFile(map[string]string{"processes": value})
	return value
}

func GetPublicIP() string {
	resp, err := HttpGet("https://api.ipify.org?format=text")
	if err != nil {
		SaveDataToFile(map[string]string{"public ip": defaultConfigValue})
		return defaultConfigValue
	}
	defer resp.Body.Close()

	ip, err := ReadAll(resp.Body)
	if err != nil {
		SaveDataToFile(map[string]string{"public ip": defaultConfigValue})
		return defaultConfigValue
	}

	SaveDataToFile(map[string]string{"public ip": string(ip)})
	return string(ip)
}

func MenuItems(config Config, usedatafile bool) map[string]string {

	items := make(map[string]string)
	var wg sync.WaitGroup
	mu := &sync.Mutex{}

	workingPath := LoadWorkingPath()
	dataPath := filepath.Join(workingPath, "data", "data.json")
	if usedatafile {
		datafile, err := os.Open(dataPath)
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

		valueMap := map[string]func() string{
			"os_ansi_color":        func() string { return osRelease.ANSI_COLOR },
			"os_pretty_name":       func() string { return osRelease.PRETTY_NAME },
			"os_bug_report_url":    func() string { return osRelease.BUG_REPORT_URL },
			"os_build_id":          func() string { return osRelease.BUILD_ID },
			"os_cpe_name":          func() string { return osRelease.CPE_NAME },
			"os_default_hostname":  func() string { return osRelease.DEFAULT_HOSTNAME },
			"os_documentation_url": func() string { return osRelease.DOCUMENTATION_URL },
			"os_home_url":          func() string { return osRelease.HOME_URL },
			"os_id":                func() string { return osRelease.ID },
			"os_id_like":           func() string { return osRelease.ID_LIKE },
			"os_image_id":          func() string { return osRelease.IMAGE_ID },
			"os_image_version":     func() string { return osRelease.IMAGE_VERSION },
			"os_version":           func() string { return osRelease.VERSION },
			"os_logo":              func() string { return osRelease.LOGO },
			"os_name":              func() string { return osRelease.NAME },
			"os_support_url":       func() string { return osRelease.SUPPORT_URL },
			"os_variant":           func() string { return osRelease.VARIANT },
			"os_variant_id":        func() string { return osRelease.VARIANT_ID },
			"os_vendor_name":       func() string { return osRelease.VENDOR_NAME },
			"os_vendor_url":        func() string { return osRelease.VENDOR_URL },
			"os_version_codename":  func() string { return osRelease.VERSION_CODENAME },
			"os_version_id":        func() string { return osRelease.VERSION_ID },
			"user":                 GetUsername,
			"hostname":             GetHostname,
			"kernel":               GetKernelVersion,
			"shell":                GetShell,
			"uptime":               GetUptime,
			"dm":                   GetDesktopManager,
			"motherboard":          GetMotherboardInfo,
			"gpu":                  GetGPUInfo,
			"cpu":                  GetCPUInfo,
			"ram":                  GetRAMInfo,
			"drive":                GetDriveInfo,
			"gpu %":                GetGPUUsage,
			"cpu %":                GetCPUUsage,
			"ram %":                GetRAMUsage,
			"drive %":              GetDriveUsage,
			"term":                 GetTerminal,
			"processes":            GetRunningProcessesCount,
			"wm":                   GetWM,
			"ip":                   GetIP,
			"public ip":            GetPublicIP,
			"resolution":           GetResolution,
		}

		for _, item := range config.Items {
			wg.Add(1)
			go func(item ConfigItem) {
				defer wg.Done()
				if valueFunc, exists := valueMap[item.Keyword]; exists {
					value := valueFunc()
					mu.Lock()
					items[item.Keyword] = value
					mu.Unlock()
				}
			}(item)
		}
		wg.Wait()
	}

	return items
}
