#!/usr/bin/env bash

docker build -t local/test/sdcli .
docker build -t local/test/sdclitests -f test.Dockerfile .
