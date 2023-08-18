package mirror

import (
	"crypto/tls"
	"github.com/pglomba/proxor/internal/config"
	"github.com/rs/zerolog"
	"net"
	"net/http"
	"net/url"
	"time"
)

func SendMirrorRequest(req *http.Request, config *config.Mirror, logger *zerolog.Logger) {
	targetURL, err := url.Parse(config.TargetUrl)
	if err != nil {
		logger.Error().Err(err).
			Str("context", "mirror").
			Msg("Failed to parse mirror URL")
	}

	requestURI := req.RequestURI

	req.Host = targetURL.Host
	req.URL.Host = targetURL.Host
	req.URL.Scheme = targetURL.Scheme
	req.RequestURI = ""

	client := http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: config.ClientTransportSkipVerify},
			DialContext: (&net.Dialer{
				Timeout:   time.Duration(config.ClientTransportTimeout) * time.Second,
				KeepAlive: time.Duration(config.ClientTransportKeepAlive) * time.Second,
			}).DialContext,
			TLSHandshakeTimeout:   time.Duration(config.ClientTransportTLSHandshakeTimeout) * time.Second,
			ResponseHeaderTimeout: time.Duration(config.ClientTransportResponseHeaderTimeout) * time.Second,
			ExpectContinueTimeout: time.Duration(config.ClientExpectContinueTimeout) * time.Second,
		},
		Timeout: time.Duration(config.ClientTimeout) * time.Second,
	}

	response, err := client.Do(req)
	if err != nil {
		logger.Error().Err(err).
			Str("context", "mirror").
			Msg("Failed to send mirrored request")
	} else {
		err = response.Body.Close()
		if err != nil {
			logger.Error().Err(err).Msg("Failed to close mirrored request body")
		} else {
			logger.Info().
				Str("proto", req.Proto).
				Str("method", req.Method).
				Int("response", response.StatusCode).
				Str("scheme", req.URL.Scheme).
				Str("host", req.Host).
				Str("uri", requestURI).
				Str("context", "mirror").
				Msg("")
		}
	}
}
