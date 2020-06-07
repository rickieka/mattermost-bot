# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOBENCH=$(GOTEST) -bench=. ./... -benchmem
GOGET=$(GOCMD) get
BUILD_PATH=build
BINARY_NAME=mattermost-bot

all: clean build testcoverage
build: 
		GOOS=linux GOARCH=amd64 \
		go build -ldflags="-X main.Version=$(shell date +%FT%T%Z)" \
		-o $(BUILD_PATH)/$(BINARY_NAME) \
		main.go

test: 
		$(GOTEST) -v ./...

# generates and prints a code coverage report
testcoverage:
		mkdir -p $(BUILD_PATH)/reports

		# generate code coverage report
		$(GOTEST) ./... -coverprofile=$(BUILD_PATH)/reports/code_coverage.out

		# print coverage report
		go tool cover -func=$(BUILD_PATH)/reports/code_coverage.out

		# generate coverage report as html
		go tool cover -html=$(BUILD_PATH)/reports/code_coverage.out -o $(BUILD_PATH)/reports/code_coverage.html

# executes all benchmark functions
benchmark: 
		$(GOBENCH)

clean: 
		$(GOCLEAN)
		rm -rf $(BUILD_PATH)

runBuild:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BUILD_PATH)/$(BINARY_NAME)
