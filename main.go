package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fatih/color"
)

const BIG_FILE_SIZE int64 = 1_000_000_000

var logFile string
var total int = 0
var blue = color.New(color.FgCyan, color.Bold)
var green = color.New(color.FgGreen, color.Bold)

func main() {
	var output string
	var help bool
	flag.StringVar(&output, "o", "list.txt", "Write output to file")
	flag.BoolVar(&help, "help", false, "Help")
	flag.Parse()
	args := flag.Args()

	if help {
		flag.PrintDefaults()
	}

	var path string

	if flag.NArg() == 1 {
		path = string(args[0])

		if flag.NFlag() == 1 && output != "" {
			logFile = output
		} else {
			logFile = "list.txt"
		}
		log.Printf("-- Searching large files in %s --", blue.Sprintf(path))
		listFiles(path)
		log.Printf("-- Found %d files of size around 1GB --", total)
		log.Printf("-- List generated: %s --", blue.Sprintf(logFile))
	} else {
		fmt.Println("\nMissing argument: PATH")
		fmt.Println("Provide a directory as argument: /tmp")
		fmt.Println("\nUsage: golarge PATH")
		fmt.Println("An util to list large files of a given directory path")
		fmt.Println("Example: ./golarge ~/")
		os.Exit(0)
	}
}

func listFiles(path string) {
	if files, err := os.ReadDir(path); err == nil {
		for _, v := range files {
			if v.IsDir() {
				listFiles(fmt.Sprintf("%s/%s", path, v.Name()))
			}
			if info, err := v.Info(); err == nil {
				name := info.Name()
				size := info.Size()
				if size >= BIG_FILE_SIZE {
					f := fmt.Sprintf("%s/%s => %dMB", path, name, size/1024/1024)
					log.Println(green.Sprintf(f))
					total++
					saveToFile(logFile, f)
				}
			}
		}
	} else {
		log.Fatal(err.Error())
	}
}

func saveToFile(dst string, file string) {
	f, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer f.Close()

	if _, err := os.Stat(dst); err == nil && total == 1 {
		f.Truncate(0)
	} else {
		log.Fatal(err.Error())
	}
	if _, err := f.WriteString(fmt.Sprintf("%s\n", file)); err != nil {
		log.Fatal(err.Error())
	}
}
