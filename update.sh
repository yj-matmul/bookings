#!/bin/bash

git pull

soda migrate

go build -o bookings cmd/web/*.go