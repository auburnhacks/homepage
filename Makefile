all:
	@echo "Building server"
	go build -o main
	@echo "Built server"
release:
	@echo "Building binary for linux..."
	CGO_ENABLED=0 \
	GOOS=linux \
	go build -a \
	-ldflags '-extldflags "-static"' \
	-installsuffix cgo \
	-o main .
	@echo "Built binary"
	@echo "Building docker container"
	docker build -t "auburn-hacks-landing" -f Dockerfile .
