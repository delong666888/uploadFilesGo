package config

import (
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Payload struct {
	Token   string `yaml:"token"`
	FileTxt string `yaml:"file_txt"`
}

type UploadConfig struct {
	ApiUrl    string `yaml:"api_url"`
	SourceDir string `yaml:"source_dir"`
	TargetDir string `yaml:"target_dir"`
	*Payload
}

func NewUploadConfig(filePath string) *UploadConfig {
	configFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal("open config file error:", err)
	}
	defer configFile.Close()
	uploadConfig := new(UploadConfig)
	decoder := yaml.NewDecoder(configFile)
	err = decoder.Decode(&uploadConfig)
	if err != nil {
		log.Fatal("decode uploadConfig error:", err)
	}
	return uploadConfig
}
