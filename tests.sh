#!/bin/sh

echo ++++++++ go run .
go build . && go run .
echo ++++++++ go run . reset
go build . && go run . reset
echo ++++++++ go run . reset too many
go build . && go run . reset too many
echo ++++++++ go run . register
go build . && go run . register
echo ++++++++ go run . register bun
go build . && go run . register bun
echo ++++++++ go run . register bun
go build . && go run . register bun
echo ++++++++ go run . login
go build . && go run . login
echo ++++++++ go run . login bun
go build . && go run . login bun
echo ++++++++ go run . login test
go build . && go run . login test

