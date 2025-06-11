package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

const CONFIG_FILE_NAME = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func Read() (Config, error) {
	config_file_path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	config_file, err := os.Open(config_file_path)
	if err != nil {
		log.Printf("Failed to open config file at '%v' due to error: %v", config_file_path, err)
		return Config{}, err
	}
	defer config_file.Close()

	jd := json.NewDecoder(config_file)
	config := Config{}
	err = jd.Decode(&config)
	if err != nil {
		log.Println("ERROR: Failed to decode json: ", err)
		return Config{}, err
	}
	return config, nil
}

func getConfigFilePath() (string, error) {
	home_dir, err := os.UserHomeDir()
	if err != nil {
		log.Println("ERROR: Could not find home directory")
		return "", err
	}
	return filepath.Join(home_dir, CONFIG_FILE_NAME), nil
}

func (cfg *Config) SetUser(user_name string) error {
	cfg.CurrentUserName = user_name

	err := write(*cfg)
	if err != nil {
		log.Println("ERROR: Failed to write updated config to file: ", err)
		return err
	}

	return nil
}

func write(cfg Config) error {
	config_file_path, err := getConfigFilePath()
	if err != nil {
		return err
	}

	config_file, err := os.OpenFile(config_file_path, os.O_WRONLY, 644)
	if err != nil {
		log.Printf("ERROR: Failed to open config file at '%v' due to error: %v", config_file_path, err)
		return err
	}
	defer config_file.Close()

	je := json.NewEncoder(config_file)
	err = je.Encode(cfg)
	if err != nil {
		log.Println("ERROR: Failed to encode json: ", err)
		return err
	}
	return nil
}

func (cfg Config) String() string {
	return fmt.Sprintf("Config { DbUrl: '%v', CurrentUserName: '%v' }", cfg.DbUrl, cfg.CurrentUserName)
}
