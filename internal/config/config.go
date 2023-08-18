package config

import (
	"errors"
	"github.com/pglomba/proxor/internal/certificates"
	"github.com/spf13/viper"
	"log"
)

type Proxy struct {
	Port            uint16
	ReadTimeout     uint16
	WriteTimeout    uint16
	IdleTimeout     uint16
	BackendUrl      string
	MirrorEnable    bool
	MirrorUrl       string
	TlsEnable       bool
	TlsMinVersion   string
	TlsCipherSuites []string
}

type Mirror struct {
	TargetUrl                            string
	ClientTimeout                        uint16
	ClientTransportTimeout               uint16
	ClientTransportKeepAlive             uint16
	ClientTransportTLSHandshakeTimeout   uint16
	ClientTransportResponseHeaderTimeout uint16
	ClientTransportSkipVerify            bool
	ClientExpectContinueTimeout          uint16
}

type Config struct {
	Proxy        Proxy
	Mirror       Mirror
	Certificates []certificates.Certificate
	LogLevel     string
	LogPath      string
}

func InitConfig() {
	// Default values
	viper.SetDefault("logLevel", "info")

	viper.SetDefault("proxy.readTimeout", 5)
	viper.SetDefault("proxy.writeTimeout", 10)
	viper.SetDefault("proxy.idleTimeout", 15)
	viper.SetDefault("proxy.tlsEnable", false)
	viper.SetDefault("proxy.tlsMinVersion", "VersionTLS12")
	viper.SetDefault("proxy.tlsCipherSuites", []string{"TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384", "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256", "TLS_RSA_WITH_AES_256_GCM_SHA384", "TLS_RSA_WITH_AES_128_GCM_SHA256"})

	viper.SetDefault("mirror.clientTimeout", 5)
	viper.SetDefault("mirror.clientTransportTimeout", 5)
	viper.SetDefault("mirror.clientTransportKeepAlive", 5)
	viper.SetDefault("mirror.clientTransportTLSHandshakeTimeout", 5)
	viper.SetDefault("mirror.clientTransportResponseHeaderTimeout", 5)
	viper.SetDefault("mirror.clientTransportSkipVerify", false)
	viper.SetDefault("mirror.clientExpectContinueTimeout", 5)

	viper.SetConfigName("config.yaml")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %s", err.Error())
	} else {
		log.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func (conf *Config) Load() error {
	err := viper.Unmarshal(&conf)
	if err != nil {
		return err
	}

	// Validate required config
	if !viper.IsSet("proxy.port") {
		return errors.New("proxy.port parameter is missing")
	}
	if !viper.IsSet("proxy.backendUrl") {
		return errors.New("proxy.backendUrl parameter is missing")
	}

	if viper.GetBool("proxy.mirrorEnable") {
		if !viper.IsSet("mirror.targetUrl") {
			return errors.New("mirror.targetUrl parameter is missing")
		}
	}

	return nil
}
