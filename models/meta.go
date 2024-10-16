package models

// Meta represents metadata of a file
type Meta struct {
	Filename string `json:"filename"`
	Absolute string `json:"absolute"`
	IsDir    bool   `json:"is_dir"`
	Size     int64  `json:"size"`
}
