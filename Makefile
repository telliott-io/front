.PHONY: clean build push deploy all

clean:
	rm -rf build

build:
	GOOS=linux GOARCH=amd64 go build -o build/front ./cmd/front
	docker build -t telliottio/front:latest ./build -f Dockerfile

push: build
	docker push telliottio/front:latest

deploy:
	kubectl apply -f deployment.yaml
	kubectl apply -f ingress.yaml

all: clean push deploy