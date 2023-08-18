package reverseproxy

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/pglomba/proxor/internal/certificates"
	"github.com/pglomba/proxor/internal/config"
	"github.com/pglomba/proxor/internal/mirror"
	"github.com/rs/zerolog"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

type ProxyServer struct {
	Server  *http.Server
	Handler http.Handler
	Config  *config.Config
	Logger  *zerolog.Logger
}

func (proxy *ProxyServer) Start() error {
	proxy.Handler = reverseProxyHandler(proxy.Config, proxy.Logger)
	port := fmt.Sprintf(":%d", proxy.Config.Proxy.Port)

	// Load certificates
	if proxy.Config.Proxy.TlsEnable {
		certs, err := certificates.LoadCertificates(proxy.Config.Certificates)
		if err != nil {
			return err
		}

		// Read TLS config
		var tlsMinVersion uint16
		if _, ok := config.MinVersion[proxy.Config.Proxy.TlsMinVersion]; ok {
			tlsMinVersion = config.MinVersion[proxy.Config.Proxy.TlsMinVersion]
		} else {
			err = errors.New(fmt.Sprintf("Wrong 'tlsMinVersion' value: %s", proxy.Config.Proxy.TlsMinVersion))
			return err
		}

		var tlsCipherSuites []uint16
		for _, tlsCipherSuite := range proxy.Config.Proxy.TlsCipherSuites {
			if _, ok := config.CipherSuites[tlsCipherSuite]; ok {
				tlsCipherSuites = append(tlsCipherSuites, config.CipherSuites[tlsCipherSuite])
			} else {
				err = errors.New(fmt.Sprintf("Wrong cipher suite value in 'tlsCipherSuites': %s", tlsCipherSuite))
				return err
			}
		}

		tlsConfig := &tls.Config{
			Certificates: certs,
			MinVersion:   tlsMinVersion,
			CipherSuites: tlsCipherSuites,
		}

		proxy.Server = &http.Server{
			Addr:         port,
			Handler:      proxy.Handler,
			TLSConfig:    tlsConfig,
			ReadTimeout:  time.Duration(proxy.Config.Proxy.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(proxy.Config.Proxy.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(proxy.Config.Proxy.IdleTimeout) * time.Second,
		}

		proxy.Logger.Info().Msgf("Starting Proxor server on %s", proxy.Server.Addr)
		if err := proxy.Server.ListenAndServeTLS("", ""); err != nil {
			return err
		}

	} else {

		proxy.Server = &http.Server{
			Addr:         port,
			Handler:      proxy.Handler,
			ReadTimeout:  time.Duration(proxy.Config.Proxy.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(proxy.Config.Proxy.WriteTimeout) * time.Second,
			IdleTimeout:  time.Duration(proxy.Config.Proxy.IdleTimeout) * time.Second,
		}

		proxy.Logger.Info().Msgf("Starting Proxor server on %s", proxy.Server.Addr)
		if err := proxy.Server.ListenAndServe(); err != nil {
			return err
		}
	}
	return nil
}

func (proxy *ProxyServer) Shutdown(timeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()
	if err := proxy.Server.Shutdown(ctx); err != nil {
		if err = proxy.Server.Close(); err != nil {
			return err
		}
	}

	return nil
}

func reverseProxyHandler(config *config.Config, logger *zerolog.Logger) http.Handler {
	logResponse := func(response *http.Response) (err error) {
		logger.Info().
			Str("proto", response.Request.Proto).
			Str("method", response.Request.Method).
			Int("response", response.StatusCode).
			Str("scheme", response.Request.URL.Scheme).
			Str("host", response.Request.Host).
			Str("uri", response.Request.RequestURI).
			Str("context", "proxy").
			Msg("")

		return nil
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error reading body: %s", err), 555)
			return
		}

		if config.Proxy.MirrorEnable {
			mirroredRequest := r.Clone(context.Background())

			r.Body = io.NopCloser(bytes.NewReader(body))

			mirroredRequest.Body = io.NopCloser(bytes.NewReader(body))

			go mirror.SendMirrorRequest(mirroredRequest, &config.Mirror, logger)
		}

		targetURL, err := url.Parse(config.Proxy.BackendUrl)
		if err != nil {
			logger.Error().Err(err).Msg("Failed to parse backend URL")
		}

		reverseProxy := httputil.NewSingleHostReverseProxy(targetURL)
		reverseProxy.ModifyResponse = logResponse
		reverseProxy.ServeHTTP(w, r)

	}
	return http.HandlerFunc(handler)
}
