package config

import (
	"github.com/spf13/viper"
	"testing"
)

func TestLoad(t *testing.T) {
	viper.AddConfigPath("../../testdata")
	var conf Config
	InitConfig()

	t.Run("Require proxy.port parameter", func(t *testing.T) {
		err := conf.Load()

		if err == nil {
			t.Fatal("expected an error")
		}
	})

	t.Run("Require proxy.backendUrl parameter", func(t *testing.T) {
		viper.SetDefault("proxy.port", 443)
		err := conf.Load()

		if err == nil {
			t.Fatal("expected an error")
		}
	})

	t.Run("Require mirror.targetUrl parameter", func(t *testing.T) {
		viper.SetDefault("proxy.port", 443)
		viper.SetDefault("proxy.backendUrl", "http://localhost")
		viper.SetDefault("proxy.mirrorEnable", true)
		err := conf.Load()

		if err == nil {
			t.Fatal("expected an error")
		}
	})
}
