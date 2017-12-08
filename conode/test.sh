#!/usr/bin/env bash

DBG_TEST=1
DBG_APP=2

. $GOPATH/src/github.com/dedis/onet/app/libtest.sh

main(){
    startTest
    setupConode
    test Build
    stopTest
}

testBuild(){
    testOK dbgRun runCo 1 --help
}

main
