package config

import (
	"encoding/json"
	"os"
)

// Config holds the configuration settings for the KV store application.
type Config struct {
	UsersFile       string `json:"users_file"`
	DBFile          string `json:"db_file"`
	PublicDir       string `json:"public_dir"`
	IndexFile       string `json:"index_file"`
	TemplatesDir    string `json:"templates_dir"`
	HTTPAddr        string `json:"http_addr"` // eg: "0.0.0.0:8080"
	EnablePprof     bool   `json:"enable_pprof"`
	PprofAddr       string `json:"pprof_addr"`       // eg: "127.0.0.1:6060"
	TranslationsDir string `json:"translations_dir"` // Directory for translation files
	DefaultLanguage string `json:"default_language"` // Default language (e.g., "zh-Hans" for Simplified Chinese)
}

var Cfg Config

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &Cfg)
}
