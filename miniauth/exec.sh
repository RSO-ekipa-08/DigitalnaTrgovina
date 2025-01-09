#!/usr/bin/env bash
docker build -t auth0-golang-web-app .
docker run --env-file .env -p 50051:50051 -it auth0-golang-web-app
