build-functions:
	./cmd/netlify-function/override-env.sh
	mkdir -p functions
	go get ./...
	go clean -cache
	go build -o functions/storage ./cmd/netlify-function
