package client

import (
	"os"
	"strings"
	"time"
)

type File struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	IsDir   bool      `json:"isdir"`
	Perms   string    `json:"perms"`
	LastMod time.Time `json:"lastmod"`
}

type FileResponse struct {
	Error string `json:"error"`
	MFile File   `json:"mfile"`
}

func FileCheck(cmd Command) FileResponse {
	file := ""
	for _, args := range cmd.Params {
		if strings.ToLower(args.Name) == "-file" {
			file = args.Value
		}
	}

	if file == "" {
		return FileResponse{Error: "Invalid value on the flag -file"}
	}

	f, err := os.Open(file)
	if err != nil {
		return FileResponse{Error: err.Error()}
	}

	fi, err := f.Stat()
	if err != nil {
		return FileResponse{Error: err.Error()}
	}

	resp := FileResponse{
		MFile: File{
			Name:    fi.Name(),
			Size:    fi.Size(),
			IsDir:   fi.IsDir(),
			Perms:   fi.Mode().String(),
			LastMod: fi.ModTime(),
		},
	}
	return resp
}
