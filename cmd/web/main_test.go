package main

import (
	"testing"
)

func TestRun(t *testing.T) {
	// dbInfoPath = "./../../static/db_info.txt"
	// dsn := loadDsn(dbInfoPath)
	_, err := run()
	if err != nil {
		t.Error("fail run()")
	}
}
