# build stage
FROM golang:alpine AS build-env
RUN apk add --update gcc make musl-dev
WORKDIR /src/
COPY . ./
RUN make project-utils
RUN make all

# final stage
FROM alpine
RUN apk add make
WORKDIR /pikamon
COPY Makefile ./
COPY migrations ./migrations/
COPY configs /etc/pikamon/
COPY --from=build-env /src/bin/* ./
COPY --from=build-env /go/bin/migrate /bin/
RUN mkdir -p data/sqlite
ENTRYPOINT make migrate-up && ls data/sqlite && ./pikamon -t $PIKAMON_TOKEN
