package main

import (
	"../TCP_scanner"
	"fmt"
	"github.com/stretchr/testify/assert"
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
	scanner.Timeout = time.Millisecond * 500 // timeout sets to 500 millisec
	scanner.Execute()
	assert.EqualValues(t, 2, len(scanner.SucceededScans))

}

func TestAsyncPortScanner_Execute(t *testing.T) {
	scanner := main.NewAsyncScanner(main.IP_ADDR, 1, 30)
	scanner.Execute()
	assert.EqualValues(t, 2, len(scanner.SucceededScans))

}

func TestAsyncPortScannerWithChan_Execute(t *testing.T) {
	scanner := main.NewAsyncChanScanner(main.IP_ADDR, 1, 100)
	fmt.Println(scanner.SucceededScans)
	scanner.Execute()
	assert.EqualValues(t, 3, len(scanner.SucceededScans))
	fmt.Println(scanner.SucceededScans)

}
