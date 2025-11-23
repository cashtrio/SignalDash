package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type ViperOptions struct {
	Filename  string
	ConfigDir string
}

func (vo *ViperOptions) InitConfig() *viper.Viper {
	wd, gErr := os.Getwd()
	if gErr != nil {
		log.Fatalf("InitConfig.go: failed to read working directory: %s", gErr.Error())
	}

	v := viper.New()
	v.AddConfigPath(wd + "/" + vo.ConfigDir)
	v.SetConfigName(vo.Filename)

	log.Println("InitConfig.go: ", wd+"/"+vo.ConfigDir)

	if rErr := v.ReadInConfig(); rErr != nil {
		log.Fatalf("InitConfig.go: failed to init viper: %s", rErr.Error())
	}

	return v
}
