package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/will666/golarge/helper"
	"github.com/will666/golarge/types"
)

const BIG_FILE_SIZE int64 = 1_000_000_000

var logging bool = false
var logFile string

var entries []types.List
var total int = 0

func main() {
	var fileOutput string
	var jsonOutput bool
	var help bool
	var path string

	flag.StringVar(&fileOutput, "o", "list.txt", "Write output to text file")
	flag.BoolVar(&jsonOutput, "j", false, "Ouput list to JSON file")
	flag.BoolVar(&help, "help", false, "Help")
	flag.Parse()
	args := flag.Args()

	if help {
		flag.PrintDefaults()
	}

	if flag.NArg() == 1 {
		path = string(args[0])

		if flag.NFlag() >= 1 && fileOutput != "" {
			logging = true
			logFile = fileOutput
		}
		log.Printf("-- Searching large files in %s --", helper.Colorize(path, "cyan"))
		listFiles(path)
		if jsonOutput {
			saveToJson(entries)
		}
		log.Printf("-- Found %d files of size around 1GB --", total)
		if logging {
			log.Printf("-- List generated: %s --", helper.Colorize(logFile, "cyan"))
		}
	} else {
		fmt.Println("\nUsage: golarge [OPTIONS] PATH")
		fmt.Println("\nUtil to find files around 1GB of size from given directory path")
		fmt.Println("\nOptions:")
		fmt.Println("  -o string   Full path of file to save large file list to (default: list.txt)")
		fmt.Println("  -j          Enable export to JSON file")
		fmt.Println("\nExamples:")
		fmt.Println("  golarge /foo/bar")
		fmt.Println("  golarge -o list.txt /foo/bar")
		fmt.Println("  golarge -o list.txt -j /foo/bar")
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
					f := fmt.Sprintf("%s/%s => %dMiB", path, name, size/(1024*1024))
					log.Println(helper.Colorize(f, "green"))
					total++
					entries = append(entries, types.List{Name: name, BasePath: path, FullPath: fmt.Sprintf("%s/%s", path, name), Size: size, Type: filepath.Ext(name)})
					if logging {
						saveToFile(logFile, f)
					}
				}
			} else {
				log.Println(helper.Colorize(err.Error(), "yellow"))
			}
		}
	} else {
		log.Println(helper.Colorize(err.Error(), "yellow"))
	}
}

func saveToFile(dst string, file string) {
	f, err := os.OpenFile(dst, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(helper.Colorize(err.Error(), "red"))
	}
	defer f.Close()

	if _, err := os.Stat(dst); err == nil && total == 1 {
		f.Truncate(0)
	} else if err != nil {
		log.Println(helper.Colorize(err.Error(), "red"))
	}

	if _, err := f.WriteString(fmt.Sprintf("%s\n", file)); err != nil {
		log.Println(helper.Colorize(err.Error(), "red"))
	}
}

func saveToJson(list []types.List) {
	if content, err := json.Marshal(list); err == nil {
		fileName := helper.ExtLess(logFile)
		os.WriteFile(fmt.Sprintf("%s.json", fileName), content, 0644)
	} else {
		log.Fatalln(err.Error())
	}
}
