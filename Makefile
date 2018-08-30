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
	docker build -t "kirandasika30/au-hacks-landing" -f Dockerfile .
	docker push "kirandasika30/au-hacks-landing:latest"
clean:
	rm main
kube:
	@echo "Running deployment"
	kubectl apply -f ./config/deployment.yaml 
	@echo "Running service"
	kubectl apply -f ./config/service.yaml 
	@echo "Running ingress service"
	kubectl apply -f ./config/ingress.yml
