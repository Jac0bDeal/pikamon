# Pikamon
A Pokemon Discord bot inspired by the late Pokecord.

[![Go Report Card](https://goreportcard.com/badge/github.com/Jac0bDeal/pikamon)](https://goreportcard.com/report/github.com/Jac0bDeal/pikamon)

## Installation
Pikamon is built on Go 1.14+, please make sure this version is used to 
build and test the application.

Build the Pikamon bot binary using the make target `bot`:
```shell script
make bot
```

Once it is built, then the bot can be run at minimum with 
the `--token` or `-t` flag:
```shell script
./bin/pikamon --token <TOKEN_HERE>
./bin/pikamon -t <TOKEN_HERE>
```
Alternatively, if the token is exported via the appropriate env variable, 
the `--token` flag is not needed:
```shell script
export PIKAMON_TOKEN <TOKEN_HERE>
./bin/pikamon
```
To acquire the dev bot token, reach out to a project admin. If a standalone bot
is desired (i.e. a completely separate instance), then 
follow the [Discord API Docs](https://discord.com/developers/docs/intro).
