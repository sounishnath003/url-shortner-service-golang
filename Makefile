AppName=url_shortner_service

install:
	go mod tidy
	go mod download -x
	go mod verify

build:
	ls -GFlash
	go build -ldflags "-s -w" -o bin/$(AppName) cmd/*.go

run: build
	./bin/$(AppName)

db:
	docker-compose -f docker.compose.yaml down
	docker-compose -f docker.compose.yaml up --build

docker-build:
	docker images
	docker rmi -f $(AppName)
	docker build -t $(AppName) -f Dockerfile .

docker-run:
	docker images
	docker run -ti $(AppName)