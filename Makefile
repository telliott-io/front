.PHONY: clean build push deploy all

clean:
	rm -rf build

build:
	GOOS=linux GOARCH=amd64 go build -o build/front ./cmd/front
	- rm -rf build/public
	- rm -rf build/views
	cp -r public build/public
	cp -r views build/views
	docker build -t telliottio/front:tilt ./build -f Dockerfile

deploy:
	kubectl apply -k deployment
	kubectl rollout restart deployment/front --namespace front

dashboard:
	echo "{\"overwrite\":true,\"dashboard\":\
	`jsonnet -J thirdparty/grafonnet-lib/ observability/dashboard.jsonnet`\
	}" \
	| \
	curl -X POST -H "Content-Type: application/json" -d @- \
	http://admin:secret@grafana.telliott.io/api/dashboards/db

all: clean push deploy