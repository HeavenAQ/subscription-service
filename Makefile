BINARY_NAME=myapp
DSN="host=localhost port=5432 user=postgres password=password dbname=postgres sslmode=disable"
REDIS="127.0.0.1:6379"

build:
	@echo "Building..."
	env CGO_ENABLED=0 go build -ldflags "-s -w" -o ${BINARY_NAME} -v ./cmd/web/
	@echo "Built"

run: build
	@echo "Starting..."
	@echo env DSN=${DSN} REDIS=${REDIS} ./${BINARY_NAME} &
	@echo "Started!"

clean:
	@echo "Cleaning..."
	@go clean
	@rm -f ${BINARY_NAME}
	@echo "Cleaned!"

start: run

stop:
	@echo "Stopping"
	@-pkill -SIGTERM -f ./${BINARY_NAME}

restart: stop start

test:
	go test -v ./...
