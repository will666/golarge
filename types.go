package main

import "sync"

type list struct {
	FileName string `json:"fileName"`
	BasePath string `json:"basePath"`
	FullPath string `json:"fullPath"`
	FileSize int64  `json:"fileSize"`
	FileType string `json:"fileType"`
}

type fileList struct {
	sync.Mutex
	data  []list
	count int
}

type txtFile struct {
	sync.Mutex
	fileName string
	data     string
}

type jsonFile struct {
	fileName string
	data     []list
}

func newList(fileName string, basePath string, fullPath string, fileSize int64, fileType string) list {
	return list{
		FileName: fileName,
		BasePath: basePath,
		FullPath: fullPath,
		FileSize: fileSize,
		FileType: fileType,
	}
}

func newFileList() *fileList {
	return &fileList{
		data:  nil,
		count: 0,
	}
}

func newTxtFile(fileName string, data string) *txtFile {
	return &txtFile{
		fileName: fileName,
		data:     data,
	}
}

func newJsontFile(fileName string, data []list) jsonFile {
	return jsonFile{
		fileName: fileName,
		data:     data,
	}
}
