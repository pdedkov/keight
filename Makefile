SERVICE=keight

CMD?=api
FLAGS?=-v

default: lint test

build:
	go build $(SERVICE)/cmd/$(CMD)

build-race:
	go build -race $(FLAGS) $(SERVICE)/cmd/$(CMD)

test:
	go test $(FLAGS) -count=1 ./...

test-race:
	go test -race $(FLAGS) ./...

integration:
	godotenv -f .env go test $(FLAGS) -count=1 -tags=integration ./...

gen:
	go generate $(FLAGS) ./...

lint:
	golangci-lint run -v ./...

.PHONY: build build-race test test-race lint gen integration
