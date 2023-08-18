package main

import (
	"github.com/pglomba/proxor/internal/config"
	"github.com/pglomba/proxor/internal/logger"
	"github.com/pglomba/proxor/internal/reverseproxy"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	var conf config.Config
	config.InitConfig()
	if err := conf.Load(); err != nil {
		log.Fatalf("Unable to load config file: %v", err)
	}

	logger, logFile := logger.New(conf.LogLevel, conf.LogPath)
	if logFile != nil {
		defer logFile.Close()
	}

	proxy := &reverseproxy.ProxyServer{
		Config: &conf,
		Logger: logger,
	}

	proxyErrCh := make(chan error)
	go func() {
		if err := proxy.Start(); err != nil {
			proxyErrCh <- err
		}
		close(proxyErrCh)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM)

	select {
	case err, ok := <-proxyErrCh:
		if ok {
			logger.Err(err).Msg("Proxor server error")
		}
	case sig := <-signalCh:
		logger.Info().Msgf("Signal %s received", sig)
		if err := proxy.Shutdown(5); err != nil {
			logger.Err(err).Msg("Failed to shutdown Proxor server")
		}
		logger.Info().Msg("Proxor server shutdown")
	}
}
