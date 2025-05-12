package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	UsersFile    string `json:"users_file"`
	DBFile       string `json:"db_file"`
	PublicDir    string `json:"public_dir"`
	IndexFile    string `json:"index_file"`
	TemplatesDir string `json:"templates_dir"`
	HTTPAddr     string `json:"http_addr"` // eg: "0.0.0.0:8080"
	EnablePprof  bool   `json:"enable_pprof"`
	PprofAddr    string `json:"pprof_addr"` // eg: "127.0.0.1:6060"
}

var Cfg Config

func LoadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &Cfg)
}
