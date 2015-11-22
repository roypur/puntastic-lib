#!/bin/bash
CURRENT_DIR="$(cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
echo $CURRENT_DIR
GOPATH=$GOPATH:$CURRENT_DIR && gomobile bind -target android $1
