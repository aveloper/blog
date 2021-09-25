package config

type (
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
	}

	Logger struct {
	}

	App struct {
		Port        int
		Environment string
		Production  bool
		DB          *DB
		Logger      *Logger
	}
)
