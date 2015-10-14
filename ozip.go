//Package ozip will look at a set of Oracle AWR text files and create one zip file per server.
//The package assumes that the format of the filename is SERVER.restOfTheFileName.txt

package main

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func main() {
	files, err := ioutil.ReadDir("./")
	check(err)
	oracleFiles := make(map[string][]string)
	for _, f := range files {
		fileName := f.Name()
		fileNameFields := strings.Split(fileName, ".")
		server := fileNameFields[0]
		oracleFiles[server] = append(oracleFiles[server], fileName)
	}

	for i, _ := range oracleFiles {
		server := i
		zipFileName := server + ".zip"
		fmt.Println("Creating " + zipFileName + "...")
		zipFile, err := os.Create(zipFileName)
		check(err)
		w := zip.NewWriter(zipFile)

		for _, file := range oracleFiles[server] {
			f, err := w.Create(file)
			dat, err := ioutil.ReadFile(file)
			check(err)
			_, err = f.Write([]byte(dat))
			check(err)
		}
		err = w.Close()
		check(err)
	}
}
