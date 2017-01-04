package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: ./%s [<JSON file> [<JSON file>]]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	host := flag.String("n", "localhost", "hostname")
	port := flag.Int("p", 8000, "listen on port")
	namespacePattern := flag.String("namespace", "^[a-z0-9\\._-]*$", "pattern for namespaces/projects")
	namePattern := flag.String("name", "^[a-z0-9\\._-]+$", "pattern for names")
	containerPattern := flag.String("container", "^[a-z0-9\\._-]+$", "pattern for containers")
	envPattern := flag.String("env", "^[A-Z0-9_-]+$", "pattern for environment variables")
	flag.Parse()
	args := flag.Args()

	//patterns are ignored in server mode - specify as follows:
	//{ customNamespacePattern="...", customNamePattern="..."", ... }
	if len(args) == 0 {
		serve(*host, *port)
		return
	}

	for _, arg := range args {
		start := time.Now()
		msg, code := processFile(arg, *namespacePattern, *namePattern, *containerPattern, *envPattern)
		secs := time.Since(start).Seconds()

		if code > 0 {
			fmt.Printf("%s: %s (%.2fs)\n", arg, msg, secs)
			return
		}
		fmt.Printf("%s\n", msg)
	}
}

func processFile(path, namespacePattern, namePattern, containerPattern, envPattern string) (string, int) {
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return fmt.Sprintf("can't read %s", path), 1
	}

	combinedResultMap, err := processBytes(bytes, namespacePattern, namePattern, containerPattern, envPattern)

	if err != nil {
		return fmt.Sprintf("can't process bytes %s", path), 1
	}

	json, err := json.MarshalIndent(combinedResultMap, "", "  ")

	if err != nil {
		return fmt.Sprintf("can't marshall JSON %v", combinedResultMap), 1
	}

	return string(json), 0
}
