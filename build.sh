#!/bin/bash

CGO_ENABLED=0 GOOS=linux go build -o portainerUtility
strip portainerUtility
