package config

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/spf13/viper"
)

var (
	cfg  *App
	once sync.Once
)

func Get() *App {
	once.Do(func() {
		setDefaults()
		readConfig()
		saveConfig()

		cfg = &App{
			Port:        viper.GetInt(keyPort),
			Environment: viper.GetString(keyEnv),
			Production:  isProduction(viper.GetString(keyEnv)),
			DB: &DB{
				Host:     viper.GetString(keyDBHost),
				Port:     viper.GetInt(keyDBPort),
				User:     viper.GetString(keyDBUser),
				Password: viper.GetString(keyDBPass),
				Name:     viper.GetString(keyDBName),
			},
			Logger: &Logger{},
		}
	})

	return cfg
}

// Reset for deleting the .blog.yaml file from home directory
func Reset() {
	//1. get the .blog.yaml file
	//2. delete the file

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to find home directory")
	}

	err = os.Remove(homeDir+"/.blog.yaml")
	if err != nil {
		log.Fatalf("Failed to delete the file")
	}
}

func saveConfig() {
	if err := viper.SafeWriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
			return
		}

		log.Panicf("Failed to write config: %v", err)
	}
}

func setDefaults() {
	for key, value := range defaults {
		viper.SetDefault(key, value)
	}
}

func readConfig() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to find home directory")
	}

	viper.SetConfigName(".blog")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(homeDir)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; We need to ask the DB details
			askOnCLI()
		} else {
			log.Fatalf("Failed to read the config file: %v", err)
		}
	}
}

func askOnCLI() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Enter following details...")

	dbHost := readString(reader, "Database hostname", defaultDBHost, false)
	dbPort := readInt(reader, "Database port", defaultDBPort, false)
	dbUser := readString(reader, "Database user", defaultDBUser, false)
	dbPass := readPassword("Database password")
	dbName := readString(reader, "Database name", defaultDBName, false)

	viper.Set(keyDBHost, dbHost)
	viper.Set(keyDBPort, dbPort)
	viper.Set(keyDBUser, dbUser)
	viper.Set(keyDBPass, dbPass)
	viper.Set(keyDBName, dbName)
}
