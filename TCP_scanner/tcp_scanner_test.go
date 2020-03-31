package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPortScanner_Execute(t *testing.T) {
	//it takes 1 second at least to complete one TCP dial
	scanner:=NewScanner(IP_ADDR,1,30)
	scanner.Execute()
	assert.EqualValues(t,2,len(scanner.succeededScans))
}
func TestPortScannerTimeout(t *testing.T) {
	scanner:=NewScanner(IP_ADDR,1,30)
	scanner.timeout=time.Millisecond*500 // timeout sets to 500 millisec
	scanner.Execute()
	assert.EqualValues(t,2,len(scanner.succeededScans))

}
func TestAsyncPortScanner_Execute(t *testing.T) {
	scanner:=NewAsyncScanner(IP_ADDR,1,30)
	scanner.timeout=time.Millisecond*500 // timeout sets to 500 millisec
	scanner.Execute()
	assert.EqualValues(t,2,len(scanner.succeededScans))

}
