package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	inputFolder string
)

type Command string

const (
	R_AD           Command = "41" // A
	R_BtnCnfg      Command = "42" // B
	R_ChanVrb      Command = "43" // C
	R_PROGRAMEVENT Command = "45" // E
	R_VODCat       Command = "47" // G
	R_HIGHLIGHT    Command = "48" // H
	R_INFO         Command = "49" // I
	R_KEY          Command = "4B" // K
	R_MISSING      Command = "4D" // M
	R_OPTION       Command = "4F" // O
	R_PULSE        Command = "50" // P
	R_RESET        Command = "52" // R
	R_STATE        Command = "53" // S
	R_TURBO        Command = "54" // T
	R_UNIT         Command = "55" // U
	R_VIDEO        Command = "56" // V
)

const (
	version      = "0.9"
	rawExt       = "raw"
	UTC_GPS_Diff = 315964800
)

type FileType int

const (
	FT_WRONG FileType = iota
	FT_RAW
)

func CheckCommand(clickString string) Command {
	return Command(clickString[0:2])
}

func init() {

	flagInputFolder := flag.String("i", "", "`Input Folder Name`")

	flag.Parse()

	if flag.Parsed() {
		inputFolder = *flagInputFolder

	} else {
		flag.Usage()
		os.Exit(-1)
	}

}

func preParseLine(line string) (received string, deviceId string, clickString string, err error) {
	var receivedIndex, deviceIndex, eventIndex int

	tokens := strings.Split(line, " ")
	parts := len(tokens)

	if parts < 2 || parts > 3 {
		return "", "", "", errors.New("Wrong line format: " + line + " - ")
	}

	switch parts {
	case 2:
		receivedIndex = -1
		deviceIndex = 0
		eventIndex = 1
	case 3:
		receivedIndex = 0
		deviceIndex = 1
		eventIndex = 2
	}

	deviceId, clickString = tokens[deviceIndex], tokens[eventIndex]

	if len(clickString) < 22 {
		return "", "", "", errors.New("Wrong line format: " + line + " - ")
	}

	if receivedIndex > -1 {
		received = tokens[receivedIndex]
	} else {
		received = "1900-01-01 00:00:00"
	}

	return received, deviceId, clickString, nil
}

func printParams() {
	fmt.Println("inputFolderName\t", inputFolder)
}

func main() {
	printParams()

	// This is our semaphore/pool
	concurrency := 5
	sem := make(chan bool, concurrency)

	files := getFilesToProcess()
	fmt.Println("Files to be processed: ", files)

	for _, gfile := range files {
		// if we still have available goroutine in the pool (out of concurrency )
		fmt.Println("Processing file: ", gfile)
		sem <- true

		// fire one file to be processed in a goroutine
		go func(fileName string) {
			// Signal end of processing at the end
			defer func() { <-sem }()
			eventsCollection := []Event{}

			file, err := os.Open(fileName)
			if err != nil {
				fmt.Println("Error opening file: ", err)
				return
			}
			defer file.Close()

			msoName := getMsoName(fileName)
			scanner := bufio.NewScanner(file)

			for scanner.Scan() {
				line := scanner.Text()
				received, deviceId, clickString, err := preParseLine(line)

				if err != nil {
					fmt.Println(err, fileName)
					continue
				}

				command := lookUpEventName(clickString[0:2])
				timestamp := convertToTime(clickString[2:10])
				receivedTS := received

				event := Event{
					command,
					timestamp,
					strings.Replace(receivedTS, "_", " ", -1),
					deviceId,
					msoName,
				}
				// fmt.Println(event)
				eventsCollection = append(eventsCollection, event)
			}

			fmt.Println("About to save...")
			fileNameToSave := formatFileNameToSave(fileName)

			processJson(fileNameToSave, eventsCollection)
		}(gfile)
	}

	// waiting for all goroutines to end
	fmt.Println("Waiting for all goroutines to complete the work")

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	fmt.Printf("Processed %d files\n", len(files))

}

// Format output file name
func formatFileNameToSave(currentFileName string) string {
	return validateOutFileName(currentFileName[:len(currentFileName)-len(".raw")])
}

func validateOutFileName(fileName string) string {
	// Check if it has extension
	// If not, add the default extension
	ext := filepath.Ext(fileName)
	if ext != "" {
		if isRawFile(fileName) {
			fileName = fileName[:len(fileName)-len(ext)]
			fileName = addProperExtension(fileName)
		}
	} else if ext == "" {
		fileName = addProperExtension(fileName)
	}

	return fileName
}

func getMsoName(fileName string) string {
	mso := fileName[strings.LastIndex(fileName, "_")+1 : strings.LastIndex(fileName, ".")]
	return mso
}

func timeFormat(timestamp time.Time) string {
	str := timestamp.Format(time.RFC3339)
	return str[0:strings.LastIndex(str, "-")]
}

type Event struct {
	Command   string
	Timestamp time.Time
	Received  string
	DeviceId  string
	Mso       string
}

// Get the list of files to process in the target folder
func getFilesToProcess() []string {
	fileList := []string{}

	// We have working directory - takes over single file name, if both provided
	err := filepath.Walk(inputFolder, func(path string, f os.FileInfo, _ error) error {
		if isRawFile(path) {
			fileList = append(fileList, path)
		}
		return nil
	})

	if err != nil {
		fmt.Println("Error getting files list: ", err)
		os.Exit(-1)
	}

	return fileList
}

func isRawFile(fileName string) bool {
	return filepath.Ext(fileName) == "."+rawExt
}

func convertToTime(timestampS string) time.Time {
	timestamp, err := strconv.ParseInt(timestampS, 16, 64)
	//fmt.Println(timestampS, timestamp)
	if err == nil {
		timestamp += UTC_GPS_Diff
		//fmt.Println(timestampS, timestamp)

		t := time.Unix(timestamp, 0)
		return t
	}
	//else {
	//fmt.Println("Error:", err)
	//}
	return time.Time{}
}

func lookUpEventName(code string) string {
	return EventCodes[code]
}

func addProperExtension(fileName string) string {
	fileName = fileName + ".json"
	return fileName
}

func processJson(filename string, eventsCollection []Event) {
	jsonString, err := generateJson(eventsCollection)
	if err == nil {
		err = saveJsonToFile(filename, jsonString)
		if err != nil {
			fmt.Println("Error writing Json file:", err)
		}
	} else {
		fmt.Println(err)
	}

}
