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

var total int = 0
var logging bool = false

func main() {
	var path string
	var logFile string
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
		path = string(args[0])

		if flag.NFlag() >= 1 && fileOutput != "" {
			logging = true
			logFile = fileOutput
		}
		log.Printf("-- Searching large files in %s --", Colorize(path, "cyan"))
		fl := newFileList()
		if err := fl.listFiles(path, logFile); err != nil {
			log.Fatal("%w", err)
		}
		if jsonOutput {
			jsonFile := newJsontFile(logFile, fl.Data)
			jsonFile.saveToJson()
		}
		log.Printf("-- Found %d files with size bigger than 1GB --", total)
		if logging {
			log.Printf("-- List generated: %s --", Colorize(logFile, "cyan"))
		}
	} else {
		fmt.Println("\nUsage: golarge [OPTIONS] PATH")
		fmt.Println("\nUtil to find files around 1GB of size from given directory path")
		fmt.Println("\nOptions:")
		fmt.Println("  -o, --output string   Full path of file to save large file list to (default: list.txt)")
		fmt.Println("  -j, --json      			 Enable export to JSON file")
		fmt.Println("\nExamples:")
		fmt.Println("  golarge /foo/bar")
		fmt.Println("  golarge -o list.txt /foo/bar")
		fmt.Println("  golarge --output list.txt --json /foo/bar")
		fmt.Println("")
		os.Exit(0)
	}
}

func (fl *FileList) listFiles(filePath string, logFile string) error {
	if files, err := os.ReadDir(filePath); err == nil {
		wg := &sync.WaitGroup{}
		for _, v := range files {
			if v.IsDir() {
				fl.listFiles(path.Join(filePath, v.Name()), logFile)
			} else {
				wg.Add(1)
				go func(v fs.DirEntry) {
					if info, err := v.Info(); err == nil {
						name := info.Name()
						size := info.Size()
						if size >= BIG_FILE_SIZE {
							data := fmt.Sprintf("%s => %dMiB", path.Join(filePath, name), size/(1024*1024))
							log.Println(Colorize(data, "green"))
							total++
							fl.MU.Lock()
							fl.Data = append(fl.Data, List{Name: name, BasePath: filePath, FullPath: path.Join(filePath, name), Size: size, Type: filepath.Ext(name)})
							fl.MU.Unlock()
							if logging {
								txtFile := newTxtFile(logFile, data)
								if err := txtFile.saveToFile(); err != nil {
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

func (tf *TxtFile) saveToFile() error {
	tf.MU.Lock()
	defer tf.MU.Unlock()

	f, err := os.OpenFile(tf.FileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		if err := fmt.Errorf("%w", err); err != nil {
			return err
		}
	}
	defer f.Close()

	if _, err := os.Stat(tf.FileName); err == nil && total == 1 {
		f.Truncate(0)
	} else if err != nil {
		if err := fmt.Errorf("%w", err); err != nil {
			return err
		}
	}

	if _, err := f.WriteString(fmt.Sprintf("%s\n", tf.Data)); err != nil {
		if err := fmt.Errorf("%w", err); err != nil {
			return err
		}
	}

	return nil
}

func (js *JsonFile) saveToJson() error {
	if content, err := json.Marshal(&js.Data); err == nil {
		fileName := ExtLess(js.FileName)
		err := os.WriteFile(fmt.Sprintf("%s.json", fileName), content, 0644)
		if err != nil {
			if err := fmt.Errorf("could not write to %s, error: %w", js.FileName, err); err != nil {
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

func newFileList() *FileList {
	return &FileList{
		Data: nil,
		MU:   sync.Mutex{},
	}
}

func newTxtFile(fileName string, data string) *TxtFile {
	return &TxtFile{
		FileName: fileName,
		MU:       sync.Mutex{},
		Data:     data,
	}
}

func newJsontFile(fileName string, data []List) *JsonFile {
	return &JsonFile{
		FileName: fileName,
		Data:     data,
	}
}
