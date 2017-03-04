#!/bin/bash

CURRENT_PATH="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
gofmt -l -s -w ${CURRENT_PATH}/../..
