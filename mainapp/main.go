package main

import (
	"fmt"
)

func main() {
	fmt.Println("Here we are")

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
