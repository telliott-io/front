.PHONY: clean build push deploy all

clean:
	rm -rf build

build:
	GOOS=linux GOARCH=amd64 go build -o build/front ./cmd/front
	- rm -rf build/public
	- rm -rf build/styles
	cp -r public build/public
	cp -r views build/views
	docker build -t telliottio/front:latest ./build -f Dockerfile

push: build
	docker push telliottio/front:latest

deploy:
	kubectl apply -f deployment.yaml
	kubectl apply -f ingress.yaml
	# Trigger a rolling update
	kubectl rollout restart deployment/front

all: clean push deploy