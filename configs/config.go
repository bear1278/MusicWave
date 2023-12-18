package configs

import "github.com/spf13/viper"

type DataBase struct {
	Driver   string `mapstructure:"driver"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	DbName   string `mapstructure:"db_name"`
}

type Spotify struct {
	RedirectURI  string `mapstructure:"redirectURI"`
	ClientID     string `mapstructure:"clientID"`
	ClientSecret string `mapstructure:"clientSecret"`
	State        string `mapstructure:"state"`
}

type Email struct {
	From     string `mapstructure:"from"`
	Password string `mapstructure:"password"`
}

type Config struct {
	Port    string `mapstructure:"port"`
	DB      DataBase
	Spotify Spotify
	Email   Email
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
	err = viper.UnmarshalKey("spotify", &cfg.Spotify)
	if err != nil {
		return nil, err
	}
	err = viper.UnmarshalKey("email", &cfg.Email)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
