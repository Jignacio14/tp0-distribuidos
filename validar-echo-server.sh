#!/bin/bash

MESSAGE= "Hola probando 1 2 3, echooooooo"

EXECUTION=$( docker run --rm --network="tp0_testing_net" busybox:latest sh -c "echo $MESSAGE | nc server 12345 " )

SERVER_PORT = 12345

if [ "$EXECUTION" == "$MESSAGE" ]; then
    echo "action: test_echo_server | result: success"
else
    echo "action: test_echo_server | result: fail"
fi