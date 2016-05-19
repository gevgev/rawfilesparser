package main

import (
	"bufio"
	"database/sql"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/denisenkom/go-mssqldb"
)

var (
	dbServerUrl string
	dbName      string
	dbUserName  string
	dbUserPassw string

	connString string

	//	ftpLocationURL	string
	//	ftpUserName		string
	//	ftpUserPassw	string
	//	ftpSubfolder	string

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
	flagDbServerUrl := flag.String("dbs", "", "`DB Server name/IP`")
	flagDbDatabase := flag.String("dbn", "Clickstream", "`Database name`")
	flagDbUserName := flag.String("dbu", "", "`Database User Name`")
	flagDbUserPassw := flag.String("dbp", "", "`Database User Password`")

	//	flagFtpLocationURL	:= flag.String("fur", "", "`FTP URL`")
	//	flagFtpUserName		:= flag.String("fun", "", "`FTP User Name`")
	//	flagFtpUserPassw	:= flag.String("fup", "", "`FTP User Password`")
	//	flagFtpSubfolder	:= flag.String("fus", "", "`FTP Subfolder`")

	flagInputFolder := flag.String("i", "", "`Input Folder Name`")

	flag.Parse()

	if flag.Parsed() {
		dbServerUrl = *flagDbServerUrl
		dbName = *flagDbDatabase
		dbUserName = *flagDbUserName
		dbUserPassw = *flagDbUserPassw

		//		ftpLocationURL	= *flagFtpLocationURL
		//		ftpUserName		= *flagFtpUserName
		//		ftpUserPassw	= *flagFtpUserPassw
		//		ftpSubfolder	= *flagFtpSubfolder

		inputFolder = *flagInputFolder

		connString = "server=" + dbServerUrl + ";user id=" + dbUserName + ";password=" + dbUserPassw + ";database=" + dbName

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
	if receivedIndex > -1 {
		received = tokens[receivedIndex]
	} else {
		received = "2016-05-14_00:35:44"
	}

	return received, deviceId, clickString, nil
}

func printParams() {
	fmt.Println("dbServerUrl\t", dbServerUrl)
	fmt.Println("dbName\t\t", dbName)
	fmt.Println("dbUserName\t", dbUserName)
	fmt.Println("dbUserPassw\t", dbUserPassw)

	//	fmt.Println("ftpLocationURL\t", ftpLocationURL)
	//	fmt.Println("ftpUserName\t", ftpUserName)
	//	fmt.Println("ftpUserPassw\t", ftpUserPassw)
	//	fmt.Println("ftpSubfolder\t", ftpSubfolder)

	fmt.Println("inputFolderName\t", inputFolder)
}

func main() {
	printParams()

	// Input Parameters:
	//	1. DB Server URL
	//	2. DB Name
	//	3. DB user name
	//	4. DB user password
	//	5. FTP Location URL
	//	6. FTP user name
	//	7. FTP user password
	//	8. FTP Subfolder to look into for ZIP files
	//	9. TMP Folder location = default assigned to /tmp/?

	// 1. Get the list of zip files from FTP
	// 		1a. --- Not necessary: --- Sort the list alphabetically
	// 2. For each zip file
	//		2a. Unpack to tmp location
	//		2b. --- Not necessary: --- Sort the list alphabetically
	//			Get the list of files from the provided location

	//		2b. For each file:
	//			2b-1: Parse events
	//			2b-2: Save to DB
	//		2c. Delete all unpacked files
	// 3. Print summary (???):
	//		3a. Number of zip files processed
	//		3b. Number of raw files processed
	//		3c.	Number of STB's
	//		3d. Number of Parsed Events

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
				}
				fmt.Println(event)
				eventsCollection = append(eventsCollection, event)
			}

			msoName := getMsoName(gfile)
			fmt.Println("About to save...")
			saveToDb(eventsCollection, msoName)
		}(gfile)
	}

	// waiting for all goroutines to end
	fmt.Println("Waiting for all goroutines to complete the work")

	for i := 0; i < cap(sem); i++ {
		sem <- true
	}

	fmt.Printf("Processed %d files\n", len(files))

}

func getMsoName(fileName string) string {
	return "'mso'"
}

func saveToDb(eventsCollection []Event, msoName string) {
	db, err := sql.Open("mssql", connString)
	if err != nil {
		fmt.Println("Cannot connect: ", err.Error())
		return
	}

	defer db.Close()

	cmd := `INSERT INTO [Clickstream].[dbo].[clickstreamEventsLog] 
           ([timestamp] 
           ,[received] 
           ,[deviceId] 
           ,[eventCode]
           ,[msoName])
     VALUES `

	fmt.Println(cmd + getValues(eventsCollection, msoName))
	//	db.Exec(cmd + getValues(eventsCollection, msoName))
}

func timeFormat(timestamp time.Time) string {
	str := timestamp.Format(time.RFC3339)
	return str[0:strings.LastIndex(str, "-")]
}

func getValues(eventsCollection []Event, msoName string) string {

	var valuesStr string

	for _, event := range eventsCollection {
		single := "(" +
			"'" + timeFormat(event.Timestamp) + "'," +
			"'" + event.Received + "'," +
			"'" + event.DeviceId + "'," +
			"'" + event.Command + "'," +
			msoName +
			")"
		valuesStr = valuesStr + ", " + single
	}
	/*
	   (<timestamp, datetime,>
	   ,<received, datetime,>
	   ,<deviceId, nvarchar(50),>
	   ,<eventCode, nvarchar(max),>
	   ,<msoName, nvarchar(50),>)" */
	return valuesStr[1:len(valuesStr)]
}

type Event struct {
	Command   string
	Timestamp time.Time
	Received  string
	DeviceId  string
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
