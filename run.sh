#!/bin/bash

go build -o bookings cmd/web/*.go
./bookings -dbname=bookings -dbuser=postgres -dbpassword= -production=false -cache=false
