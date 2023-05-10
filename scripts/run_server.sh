#!/bin/bash

SCIRPTS_DIR="$( cd -- "$(dirname "$0")" >/dev/null 2>&1 ; pwd -P )"

# Move into project root, then into front end and start node
cd $SCIRPTS_DIR && cd ../
cd server && go run main.go

# Move into project root, then into server and start go server
cd $SCIRPTS_DIR && cd ../
cd frontend && npm run start

