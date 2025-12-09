package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type environment struct {
	Instance *viper.Viper
}

func NewEnvironment(configDir, filename string) *environment {
	wd, gErr := os.Getwd()
	if gErr != nil {
		log.Fatalf("NewEnvironment.go: failed to read working directory: %s", gErr.Error())
	}

	v := viper.New()
	v.AddConfigPath(wd + "/" + configDir)
	v.SetConfigName(filename)

	log.Println("NewEnvironment.go: ", wd+"/"+configDir)

	if rErr := v.ReadInConfig(); rErr != nil {
		log.Fatalf("NewEnvironment.go: failed to init viper: %s", rErr.Error())
	}

	return &environment{
		Instance: v,
	}
}
