#!/bin/bash
set -e

go build cli.go
go build client.go
go build server.go

mv cli DockerStuff/server
mv server DockerStuff/server
mv client DockerStuff/client

touch DockerStuff/client/ClientID
cd DockerStuff
docker-compose build
cd ..
