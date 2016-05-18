package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	dbServerUrl		string
	dbName			string
	dbUserName		string
	dbUserPassw		string

//	ftpLocationURL	string
//	ftpUserName		string
//	ftpUserPassw	string
//	ftpSubfolder	string

	inputFolder		string
	tmpFolderName	string
)

func init() {
	flagDbServerUrl		:= flag.String("dbs", "", "`DB Server name/IP`")
	flagDbName			:= flag.String("dbn", "", "`Database Name`")
	flagDbUserName		:= flag.String("dbu", "", "`Database User Name`")
	flagDbUserPassw		:= flag.String("dbp", "", "`Database User Password`")

//	flagFtpLocationURL	:= flag.String("fur", "", "`FTP URL`")
//	flagFtpUserName		:= flag.String("fun", "", "`FTP User Name`")
//	flagFtpUserPassw	:= flag.String("fup", "", "`FTP User Password`")
//	flagFtpSubfolder	:= flag.String("fus", "", "`FTP Subfolder`")

	flagInputFolder		:= flag.String("i", "", "`Input Folder Name`")
	flagTmpFolderName	:= flag.String("tmp", "tmp", "`Tmp folder location`")

	flag.Parse()

	if flag.Parsed() {
		dbServerUrl		= *flagDbServerUrl
		dbName			= *flagDbName
		dbUserName		= *flagDbUserName
		dbUserPassw		= *flagDbUserPassw

//		ftpLocationURL	= *flagFtpLocationURL
//		ftpUserName		= *flagFtpUserName
//		ftpUserPassw	= *flagFtpUserPassw
//		ftpSubfolder	= *flagFtpSubfolder

		inputFolder 	= *flagInputFolder
		tmpFolderName	= *flagTmpFolderName

		} else {
			flag.Usage()
			os.Exit(-1)
		}

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
	fmt.Println("tmpFolderName\t", tmpFolderName)
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

	//		2b. For each file:
	//			2b-1: Parse events
	//			2b-2: Save to DB
	//		2c. Delete all unpacked files
	// 3. Print summary (???):
	//		3a. Number of zip files processed
	//		3b. Number of raw files processed
	//		3c.	Number of STB's
	//		3d. Number of Parsed Events
}
