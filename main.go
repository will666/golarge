package main

import (
	"fmt"
	"log"
	"os"
)

const BIG_FILE_SIZE int64 = 1000_000_000
const FILE_LIST = "./list.txt"

var total int = 0

func main() {
	args := os.Args
	// fmt.Println(args)

	if len(args) >= 2 {
		path := args[1]

		log.Println("Searching large files on", path)
		fmt.Println("")

		f, err := os.OpenFile(FILE_LIST, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		if _, err := os.Stat(FILE_LIST); err == nil {
			f.Truncate(0)
		} else {
			log.Fatal(err)
		}

		listFiles(path)

		fmt.Println("")
		log.Println("[Found", total, "files over 1GB]")
	} else {
		fmt.Println("")
		fmt.Println("Usage: golarge PATH")
		fmt.Println("")
		fmt.Println("An util to list large files from a given directory path")
		fmt.Println("")
		fmt.Println("Example: golarge /tmp")
		fmt.Println("")
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
					fmt.Printf("%s/%s => %dMB\n", path, name, size/1024/1024)
					total++
					saveToFile(FILE_LIST, fmt.Sprintf("%s/%s", path, name))
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
	if _, err := f.WriteString(fmt.Sprintf("%s\n", file)); err != nil {
		log.Fatal(err)
	}
}
