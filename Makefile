B=$(shell git rev-parse --abbrev-ref HEAD)
BRANCH=$(subst /,-,$(B))
GITREV=$(shell git describe --abbrev=7 --always --tags)
REV=$(GITREV)-$(BRANCH)-$(shell date +%Y%m%d)

# get current user name
USER=$(shell whoami)
# get current user group
GROUP=$(shell id -gn)

.PHONY: build run test deploy status remove
build:
	- cd app; go build -v -o ../logserver; cd ..
	- make info

run: build
	- ./logserver

test:
	go test -v ./...

info:
	- @echo "revision $(REV)"
	- @echo "branch $(BRANCH)"

deploy: build
	- sudo systemctl stop logserver.service || true
	sudo cp logserver /usr/bin/
	sed -i "s/%USER%/$(USER)/g" logserver.service
	sudo cp logserver.service /etc/systemd/system/
	sudo mkdir -p /etc/logserver
	sudo chown ${USER}:${GROUP} /etc/logserver
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