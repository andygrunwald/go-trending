package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

// A small helper tool to update the local test data for unit tests.
//
// This tool will download HTML from GitHub.com and writes
// it into the local testdata folder.
// This folder is used as fixtures in our unit tests.
//
// Attention: If remote HTML structure changes, there is a high possibility
// that unit tests need to be adjusted.

const (
	BASE_REPOSITORY_URL = "https://github.com/trending"
	BASE_DEVELOPERS_URL = "https://github.com/trending/developers"

	DIR_TESTDATA    = "../testdata"
	FILE_REPOSITORY = "github.com_trending.html"
	FILE_DEVELOPERS = "github.com_trending_developers.html"
)

func main() {
	contentToProcess := map[string]string{
		BASE_REPOSITORY_URL: DIR_TESTDATA + string(os.PathSeparator) + FILE_REPOSITORY,
		BASE_DEVELOPERS_URL: DIR_TESTDATA + string(os.PathSeparator) + FILE_DEVELOPERS,
	}

	log.Println("Starting to update package test data for this package by downloading HTML")
	log.Printf("from GitHub.com and writing this into local files inside the %s directory.", DIR_TESTDATA)
	log.Println("")
	log.Println("The possibility that the unit tests are failing afterwards, due to HTML structure changes")
	log.Println("by GitHub is high. Please execute unit tests after updating the test data and adjust ")
	log.Println("this library accordingly to match the new structure.")
	log.Println("")

	for u, f := range contentToProcess {
		log.Printf("Calling URL %s", u)
		resp, err := http.Get(u)
		if err != nil {
			log.Fatalf("%s", err)
		}
		defer resp.Body.Close()

		log.Printf("Opening file %s", f)
		handle, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE, 0655)
		if err != nil {
			log.Fatalf("%s", err)
		}
		defer handle.Close()

		log.Printf("Writing content of URL %s into file %s", u, f)
		n, err := io.Copy(handle, resp.Body)
		if err != nil {
			log.Fatalf("%s", err)
		}
		log.Printf("Wrote %d bytes into file %s", n, f)
	}

	log.Println("")
	log.Println("Update package test data was successful")
}
