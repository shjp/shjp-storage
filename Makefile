build-functions:
	mkdir -p functions
	go get ./...
	go clean -cache
	go build -o functions/storage ./cmd/netlify-function
