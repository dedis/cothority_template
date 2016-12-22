#!/usr/bin/env bash

DBG_SHOW=1
# Debug-level for app
DBG_APP=1
DBG_SRV=1
# Uncomment to build in local dir
#STATICDIR=test
# Needs 4 clients
NBR_CLIENTS=4

. $GOPATH/src/github.com/dedis/cothority/libcothority/cothority.sh

main(){
    startTest
#	test Build
#	test Time
	test Count
    stopTest
}

testCount(){
	runCoBG 1 2
	testFail runApp counter
	testOK runApp counter group.toml
	testGrep ": 0" runApp counter group.toml
	runApp time group.toml
	testGrep ": 1" runApp counter group.toml
}

testTime(){
	runCoBG 1 2
	testFail runApp time
	testOK runApp time group.toml
	testGrep Time runApp time group.toml
}

testBuild(){
    testOK runApp --help
}

runApp(){
    dbgRun ./$APP -d $DBG_APP $@
}

main
