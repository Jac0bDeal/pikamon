all: clean build

bot:
	go build -o bin/pikamon ./cmd/pikamon

bot-pi:
	env GOOS=linux GOARCH=arm GOARM=7 go build -o bin/pi/pikamon ./cmd/pikamon

clean:
	rm bin/*
