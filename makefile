BUILDTARGET=tucowDemo

format:
	go fmt ./...

test:
	go test ./... -coverprofile cover.out
	go tool cover -func cover.out

vendorize:
	go mod vendor

build: format vendorize
	go build -mod=vendor -o $(BUILDTARGET)

run:
	docker-compose up --build

