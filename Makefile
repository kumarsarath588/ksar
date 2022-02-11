setup_deps:
	go mod tidy
	go mod vendor

build:
	go build -o ./tabsquare main.go

run:
	go run main.go

setup_docker_centos:
	sudo yum install -y --quiet yum-utils
	sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo
	sudo yum install -y --quiet docker-ce docker-compose

docker_image_build: setup_docker_centos
	docker build --no-cache -t tabsquare:1.0 .

docker_compose_up: setup_docker_centos
	docker-compose up -d

docker_compose_down: setup_docker_centos
	docker-compose down --volumes

all: setup_deps build