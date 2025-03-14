ifneq (,$(wildcard .env))
    include .env
endif

CONNECTION = $(REMOTE_USERNAME)@$(REMOTE_HOSTNAME)
MODE ?= release

DOCKERFILE := $(if $(filter $(MODE),debug),Dockerfile-debug,Dockerfile)

#.PHONY: run-docker
#run-docker:
#	@docker rm -f vdovin-auth-container
#	@docker run --name vdovin-auth-container -p 2345:2345 -p 8080:8080 vdovin-auth:latest

.PHONY: build-docker
build-docker:
	@docker buildx build --platform linux/amd64 -t vdovin-auth:latest ../. -f $(DOCKERFILE)

.PHONY: docker-image-tar
docker-image-tar: build-docker
	@docker save vdovin-auth:latest > vdovin-auth-image.tar

.PHONY: upload-image-tar
upload-image-tar:
	@echo ">> Uploading image tar"
	@ssh $(CONNECTION) "mkdir -p ~/images"
	@rsync -a --progress vdovin-auth-image.tar $(CONNECTION):~/images
	@ssh $(CONNECTION) "docker rmi -f vdovin-auth"
	@ssh $(CONNECTION) "docker load < ~/images/vdovin-auth-image.tar"

.PHONY: run-remote-docker
run-remote-docker:
	@echo ">> Starting remote docker"
	@ssh $(CONNECTION) "docker rm -f vdovin-auth-container"
	@ssh $(CONNECTION) "docker run -d --name vdovin-auth-container $(if $(filter $(MODE),debug),-p 2345:2345,) -p 8080:8080 vdovin-auth:latest"

# Usage: make remote-deploy MODE=release
.PHONY: remote-deploy
remote-deploy: docker-image-tar upload-image-tar run-remote-docker
