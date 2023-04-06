.PHONY: build run test deploy status remove
build:
	go build -v ./cmd/logserver

run: 
	go build -v ./cmd/logserver
	./logserver

test:
	go test -v ./...

deploy:
	make build
	- sudo systemctl stop logserver.service || true
	sudo cp logserver /usr/bin/
	sudo cp logserver.service /etc/systemd/system/
	sudo mkdir -p /etc/logserver
	sudo chown pi:pi /etc/logserver
	sudo cp config/logserver.toml /etc/logserver/
	sudo systemctl daemon-reload
	sudo systemctl enable logserver.service
	sudo systemctl start logserver.service

status:
	sudo systemctl status logserver.service

remove:
	sudo systemctl stop logserver.service
	sudo rm /usr/bin/logserver
	sudo rm /etc/logserver/logserver.toml
	sudo rm /etc/logserver -rf
	sudo rm /etc/systemd/system/logserver.service


.DEFAULT_GOAL: build