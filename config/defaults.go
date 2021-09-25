package config

const (
	defaultPort   = 9876
	defaultEnv    = "prod"
	defaultDBHost = "localhost"
	defaultDBPort = 5432
	defaultDBUser = "postgres"
	defaultDBName = "blog"
)

var defaults = map[string]interface{}{
	keyPort:   defaultPort,
	keyEnv:    defaultEnv,
	keyDBHost: defaultDBHost,
	keyDBPort: defaultDBPort,
	keyDBUser: defaultDBUser,
	keyDBName: defaultDBName,
}

func isProduction(env string) bool {
	return env == "prod"
}
