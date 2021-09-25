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
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; We need to ask the DB details
			askOnCLI()
		} else {
			log.Panicf("Failed to read the config file: %v", err)
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
