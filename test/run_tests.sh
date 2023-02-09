#!/bin/sh

docker build -t local/test/sdcli .
docker build -t local/test/sdclitests test
docker run -i local/test/sdclitests
