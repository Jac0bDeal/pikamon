all: clean bot bot-pi

bot:
	@echo "Building bot binary for use on local system..."
	@go build -o bin/pikamon ./cmd/pikamon

bot-pi:
	@echo "Building bot binary for use on raspbian..."
	@env GOOS=linux GOARCH=arm GOARM=7 go build -o bin/pi/pikamon ./cmd/pikamon

clean:
	@echo "Cleaning bin/..."
	@rm -rf bin/*
