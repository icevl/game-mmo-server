package http

import (
	"errors"
	"os"
	"server/config"
)

func getWorldFile() (*os.File, error) {
	file, err := os.Open(config.WorldFilePath)
	if err != nil {
		return nil, errors.New("error opening file")
	}

	return file, nil
}
