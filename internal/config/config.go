package config

import (
	"fmt"
	"os"
	"rest-geoip/internal/random"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type config struct {
	Program struct {
		APIKey        string `mapstructure:"API_KEY"`
		EnableLogging bool   `mapstructure:"ENABLE_LOGGING"`
		EnableWeb     bool   `mapstructure:"ENABLE_WEB"`
		ListenAddress string `mapstructure:"LISTEN_ADDRESS"`
		ListenPort    string `mapstructure:"LISTEN_PORT"`
		ReleaseMode   bool   `mapstructure:"RELEASE_MODE"`
	} `mapstructure:"PROGRAM"`
	Maptiler struct {
		Token string `mapstructure:"TOKEN"`
	} `mapstructure:"MAPTILER"`
	Maxmind struct {
		LicenseKey string `mapstructure:"LICENSE_KEY"`
		DBLocation string `mapstructure:"DB_LOCATION"`
		DBFileName string `mapstructure:"DB_FILE_NAME"`
	} `mapstructure:"MAXMIND"`
}

var programConfig config
var once = sync.Once{}

func Init() {
	once.Do(func() {
		replacer := strings.NewReplacer(".", "__")
		viper.SetEnvKeyReplacer(replacer)
		viper.SetEnvPrefix("goip")

		// When you explicitly provide the ENV variable name (the second parameter), it does not automatically add the prefix
		viper.BindEnv("program.api_key", "GOIP_PROGRAM__API_KEY")
		viper.BindEnv("program.enable_logging", "GOIP_PROGRAM__ENABLE_LOGGING")
		viper.BindEnv("program.enable_web", "GOIP_PROGRAM__ENABLE_WEB")
		viper.BindEnv("program.listen_address", "GOIP_PROGRAM__LISTEN_ADDRESS")
		viper.BindEnv("program.listen_port", "GOIP_PROGRAM__LISTEN_PORT")
		viper.BindEnv("program.release_mode", "GOIP_PROGRAM__RELEASE_MODE")
		viper.BindEnv("maptiler.token", "GOIP_MAPTILER__TOKEN")
		viper.BindEnv("maxmind.license_key", "GOIP_MAXMIND__LICENSE_KEY")
		viper.BindEnv("maxmind.db_location", "GOIP_MAXMIND__DB_LOCATION")
		viper.BindEnv("maxmind.db_file_name", "GOIP_MAXMIND__DB_FILE_NAME")

		viper.ReadInConfig()
		viper.Unmarshal(&programConfig)

		/*
			https://github.com/spf13/viper/issues/1502
			BindEnv doesn't work well with SetDefault

			if viper.GetString("PROGRAM__API_KEY") == "" {
				generatedKey, err := random.GenerateKey(512)
				if err != nil {
					fmt.Println("Error: no api key set and we weren't able to generate a key. Exiting")
					os.Exit(1)
				}

				viper.SetDefault("PROGRAM__API_KEY", generatedKey)
				fmt.Printf("Generated API key: %s\n", viper.GetString("PROGRAM__API_KEY"))
			}
			viper.SetDefault("PROGRAM__ENABLE_LOGGING", false)
			viper.SetDefault("PROGRAM__ENABLE_WEB", true)
			viper.SetDefault("PROGRAM__LISTEN_ADDRESS", "0.0.0.0")
			viper.SetDefault("PROGRAM__LISTEN_PORT", "1323")
			viper.SetDefault("PROGRAM__RELEASE_MODE", "true")
			viper.SetDefault("MAPTILER__TOKEN", "if_you're_seeing_this_replace_me_in_env_var")
			viper.SetDefault("MAXMIND__DB_LOCATION", "/opt/")
			viper.SetDefault("MAXMIND__DB_NAME", "GeoLite2-City.mmdb")
		*/

		if programConfig.Program.ListenAddress == "" {
			programConfig.Program.ListenAddress = "0.0.0.0"
		}
		if programConfig.Program.ListenPort == "" {
			programConfig.Program.ListenPort = "1323"
		}
		if programConfig.Maptiler.Token == "" {
			programConfig.Maptiler.Token = "if_you're_seeing_this_replace_me_in_env_var"
		}
		if programConfig.Maxmind.DBLocation == "" {
			programConfig.Maxmind.DBLocation = "/opt/"
		}
		if programConfig.Maxmind.DBFileName == "" {
			programConfig.Maxmind.DBFileName = "GeoLite2-City.mmdb"
		}
		if programConfig.Program.APIKey == "" {
			generatedKey, err := random.GenerateKey(512)
			if err != nil {
				fmt.Println("Error: no api key set and we weren't able to generate a key. Exiting")
				os.Exit(1)
			}

			programConfig.Program.APIKey = generatedKey
			fmt.Printf("Generated API key: %s\n", programConfig.Program.APIKey)
		}
	})
}

func Details() config {
	return programConfig
}
