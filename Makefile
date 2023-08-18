build :
	cd cmd/ && CGO_ENABLED=0 go build -o ../bin/proxor && cd ../bin && sudo setcap CAP_NET_BIND_SERVICE=+ep proxor
run : build
	cd bin/ && ./proxor