//Package ozip will look at a set of Oracle AWR text files and create one zip file per database, suitable
// for submission to Mitrend.

package main

import (
	"archive/zip"
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type hostMap struct {
	hostName     string
	databaseName string
}

func main() {
	dbAndHostPattern := regexp.MustCompile("DB Name.*Host")
	dbNamePattern := regexp.MustCompile("DB Name")
	hostNamePattern := regexp.MustCompile("Host Name")
	dashPattern := regexp.MustCompile("----")
	txtPattern := regexp.MustCompile("\\.txt$")
	files, err := ioutil.ReadDir("./")
	check(err)
	host2DB := make(map[hostMap]bool)
	db2File := make(map[string][]string)
	for _, f := range files {
		fileName := f.Name()
		if txtPattern.MatchString(fileName) {
			file, err := os.Open(fileName)
			check(err)
			scanner := bufio.NewScanner(file)
			var lineStatus string
			lineStatus = ""
			hostName := ""
			databaseName := ""
			for scanner.Scan() {
				resultText := scanner.Text()
				if dbAndHostPattern.MatchString(resultText) {
					lineStatus = "dbAndHost"
				} else if dbNamePattern.MatchString(resultText) {
					lineStatus = "db"
				} else if hostNamePattern.MatchString(resultText) {
					if lineStatus == "dbDataPulled" {
						lineStatus = "host"
					}
				} else if dashPattern.MatchString(resultText) {
					if lineStatus == "dbAndHost" {
						lineStatus = "dbAndHostdata"
					} else if lineStatus == "db" {
						lineStatus = "dbData"
					} else if lineStatus == "host" {
						lineStatus = "hostData"
					}
				} else if lineStatus == "dbAndHostdata" {
					lineData := strings.Fields(resultText)
					databaseName = lineData[0]
					hostName = lineData[6]
					lineStatus = ""
					break
				} else if lineStatus == "dbData" {
					lineData := strings.Fields(resultText)
					databaseName = lineData[0]
					lineStatus = "dbDataPulled"
				} else if lineStatus == "hostData" {
					lineData := strings.Fields(resultText)
					hostName = lineData[0]
					lineStatus = ""
					break
				}
			}

			db2File[databaseName] = append(db2File[databaseName], fileName)

			if !host2DB[hostMap{hostName, databaseName}] {
				host2DB[hostMap{hostName, databaseName}] = true
			}
		}
	}

	for database := range db2File {
		zipFileName := database + ".zip"
		fmt.Println("Creating " + zipFileName + "...")
		zipFile, err := os.Create(zipFileName)
		check(err)
		w := zip.NewWriter(zipFile)

		for _, file := range db2File[database] {
			f, err := w.Create(file)
			dat, err := ioutil.ReadFile(file)
			check(err)
			_, err = f.Write([]byte(dat))
			check(err)
		}
		err = w.Close()
		check(err)
	}
	//for hostMap := range host2DB {
	//	fmt.Printf("%s,%s\n", hostMap.hostName, hostMap.databaseName)s
	//}
}
