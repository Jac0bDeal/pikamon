all: clean build

build:
	go build -o bin/pikamon ./cmd/pikamon

clean:
	rm bin/*

rebuild: clean build
