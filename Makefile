app/ir-remote: app/main.go
	@cd app && \
	env GOOS=linux GOARCH=arm GOARM=5 go build

docker:
	docker build -t ir-remote .
.PHONY: docker

docker-run:
	docker run -p 127.0.0.1:8080:80 ir-remote
.PHONY: docker-run
