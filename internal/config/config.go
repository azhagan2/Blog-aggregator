package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// This is the config struct, where used to store the metada of the database

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

// Declaring a constant for storing the file name which is in root directory

const configFileName = ".gatorconfig.json"

// It basically set the username of the user, when this func is called from the main func and calls the write func

func (c *Config) SetUser(username string) error {
	c.CurrentUserName = username
	return write(*c)
}

// This is basically a func to get the full URL for the config JSON file

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, configFileName), nil
}

/* This write func is being called from the SetUser func, which is used to write in the config file. Basically,
this func gets the GO struct from the SetUser func, then converts it into again the JSON file using json.NewEncoder */

func write(cfg Config) error {
	ConfigFile, err := getConfigFilePath()
	if err != nil {
		return fmt.Errorf("error getting the filepath when writing: %w", err)
	}

	file, err := os.Create(ConfigFile)
	if err != nil {
		return fmt.Errorf("error creating/ opening a config file: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(&cfg); err != nil {
		return fmt.Errorf("error encoding to config to JSON: %w", err)
	}

	return nil
}

/* This Read func reads from the actual config file, by decoding the JSON to GO struct and return, basically this is called
from the main func, then it returns the GO struct, converted from the JSON file using json.NewDecoder(). */

func Read() (*Config, error) {

	configFile, err := getConfigFilePath()
	if err != nil {
		return &Config{}, fmt.Errorf("error getting the filepath: %w", err)
	}

	file, err := os.Open(configFile)
	if err != nil {
		return &Config{}, fmt.Errorf("error Opening the file: %w", err)
	}

	defer file.Close()

	var user Config

	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&user); err != nil {
		fmt.Println("Error decoding Json:", err)
		return &Config{}, fmt.Errorf("error decoding Json: %w", err)
	}
	return &user, nil
}
