MAIN_SERVICE_BINARY=main_service
STAFF_BINARY=staff_service
SURVEY_BINARY=survey_service

PROJECT_DIR := ${CURDIR}

DOCKER_DIR := ${CURDIR}/docker

## build: Build compiles project
build:
	go build -o ${MAIN_SERVICE_BINARY} cmd/main_service/start.go
	go build -o ${STAFF_BINARY} cmd/staff_service/start.go
	go build -o ${SURVEY_BINARY} cmd/survey_service/start.go

## build-docker: Builds all docker containers
build-docker:
	docker build -t dependencies -f ${DOCKER_DIR}/builder.Dockerfile .
	docker build -t main_service -f ${DOCKER_DIR}/main_service.Dockerfile .
	docker build -t staff_service -f ${DOCKER_DIR}/staff.Dockerfile .
	docker build -t survey_service -f ${DOCKER_DIR}/survey.Dockerfile .

## run-and-build: Build and run docker
build-and-run: build-docker
	docker-compose up

## run-background: Run process in background(available after build)
run-background:
	docker-compose up -d

## stop: Stop all containers on machine
stop:
	docker stop $(docker ps -a -q)

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command to run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
