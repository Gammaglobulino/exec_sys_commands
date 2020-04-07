package main

import (
	"../TCP_scanner"
	"fmt"
	"testing"
	"time"
)

/*
func TestPortScanner_Execute(t *testing.T) {
	//it takes 1 second at least to complete one TCP dial
	scanner:=NewScanner(IP_ADDR,1,30)
	scanner.Execute()
	assert.EqualValues(t,2,len(scanner.succeededScans))
}
*/
func TestPortScannerTimeout(t *testing.T) {
	scanner := main.NewScanner(main.IP_ADDR, 1, 30)
	scanner.Timeout = time.Millisecond * 200 // if it fail, consider to rise up the timeout
	scanner.Execute()
	nports := len(scanner.SucceededScans)
	if nports != 2 {
		t.Fatalf("Expected 2 open ports, acutal %d", nports)
	}
}

func TestAsyncPortScanner_Execute(t *testing.T) {
	scanner := main.NewAsyncScanner(main.IP_ADDR, 1, 30)
	scanner.Execute()
	nports := len(scanner.SucceededScans)
	if nports != 2 {
		t.Fatalf("Expected 2 open ports, acutal %d", nports)
	}

}

func TestAsyncPortScannerWithChan_Execute(t *testing.T) {
	scanner := main.NewAsyncChanScanner(main.IP_ADDR, 1, 100)
	fmt.Println(scanner.SucceededScans)
	scanner.Execute()
	nports := len(scanner.SucceededScans)
	if nports != 3 {
		t.Fatalf("Expected 3 open ports, acutal %d", nports)
	}
}
