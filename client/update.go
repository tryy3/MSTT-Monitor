package client

import (
	"strings"
)

type UpdateResponse struct {
	Error   string
	Message string
}

func UpdateCheck(cmd Command) UpdateResponse {
	url := ""

	for _, args := range cmd.Params {
		if strings.ToLower(args.Name) == "-url" {
			url = args.Value
		}
	}

	if url == "" {
		return UpdateResponse{Error: "Please supply the -url flag"}
	}

	err := updater.Update(url)
	if err != nil {
		return UpdateResponse{Error: err.Error()}
	}

	return UpdateResponse{Message: "Updated version: " + version}
}
