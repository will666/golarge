package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"sync"
)

// const BIG_FILE_SIZE int64 = 1_000_000_000
const BIG_FILE_SIZE int64 = (1024 * 1024 * 1024)
const DEBUG = false

func main() {
	if DEBUG {
		// stat()
		PrintMemUsage()
	}

	var fileOutput string
	var jsonOutput bool
	var help bool

	flag.StringVar(&fileOutput, "o", "list.txt", "Write output to text file")
	flag.StringVar(&fileOutput, "output", "list.txt", "Write output to text file")
	flag.BoolVar(&jsonOutput, "j", false, "Ouput list to JSON file")
	flag.BoolVar(&jsonOutput, "json", false, "Ouput list to JSON file")
	flag.BoolVar(&help, "help", false, "Help")
	flag.Parse()
	args := flag.Args()

	if help {
		flag.PrintDefaults()
	}

	if flag.NArg() == 1 {
		dirPath := string(args[0])
		logging := false
		var logFile string

		if flag.NFlag() >= 1 && fileOutput != "" {
			logging = true
			logFile = fileOutput
		}
		log.Printf("-- Searching large files in %s --", Colorize(dirPath, "cyan"))
		fl := newFileList()
		if err := fl.listFiles(dirPath, logging, logFile); err != nil {
			log.Fatal("%w", err)
		}
		if jsonOutput {
			jsonFile := newJsontFile(logFile, fl.data)
			jsonFile.saveToJson()
		}
		log.Printf("-- Found %d files with size bigger than 1GiB --", fl.count)
		if logging {
			log.Printf("-- List generated: %s --", Colorize(logFile, "cyan"))
			if jsonOutput {
				log.Printf("-- List generated: %s --", Colorize(ExtLess(logFile)+".json", "cyan"))
			}
		}
	} else {
		fmt.Println("\nUsage: golarge [OPTIONS] PATH")
		fmt.Println("\nUtil to find files around 1GiB of size from given directory path")
		fmt.Println("\nOptions:")
		fmt.Println("  -o, --output string   	 Output path (default: list.txt)")
		fmt.Println("  -j, --json      			 Enable export to JSON file")
		fmt.Println("\nExamples:")
		fmt.Println("  golarge /foo/bar")
		fmt.Println("  golarge -o list.txt /foo/bar")
		fmt.Println("  golarge --output list.txt --json /foo/bar")
		fmt.Println("")
		os.Exit(0)
	}

	if DEBUG {
		PrintMemUsage()
	}
}

func (fl *fileList) listFiles(filePath string, logging bool, logFile string) error {
	if files, err := os.ReadDir(filePath); err == nil {
		wg := sync.WaitGroup{}
		for _, v := range files {
			if v.IsDir() {
				fl.listFiles(path.Join(filePath, v.Name()), logging, logFile)
			} else {
				wg.Add(1)
				go func(v fs.DirEntry) {
					if info, err := v.Info(); err == nil {
						fileSize := info.Size()
						if fileSize >= BIG_FILE_SIZE {
							fileName := info.Name()
							data := fmt.Sprintf("%s => %dMiB", path.Join(filePath, fileName), fileSize/(1024*1024))
							log.Println(Colorize(data, "green"))
							nl := newList(fileName, filePath, path.Join(filePath, fileName), fileSize, filepath.Ext(fileName))
							fl.Lock()
							fl.data = append(fl.data, nl)
							fl.count++
							fl.Unlock()
							if logging {
								txtFile := newTxtFile(logFile, data)
								if err := txtFile.saveToFile(fl.count); err != nil {
									log.Fatal("%w", err)
								}
							}
						}
					} else {
						log.Println(Colorize(err.Error(), "yellow"))
					}
					wg.Done()
				}(v)
			}
		}
		wg.Wait()
	} else {
		log.Println(Colorize(err.Error(), "yellow"))
	}

	return nil
}

func (tf *txtFile) saveToFile(entryCount int) error {
	tf.Lock()
	defer tf.Unlock()

	f, err := os.OpenFile(tf.fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if err := fmt.Errorf("%w", err); err != nil {
			return err
		}
	}
	defer f.Close()

	if _, err := os.Stat(tf.fileName); err == nil && entryCount == 1 {
		f.Truncate(0)
	} else if err != nil {
		if err := fmt.Errorf("%w", err); err != nil {
			return err
		}
	}

	if _, err := f.WriteString(fmt.Sprintf("%s\n", tf.data)); err != nil {
		if err := fmt.Errorf("%w", err); err != nil {
			return err
		}
	}

	return nil
}

func (js jsonFile) saveToJson() error {
	if content, err := json.Marshal(js.data); err == nil {
		fileName := ExtLess(js.fileName)
		err := os.WriteFile(fmt.Sprintf("%s.json", fileName), content, 0644)
		if err != nil {
			if err := fmt.Errorf("could not write to %s, error: %w", js.fileName, err); err != nil {
				return err
			}
		}
	} else {
		if err := fmt.Errorf("JSON error: %w", err); err != nil {
			return err
		}
	}

	return nil
}
