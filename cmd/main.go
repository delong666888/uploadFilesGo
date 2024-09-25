package main

import (
	"fmt"

	uploadFilesGo "github.com/delong666888/uploadFilesGo"
	"github.com/delong666888/uploadFilesGo/config"
)

func main() {
	uploadConfig := config.NewUploadConfig("config.yaml")
	uploadFilesGo.CheckExistingFiles(uploadConfig)
	fmt.Println("Watching directory:", uploadConfig.SourceDir)
	uploadFilesGo.WatchDirectory(uploadConfig)
}
