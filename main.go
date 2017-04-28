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
// TODO make so that if transaction not committed other transactions will not see new values until commit
