go build -v ./cmd/logserver
@echo off

IF "%1"=="run" GOTO RUN
GOTO END

:RUN
logserver.exe
GOTO END

:END