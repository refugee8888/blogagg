package config

type Config struct {
	DBurl           string `json:"connection_string_goes_here"`
	CurrentUserName string `json:"username_goes_here"`
}

const configFileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	return "", nil
}

func Read(filepath string) (error, Config) {
	return nil, Config{}
}

func SetUser(cfg *Config) error {
	return nil
}
