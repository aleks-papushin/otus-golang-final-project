.PHONY: docker/start-server
docker/start-server:
	docker build -t stat-server .
	docker run -p 50000:50000 stat-server

.PHONY: start-server
start-server:
	go build -o .build/server ./cmd/server
	chmod +x .build/server
	.build/server

.PHONY: start-client
start-client:
	go build -o .build/client ./cmd/client
	chmod +x .build/client
	.build/client
