package types

type List struct {
	Name     string `json:"name"`
	BasePath string `json:"basePath"`
	FullPath string `json:"fullPath"`
	Size     int64  `json:"size"`
}
