<h1 align="center">
  <br>
  Proxor
  <br>
</h1>

<h4 align="center">A simple reverse proxy with HTTP request mirroring capability</h4>

<p align="center">
	<a href="https://github.com/pglomba/proxor/actions/workflows/ci.yaml"><img src="https://github.com/pglomba/proxor/actions/workflows/ci.yaml/badge.svg"></a>
</p>

### Features
* Reverse proxy requests to the backend URL
* Mirror requests to the mirror URL
* TLS support
* Config file based configuration

### Config
Proxor service looks for a configuration file `config.yaml` in a directory the tool's binary file is. 

Example config.yaml
```yaml
proxy:
  port: 443
  readTimeout: 5
  writeTimeout: 10
  idleTimeout: 15
  backendUrl: "https://foo.bar.net"
  tlsEnable: True
  mirrorEnable: True
  tlsMinVersion: "VersionTLS12"
  tlsCipherSuites:
    - "TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384"
    - "TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256"
    - "TLS_RSA_WITH_AES_256_GCM_SHA384"
    - "TLS_RSA_WITH_AES_128_GCM_SHA256"
certificates:
  - name: "bar.net"
    certificateFile: "./bar.net.crt"
    certificateChainFile: "./bar.net-intermediate.crt"
    certificateKey: "./bar.net.key"
logLevel: "info"
mirror:
  targetUrl: "https://mirror"
  clientTimeout: 5
  clientTransportTimeout: 5
  clientTransportKeepAlive: 30
  clientTransportTLSHandshakeTimeout: 10
  clientTransportResponseHeaderTimeout: 10
  clientTransportSkipVerify: False
  clientExpectContinueTimeout: 5
```
### Run
Proxor service requires `CAP_NET_BIND_SERVICE` capability to bind a socket to privileged ports.
```bash
$ git clone "https://github.com/pglomba/proxor.git"
$ make run
```


