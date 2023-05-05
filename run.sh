#!/usr/bin/env bash

go build -o bookings cmd/web/*.go && ./bookings -production=false -cache=false -dbuser=postgres -dbname=bookings -dbpassword=postgres
