# Include variable from .envrc file
include .env

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

.PHONY: confirm
confirm:
	@echo -n 'Are you sure? [y/N] ' && read ans && [ $${ans:-N} = y ]

## run: run the web application
.PHONY: run
run:
	go run ./cmd/web

## air: autorun using AIR
.PHONY: air
air:
	air -c .air.toml

## build: build minimzed app and save to .tmp directory
.PHONY: build
build:
	go build -ldflags='-s' -o=./tmp/${SVC} .
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./tmp/linux_amd64/${SVC} .

## rebuild: rebuild, copy, and restart service
.PHONY: rebuild
rebuild: build svc/stop copy svc/start

## copy: copy new build into place
.PHONY: copy
copy:
	cp tmp/linux_amd64/${SVC} /srv/${SVC}/

## svc: create service
.PHONY: svc
svc:
	cp ${SVC}.service /etc/systemd/system/${SVC}.service

## svc/start: start service 
.PHONY: svc/start
svc/start:
	systemctl start ${SVC}.service

## svc/stop: stop service 
.PHONY: svc/stop
svc/stop:
	systemctl stop ${SVC}.service
	
## host/ssh: ssh to host server
.PHONY: host/ssh
host/ssh:
	ssh -i ${USERKEY} root@${HOST}

## host/push: copy linux build to host server
.PHONY: host/push
host/push:
	scp -i ${USERKEY} ./tmp/linux_amd64/${SVC} root@${HOST}:/srv/${SVC}/${SVC}

## host/svc: copy service file to host server
.PHONY: host/svc
host/svc:
	scp -i ${USERKEY} ./${SVC}.service root@${HOST}:/etc/systemd/system/${SVC}.service

## host/start: start service
.PHONY: host/start
host/start:
	ssh -i ${USERKEY} root@${HOST} 'systemctl start ${SVC}.service'

## host/stop: stop service
.PHONY: host/stop
host/stop:
	ssh -i ${USERKEY} root@${HOST} 'systemctl stop ${SVC}.service'
