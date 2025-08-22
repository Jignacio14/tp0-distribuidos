#!/bin/bash

MESSAGE="Hola probando 1 2 3, echooooooo"
PORT=12345
EXECUTION=$(docker run --rm --network="tp0_testing_net" busybox:latest sh -c "echo $MESSAGE | nc server $PORT " )
RESULT="fail"

if [ "$EXECUTION" == "$MESSAGE" ]; then
    RESULT="success"
fi

echo "action: test_echo_server | result: $RESULT"