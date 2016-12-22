#!/usr/bin/env bash

DBG_TEST=1
# Debug-level for app
DBG_APP=2

. $GOPATH/src/github.com/dedis/onet/app/libtest.sh

main(){
    startTest
    buildCothority include.go
	test Build
	test Main
    stopTest
}

testMain(){
	testGrep Main runTmpl main
}

testBuild(){
    testOK dbgRun runTmpl --help
}

runTmpl(){
    dbgRun ./$APP -d $DBG_APP $@
}

main
