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
	"strconv"
	"sync"
	"time"
)

// const BIG_FILE_SIZE int64 = 1_000_000_000
// const LOG_PATH string = "./logs/" // testing
const BIG_FILE_SIZE int64 = (1024 * 1024 * 1024)
const DEBUG bool = false
const LOG_PATH string = "./tmp/"
const ERR_LOG_FILE string = LOG_PATH + "golarge_errors.log"

var unaccessibleDirectories int = 0
var unreadableFiles int = 0
var err_msg_pool []string

func main() {
	if DEBUG {
		// stat()
		PrintMemUsage()
	}

	tStart := time.Now()

	var fileOutput string
	var jsonOutput bool
	var help bool
	var verbose bool = false
	var concurrent = false

	flag.StringVar(&fileOutput, "o", "list.txt", "Write output to text file")
	flag.BoolVar(&jsonOutput, "j", false, "Ouput list to JSON file")
	flag.BoolVar(&verbose, "v", false, "Display warnings & errors")
	flag.BoolVar(&concurrent, "t", false, "Process file analysis concurrently")
	flag.BoolVar(&help, "h", false, "Display this help and exit")
	flag.Parse()
	args := flag.Args()

	if help {
		// flag.PrintDefaults()
		usage()
		os.Exit(0)
	}

	if flag.NArg() == 1 {
		dirPath := string(args[0])
		logging := false
		var logFile string

		if flag.NFlag() >= 1 && fileOutput != "" {
			logging = true
			logFile = fileOutput
		}
		fmt.Printf("\nLooking for files in %s\n\n", Colorize(dirPath, "blue", ""))
		fl := newFileList()
		if err := fl.listFiles(dirPath, logging, logFile, verbose, concurrent); err != nil {
			log.Fatal("%w", err)
		}
		if jsonOutput {
			jsonFile := newJsontFile(logFile, fl.data)
			jsonFile.saveToJson()
		}

		if !verbose {
			var d string
			for _, v := range err_msg_pool {
				d += v
			}
			txtFile := newTxtFile(ERR_LOG_FILE, d)
			if err := txtFile.saveToFile(1); err != nil {
				log.Fatal("%w", err)
			}
		}

		fmt.Println("")
		fmt.Println("Results:")
		fmt.Printf("  - found: %s %s\n", Colorize(strconv.Itoa(fl.count), "green", ""), Colorize("files with size bigger than 1GiB", "green", ""))
		fmt.Printf("  - unreadable file(s):       %s\n", Colorize(strconv.Itoa(unreadableFiles), "blue", ""))
		fmt.Printf("  - unaccessible directories: %s\n", Colorize(strconv.Itoa(unaccessibleDirectories), "blue", ""))
		fmt.Printf("  - processing time:          %.2fs\n\n", time.Now().Sub(tStart).Seconds())
		if logging {
			fmt.Printf("Logs:\n")
			fmt.Printf("  - List generated: %s\n", Colorize(logFile, "green", ""))
			if jsonOutput {
				fmt.Printf("  - List generated: %s\n", Colorize(ExtLess(logFile)+".json\n", "cyan", ""))
			}
			fmt.Println("")
		}
	} else {
		usage()
		os.Exit(1)
	}

	if DEBUG {
		PrintMemUsage()
	}
}

func (fl *fileList) processFile(filePath string, v fs.DirEntry, logging bool, logFile string, verbose bool, concurrent bool) {
	if info, err := v.Info(); err == nil {
		fileSize := info.Size()
		if fileSize >= BIG_FILE_SIZE {
			fileName := info.Name()
			data := fmt.Sprintf("%s => %dMiB", path.Join(filePath, fileName), fileSize/(1024*1024))
			fmt.Println(Colorize(data, "green", ""))
			nl := newList(fileName, filePath, path.Join(filePath, fileName), fileSize, filepath.Ext(fileName))
			if concurrent {
				fl.Lock()
			}
			fl.data = append(fl.data, nl)
			fl.count++
			if concurrent {
				fl.Unlock()
			}
			if logging {
				txtFile := newTxtFile(logFile, data)
				if err := txtFile.saveToFile(fl.count); err != nil {
					log.Fatal("%w", err)
				}
			}
		}
	} else {
		unreadableFiles++

		if !verbose {
			t := time.Now()
			data := fmt.Sprintf("[warning] %d/%d/%d %d:%d:%d - %s", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), err.Error())
			err_msg_pool = append(err_msg_pool, data)
		} else {
			fmt.Println("[warning]", Colorize(err.Error(), "yellow", ""))
		}
	}
}

func (fl *fileList) listFiles(filePath string, logging bool, logFile string, verbose bool, concurrent bool) error {
	if files, err := os.ReadDir(filePath); err == nil {
		if concurrent {
			wg := sync.WaitGroup{}
			for _, v := range files {
				if v.IsDir() {
					fl.listFiles(path.Join(filePath, v.Name()), logging, logFile, verbose, concurrent)
				} else {
					wg.Add(1)
					go func(v fs.DirEntry) {
						fl.processFile(filePath, v, logging, logFile, verbose, concurrent)
						wg.Done()
					}(v)
				}
			}
			wg.Wait()
		} else {
			for _, v := range files {
				if v.IsDir() {
					fl.listFiles(path.Join(filePath, v.Name()), logging, logFile, verbose, concurrent)
				} else {
					fl.processFile(filePath, v, logging, logFile, verbose, concurrent)
				}
			}
		}
	} else {
		unaccessibleDirectories++

		if !verbose {
			t := time.Now()
			data := fmt.Sprintf("[warning] %d/%d/%d %d:%d:%d - %s\n", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), err.Error())
			err_msg_pool = append(err_msg_pool, data)
		} else {

			fmt.Println("[warning]", Colorize(err.Error(), "yellow", ""))
		}
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
			if err := fmt.Errorf("[error] could not write to %s, error: %w", js.fileName, err); err != nil {
				return err
			}
		}
	} else {
		if err := fmt.Errorf("[error] JSON: %w", err); err != nil {
			return err
		}
	}

	return nil
}

func usage() {
	fmt.Println("")
	fmt.Println("Usage: golarge [OPTIONS] <PATH>")
	fmt.Println("")
	fmt.Println("Look for files bigger than 1GiB from given directory path")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("  -o <string>     Output path (default: list.txt)")
	fmt.Println("  -j      	  Enable export to JSON file")
	fmt.Println("  -v      	  Display warnings & error instead of logging to file")
	fmt.Println("  -t      	  Enable concurrent file processing")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  golarge /foo/bar")
	fmt.Println("  golarge -v /bar")
	fmt.Println("  golarge -o list.txt -j /foo/bar")
	fmt.Println("  golarge -t /foo")
	fmt.Println("")
}
