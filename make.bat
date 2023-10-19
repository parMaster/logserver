go build -v ./app/...
@echo off

IF "%1"=="run" GOTO RUN
IF "%1"=="test" GOTO TEST
GOTO END

:RUN
logserver
GOTO END

:TEST
go test -v ./...
GOTO END

:END