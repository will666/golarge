package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const BIG_FILE_SIZE int64 = 1_000_000_000

var FILE_LIST string

var total int = 0

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
			FILE_LIST = output
		} else {
			FILE_LIST = "list.txt"
		}
		log.Printf("-- Searching large files in %s --", path)
		listFiles(path)
		log.Printf("-- Found %d files of size around 1GB --", total)
		log.Printf("-- List generated: %s --", FILE_LIST)
	} else {
		fmt.Println("Usage: golarge PATH")
		fmt.Println("An util to list large files from a given directory path")
		fmt.Println("Example: golarge ~/")
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
					res := fmt.Sprintf("%s/%s => %dMB", path, name, size/1024/1024)
					log.Println(res)
					total++
					saveToFile(FILE_LIST, res)
				}
			}
		}
	}
}

func saveToFile(dst string, file string) {
	f, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if _, err := os.Stat(dst); err == nil && total == 1 {
		f.Truncate(0)
	} else {
		log.Fatal(err)
	}
	if _, err := f.WriteString(fmt.Sprintf("%s\n", file)); err != nil {
		log.Fatal(err)
	}
}
