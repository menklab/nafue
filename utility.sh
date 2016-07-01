#!/bin/bash

COMMAND="$1"

function deps() {
         rm -rf vendor
         go get ./...
         govendor init
         govendor add +external
}

if [ "$COMMAND" = "deps" ]; then
           echo "manage deps"
           deps
fi


