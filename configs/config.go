package configs

import "github.com/spf13/viper"

type DataBase struct {
	Driver   string `mapstructure:"driver"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"db_name"`
}

type Config struct {
	Port string `mapstructure:"port"`
	DB   DataBase
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")
	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	var cfg Config
	cfg.Port = viper.GetString("port")
	err = viper.UnmarshalKey("db", &cfg.DB)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
