#!/usr/bin/env bash

DBG_TEST=1
DBG_APP=2

. ../libtest.sh

main(){
    startTest
    setupConode github.com/dedis/cothority_template/service
    run testBuild
    stopTest
}

testBuild(){
    cp co1/public.toml .
    testOK dbgRun runCo 1 --help
}

main
