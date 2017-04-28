package main

import (
	"dtrans"
)

func main() {
	dtrans.Start()
}

// TODO rollback on abort
// TODO load testing & deadlock testing
// TODO publish results to all servers on commit
// TODO HANDLE DEADLOCK
