#!/usr/bin/env bash

# Application
export PORT=7777

# MySQL variables
export DATABASE_FILE_PATH="./qonto_accounts.sqlite"
export DATABASE_MAX_CONN_LIFETIME="0"
export DATABASE_MAX_OPEN_CONNS="0"
export DATABASE_MAX_IDLE_CONNS="0"

# execution
./bin/service-qonto api --addr 127.0.0.1
