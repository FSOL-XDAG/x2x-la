package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
)

type ProxyRecord struct {
	ID            uint   // Proxy ID  : 'Proxy [ID]' in logs
	DateLogin     uint64 // Unix Timestamp login (converted from layout "2006/01/02 15:04:05")
	DateLogoff    uint64 // Unix Timestamp logoff (converted from layout "2006/01/02 15:04:05")
	MinerName     string // Miner
	WalletAddress string // XDAG WalletAddress
}

const (
	// App Main informations
	AppName      string = "x2x LogAnalyzer"
	AppExe       string = "x2x-la"
	AppURL       string = "https://github.com/FSOL-XDAG/x2x-loganalyzer"
	AppHelp      string = "Getting more options & informations : --help"
	AppCopyright string = "Copyright (C) 2023 FSOL"

	// App Versioning
	AppVerMajor string = "0"
	AppVerMinor string = "1"
	AppVerPatch string = "0"
)

func DisplayProgramTitle() {
	// Define colors
	titleColor1 := color.New(color.FgCyan).SprintFunc()
	titleColor2 := color.New(color.FgHiCyan).SprintFunc()
	titleColor3 := color.New(color.FgCyan).SprintFunc()
	titleColor4 := color.New(color.FgWhite).SprintFunc()

	// Display text
	fmt.Printf("\n---< %s / %s >--- \n", titleColor2(AppName+" v"+AppVerMajor+"."+AppVerMinor+"."+AppVerPatch), titleColor1(AppCopyright))
	fmt.Printf("---<   %s   >---\n", titleColor3(AppURL))
	fmt.Printf("---<   %s   >---\n", titleColor4(AppHelp))
}

func DisplaySubTitle(title string) {
	// Define colors
	titleColor := color.New(color.FgHiYellow).SprintFunc()

	// Set vars
	titleLen := len(title)
	borderLen := 110 - titleLen - 4

	// Display text
	fmt.Printf("\n[ %s ]%s\n", titleColor(title), strings.Repeat("-", borderLen))
}

func DisplayItem(title string, value string) {
	// Define colors
	Color1 := color.New(color.FgHiRed).SprintFunc()
	Color2 := color.New(color.FgYellow).SprintFunc()

	// Display text
	fmt.Printf("  %s %s \t: %s\n", Color1("->"), title, Color2(value))
}

func DisplayError(title string) {
	// Define colors
	Color1 := color.New(color.FgHiRed).SprintFunc()
	Color2 := color.New(color.FgYellow).SprintFunc()

	// Display text
	fmt.Printf("\n  %s \t: %s\n", Color1("!!! Error !!!"), Color2(title))
}

func GetLogs(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		DisplayError("Unable to find '" + fileName + "' file.")
		return nil, fmt.Errorf("")
	}
	defer file.Close()

	lines := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// Filtering useless lines
		if !strings.Contains(line, "job blob:") &&
			!strings.Contains(line, "seed:") &&
			!strings.Contains(line, "--read:") &&
			!strings.Contains(line, "nonce:") &&
			!strings.Contains(line, "XDAG_FIELD_HEAD") &&
			!strings.Contains(line, "Goroutine") &&
			!strings.Contains(line, "new target:") &&
			!strings.Contains(line, "Broadcasting") {
			// Get Line
			lines = append(lines, line)
		}
	}

	return lines, nil
}

func GetProxyVer(lines []string) string {
	var ver string

	if len(lines) >= 3 {
		ver = lines[2]
		re := regexp.MustCompile(`v([\d\.]+)`)
		match := re.FindStringSubmatch(ver)
		if len(match) > 1 {
			ver = match[1]
		}
	}

	return ver
}

func GetDuration(lines []string) string {
	timestampBegin := GetDateAndTime(1, lines)
	timestampEnd := GetDateAndTime(len(lines), lines)

	diffInSeconds := timestampEnd - timestampBegin
	diffDuration := time.Duration(diffInSeconds) * time.Second

	result :=
		strconv.Itoa(int(diffDuration.Hours()/24)) + "d " +
			strconv.Itoa(int(diffDuration.Hours())%24) + "h " +
			strconv.Itoa(int(diffDuration.Minutes())%60) + "m " +
			strconv.Itoa(int(diffDuration.Seconds())%60) + "s"

	return result

}

