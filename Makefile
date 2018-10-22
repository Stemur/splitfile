# Go Parameters

GOCMD=Go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOBENCH=$(GOCMD) test -bench=.
BINARY_NAME=splitfile
BINARY_WINDOWS=$(BINARY_NAME).exe

build:
		$(GOBUILD) -o $(BINARY_NAME) -v

run:
		$(GOBUILD) -o $(BINARY_NAME) -v ./...
		./$(BINARY_NAME)

test:
	$(GOTEST) -v ./...

benchmark:
	$(GOBENCH)

# Cross compilation
build-windows:
	GOOS=windows GOARCH=386 $(GOBUILD) -o $(BINARY_WINDOWS) -v
