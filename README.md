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

Run the db migrations ([golang-migrate](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate) is needed for this):
```shell script
make migrate-up
```

Once it is built and the migrations applied, then the bot can be run at minimum with 
the `--token` or `-t` flag:
```shell script
./bin/pikamon --token <TOKEN_HERE>
./bin/pikamon -t <TOKEN_HERE>
```

To acquire the dev bot token, reach out to a project admin. If a standalone bot
is desired (i.e. a completely separate instance), then 
follow the [Discord API Docs](https://discord.com/developers/docs/intro).

## Docker
A Docker image containing the bot and migration utility can be built using
```shell script
make docker-image
```

To run the image, pass the token in via the `PIKAMON_TOKEN` environment variable:
```shell script
docker run -e PIKAMON_TOKEN=<TOKEN_HERE> pikamon
```

To make the store persistent, just mount `/pikamon/data` to an external volume:
```shell script
docker run -e PIKAMON_TOKEN=<TOKEN_HERE> -v data:/pikamon/data pikamon
```

## Configuration
Pikamon requires a config file named `pikamon.yml` in either `/etc/pikamon` or `./configs`.

## Architecture
Pikamon is a Bot built around [github.com/bwmarrin/discordgo](https://github.com/bwmarrin/discordgo), following a model
in which handlers are registered as listeners to the discord api. When the event a handler
is listening for is pushed from the api, then that handler is called with that event and it's corresponding session.

### Handlers
There are two handlers that listen for message create events: `commands` and `spawn`. The `commands` handler simply
parses non-Pikamon messages for command text, and if found executes a call chain of commands. The `spawn` handler
listens for no-Pikamon messages to spawn Pokemon (plans for items in the future). These handlers have access to
the other components shared by the bot (the [cache and store](#data-storage)).

### Data Storage
Additionally, there is a cache used for tracking temporary events such as which pokemon are spawned in different
channels and a store that persists caught pokemon by a user.