func GetProxyPool(lines []string) string {
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.Contains(line, "Connected to pool server:") {
			startIndex := strings.Index(line, "<") + 1
			endIndex := strings.Index(line, ">")
			return line[startIndex:endIndex]
		}
	}
	return "<not connected to any pool>"
}

func GetProxyPort(lines []string) string {
	var port string

	if len(lines) >= 7 {
		line7 := lines[6]                    // Line 7
		idx := strings.LastIndex(line7, ":") // Search last ":"
		if idx != -1 {
			port = line7[idx+1:] // extract number at the EOL
		}
	}

	return port
}

func WriteLogsToFile(lines []string, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		DisplayError("Unable to create '" + fileName + "' file.")
		return fmt.Errorf("")
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		_, err := writer.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

func GetDateAndTime(lineNumber int, lines []string) int64 {
	// Retrive line
	line := lines[lineNumber-1]

	// Extract TimeStamp
	var dateStr, timeStr string
	lineParts := strings.Split(line, " ")
	if len(lineParts) > 3 {
		dateStr = lineParts[1]
		timeStr = lineParts[2]
	}

	// Convert to Unix TimeStamp
	layout := "2006/01/02 15:04:05"
	dateTimeStr := fmt.Sprintf("%s %s", dateStr, timeStr)
	dateTime, err := time.Parse(layout, dateTimeStr)
	if err != nil {
		DisplayError("Unix Timestamp error.")
	}
	TimeStamp := dateTime.Unix()

	return TimeStamp
}

func GetMiners(lines []string, message string, Details bool) int {
	Miners := 0
	for _, line := range lines {
		if strings.Contains(line, message) {
			Miners++
			if Details {
				name := strings.SplitN(line, " ", 2)[1]
				color.Green("\t%s", name)
			}
		}
	}
	return Miners
}

func DisplayMinersOnline(lines []string) {
	connectedMiners := GetMiners(lines, "Connected to pool server:", false)
	var bBreak bool

	// For each Proxy between 1 and connectedMiners
	for i := 1; i <= connectedMiners; i++ {
		bBreak = false
		MinerStudy := ""
		// On passe en revue tous les enregistrements.
		for _, line := range lines {
			if strings.Contains(line, "Proxy ["+strconv.Itoa(i)) {
				MinerStudy = line
				for _, line := range lines {
					if strings.Contains(line, "ConnID =  "+strconv.Itoa(i)) {
						bBreak = true
						break
					}
				}
				if !bBreak {
					// Remove [XDAG_PROXY] at the begining of each line
					MinerStudy := strings.SplitN(MinerStudy, " ", 2)[1]
					color.Green("\t%s", MinerStudy)
				}
			}
			if bBreak {
				break
			}
		}

	}
}

func GetShares(lines []string) int {
	for i := len(lines) - 1; i >= 0; i-- {
		line := lines[i]
		if strings.Contains(line, "shares:") {
			startIndex := strings.Index(line, "shares:") + 7
			endIndex := strings.Index(line, "(")
			sharesStr := line[startIndex:endIndex]
			shares, err := strconv.Atoi(strings.TrimSpace(sharesStr))
			if err != nil {
				DisplayError("GetShare / strconv.Atoi bad convertion.")
				return 0 // In case of error, return 0
			}
			return shares
		}
	}
	return 0 // If any line with "shares:" found, return 0
}

func DisplaySharesStats(startIndex, endIndex int, lines []string) {
	timestampBegin := GetDateAndTime(startIndex, lines)
	timestampEnd := GetDateAndTime(endIndex, lines)

	diffInSeconds := timestampEnd - timestampBegin
	diffDuration := time.Duration(diffInSeconds) * time.Second

	shares := GetShares(lines)

	DisplayItem("Shares founds", strconv.Itoa(shares))
	DisplayItem("Shares / minute", strconv.FormatFloat(float64(shares)/diffDuration.Minutes(), 'f', 2, 64))
	DisplayItem("Shares / hours", strconv.FormatFloat(float64(shares)/diffDuration.Hours(), 'f', 2, 64))

	if shares > 1 {
		DisplayItem("New share each", strconv.FormatFloat(float64(diffInSeconds)/float64(shares-1), 'f', 2, 64)+" seconds")
	} else {
		fmt.Println("  -> Cannot compute seconds per share.")
	}
}

func main() {
	// Display Program Title
	DisplayProgramTitle()

	// Default Args
	argMinersDetails := false   // -d
	argExportFilterLog := false // -e
	// Default log filename
	fileName := "proxy.log"

	args := os.Args[1:]

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--f":
			if i+1 < len(args) && args[i+1] != "" && !strings.HasPrefix(args[i+1], "-") {
				fileName = args[i+1][0:]
				i++
			} else {
				DisplayError("--f argument requires a filename.")
				return
			}
		case "--d":
			argMinersDetails = true
		case "--e":
			argExportFilterLog = true
		case "--h", "--help":
			DisplaySubTitle("How to use " + AppName)
			fmt.Println("\nUsage : \t" + AppExe + " [OPTIONS]")
			fmt.Println("\nOptions :")
			DisplayItem("--d, --display", "Display more details about miners.")
			DisplayItem("--e, --export", "Export filtered log in 'filtered_<yourfile.log>'.")
			DisplayItem("--f, --file", "Define log file to analyze (proxy.log by default).")
			DisplayItem("--h, --help", "Display this help.")
			fmt.Println("\nUse cases :")
			fmt.Println("  " + AppExe + "\t\t: analyse <proxy.log>.")
			fmt.Println("  " + AppExe + " --f mylog.txt\t: analyse <mylog.txt>.")
			fmt.Println("  " + AppExe + " --display \t: analyse <proxy.log> + display more informations miners.")
			fmt.Println("  " + AppExe + " --e > x2x.txt \t: analyse <proxy.log> + export filtered logs in <filtered_proxy.log> + store results in <x2x.txt>.")
			return
		default:
			DisplayError("unrecognized argument '" + args[i] + "'.")
			return
		}
	}

	// Load log to memory
	DisplaySubTitle("Loading " + fileName)
	lines, err := GetLogs(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Relevant records
	DisplayItem("Relevant records", strconv.Itoa(len(lines)))

	// Export Filtered logfile
	if argExportFilterLog && len(lines) > 0 {
		err := WriteLogsToFile(lines, "filtered_"+fileName)
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// Call GetProxyVer to retreive Proxy Version and display it
	DisplayItem("Proxy version", GetProxyVer(lines))

	// Call DisplayDuration to display Proxy Duration
	if len(lines) > 0 {
		DisplayItem("Proxy running for", GetDuration(lines))
		//DisplayDuration(1, len(lines), lines)
	}

	DisplaySubTitle("Proxy communications")

	// Retreive pool proxy
	DisplayItem("Connected to pool", GetProxyPool(lines))

	// Retrieve proxy listening port
	DisplayItem("Listening port", GetProxyPort(lines))

	DisplaySubTitle("Miners connecting summary")

	// Get connected miners without display them
	connectedMiners := GetMiners(lines, "Connected to pool server:", false)
	// Get stopped miners without display them
	stoppedMiners := GetMiners(lines, "Conn Stoped()...ConnID", false)
	// Get shutdowned miners without display them
	shutdownedMiners := GetMiners(lines, "Conn Closed() ...ConnID", false)

	// Display onlineMiners count
	DisplayItem("Online\t", strconv.Itoa(connectedMiners-(stoppedMiners+shutdownedMiners)))
	//if argMinersDetails {
	DisplayMinersOnline(lines)
	//}

	// Dispaly connectedMiners Count
	DisplayItem("Connected 2 proxy", strconv.Itoa(connectedMiners))
	// Display connectedMiner List
	if argMinersDetails {
		connectedMiners = GetMiners(lines, "Connected to pool server:", true)
	}

	// Dispaly stoppedMiners Count
	DisplayItem("Shutdown by pool", strconv.Itoa(stoppedMiners))
	// Display stoppedMiner List
	if argMinersDetails {
		stoppedMiners = GetMiners(lines, "Conn Stoped()...ConnID", true)
	}

	// Dispaly shutdownedMiners Count
	DisplayItem("Shutdown by miner", strconv.Itoa(shutdownedMiners))
	// Dispaly shutdownedMiners List
	if argMinersDetails {
		shutdownedMiners = GetMiners(lines, "Conn Closed() ...ConnID", true)
	}

	DisplaySubTitle("Shares Statistics")

	// Call DisplaySharesStats to display Proxy Duration
	if len(lines) > 0 {
		DisplaySharesStats(1, len(lines), lines)
	}

}
