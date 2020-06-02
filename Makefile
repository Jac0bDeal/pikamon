all: clean build

bot:
	go build -o bin/pikamon ./cmd/pikamon

clean:
	rm bin/*
