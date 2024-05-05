compose-up:
	docker-compose up -d

compose-down:
	docker-compose down

docker-build:
	docker build -t docker.io/artemchikus/find-ship:latest -f Dockerfile .
	docker build -t docker.io/artemchikus/find-ship-migrations:latest -f .\storage\Dockerfile .
	docker build -t docker.io/artemchikus/find-ship-update-urls:latest -f .\cmd\cron_lobs\update_urls\Dockerfile .

docker-push:
	docker push docker.io/artemchikus/find-ship:latest
	docker push docker.io/artemchikus/find-ship-migrations:latest
	docker push docker.io/artemchikus/find-ship-update-urls:latest

helm-install:
	helm upgrade --install find-ship .\helm

helm-uninstall:
	helm uninstall find-ship

service-url:
	minikube service find-ship --url

build:
	go build -o find-ship cmd/main.go

run: build
	docker run -e POSTGRES_PASSWORD=password -e POSTGRES_USER=postgres -e POSTGRES_DB=db -d -p 5432:5432 postgres
	./find-ship -conf test/conf.yaml

swagger:
	cd cmd/app/api
	swag init --parseDependency  -generalInfo router.go --parseInternal

busybox:
	kubectl run temporary --image=busybox:latest -i --tty