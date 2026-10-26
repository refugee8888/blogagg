package config

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	DBurl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	filePath := filepath.Join(home, configFileName)
	return filePath, nil
}

func write(cfg Config) error {
	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	b, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath, b, os.ModePerm)
}

func Read() (Config, error) {
	filepath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	f, err := os.Open(filepath)
	if err != nil {
		return Config{}, err
	}
	defer f.Close()

	r, err := io.ReadAll(f)
	if err != nil {
		return Config{}, err
	}

	var dump Config
	if err := json.Unmarshal(r, &dump); err != nil {
		return Config{}, err
	}
	//fmt.Println(dump)

	return dump, nil
}

func (cfg *Config) SetUser(user string) error {
	cfg.CurrentUserName = user
	err := write(*cfg)
	if err != nil {
		return err
	}
	return nil
}
