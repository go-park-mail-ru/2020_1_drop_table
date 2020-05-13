MAIN_SERVICE_BINARY=main_service
STAFF_BINARY=staff_service
SURVEY_BINARY=survey_service

PROJECT_DIR := ${CURDIR}

DOCKER_DIR := ${CURDIR}/docker

## build: Build compiles project
build:
	go build -o ${MAIN_SERVICE_BINARY} internal/app/main/start.go
	go build -o ${STAFF_BINARY} internal/microservices/staff/main/start.go
	go build -o ${SURVEY_BINARY} internal/microservices/survey/main/start.go

## build-docker: Builds all docker containers
build-docker:
	docker build -t dependencies -f ${DOCKER_DIR}/builder.Dockerfile .
	docker build -t main_service -f ${DOCKER_DIR}/main_service.Dockerfile .
	docker build -t staff_service -f ${DOCKER_DIR}/staff.Dockerfile .
	docker build -t survey_service -f ${DOCKER_DIR}/survey.Dockerfile .

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command to run:"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
