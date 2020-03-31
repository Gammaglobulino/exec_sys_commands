package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

type Command interface {
	Execute()
}

type PortScanner struct {
	BaseDNS        string
	PortStart      int
	PortEnd        int
	SucceededScans []string
	Timeout        time.Duration
}

func NewScanner(bDNs string, start int, end int) *PortScanner {
	return &PortScanner{
		BaseDNS:   bDNs,
		PortStart: start,
		PortEnd:   end,
		Timeout:   time.Second,
	}

}
func (ps *PortScanner) newPortRange(start int, end int) {
	ps.PortStart = start
	ps.PortEnd = end
}

func (ps *PortScanner) newDNS(dns string) {
	ps.BaseDNS = dns
}

// using decorator to extend the PortScanner adding async scan
type AsyncPortScanner struct {
	PortScanner
}

func NewAsyncScanner(bDNs string, start int, end int) *AsyncPortScanner {
	return &AsyncPortScanner{
		PortScanner{BaseDNS: bDNs,
			PortStart: start,
			PortEnd:   end,
			Timeout:   time.Second,
		},
	}
}

func (aps *AsyncPortScanner) Execute() {
	var wg sync.WaitGroup

	for i := aps.PortStart; i <= aps.PortEnd; i++ {
		wg.Add(1)
		fmt.Println(i)
		go func(group *sync.WaitGroup) {
			defer group.Done()
			addr := fmt.Sprintf(aps.BaseDNS+":%d", i)
			//TODO remove that log
			fmt.Println(addr)
			conn, err := net.DialTimeout("tcp", addr, aps.Timeout)
			if err != nil {
				return
			}
			aps.SucceededScans = append(aps.SucceededScans, addr)
			conn.Close()
		}(&wg)
	}
	wg.Wait()

}

func (ps *PortScanner) Execute() {
	for i := ps.PortStart; i <= ps.PortEnd; i++ {
		addr := fmt.Sprintf(ps.BaseDNS+":%d", i)
		//TODO remove that log
		fmt.Println(addr)
		conn, err := net.DialTimeout("tcp", addr, ps.Timeout)
		if err != nil {
			continue
		}
		ps.SucceededScans = append(ps.SucceededScans, addr)
		conn.Close()
	}
}

const (
	IP_ADDR = "scanner.nmap.org"
)

func main() {

}
