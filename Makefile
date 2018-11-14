SHA=$(shell git rev-parse HEAD)
all:
	@echo "Building server"
	go build -o main
	@echo "Built server"
release:
	@echo "Building binary for linux..."
	CGO_ENABLED=0 \
	GOOS=linux \
	go build -v -a \
	-ldflags '-extldflags "-static"' \
	-installsuffix cgo \
	-o frontend .
	@echo "Built binary"
	@echo "Building docker container"
	docker build -t "kirandasika30/au-hacks-landing" -f Dockerfile .
	docker push "kirandasika30/au-hacks-landing:latest"
release_sha:
	@echo "Building binary for linux..."
	GOOS=linux go build -v -a --ldflags '-extldflags	"-static"' -tags netgo -installsuffix netgo -o homepage
	@echo "Built binary"
	@echo "Building docker container"
	docker build -t "kirandasika30/homepage:$(SHA)" -f min.Dockerfile .
	docker push "kirandasika30/homepage:$(SHA)"
clean:
	rm main
kube:
	@echo "Running deployment"
	kubectl apply -f ./config/deployment.yaml 
	@echo "Running service"
	kubectl apply -f ./config/service.yaml 
	@echo "Running ingress service"
	kubectl apply -f ./config/ingress.yml
