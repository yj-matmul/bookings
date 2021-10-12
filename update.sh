#!/bin/bash

git pull

soda migrate

go build -o bookings cmd/web/*.go

nohup sh run.sh 1>/dev/null 2>&1 &
