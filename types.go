package main

import "sync"

type List struct {
	Name     string `json:"name"`
	BasePath string `json:"basePath"`
	FullPath string `json:"fullPath"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
}

type FileList struct {
	Data []List
	MU   sync.Mutex
}

type TxtFile struct {
	FileName string
	MU       sync.Mutex
	Data     string
}

type JsonFile struct {
	FileName string
	Data     []List
}
