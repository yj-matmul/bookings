#!/bin/bash

app/bookings/bookings -dbname=bookings -dbuser=postgres -dbpassword= -production=false -cache=false -logpath=/app/bookings/logs/application_log.txt
