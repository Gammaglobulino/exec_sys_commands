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
	sync.Mutex
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

type AsyncChanPortScanner struct {
	totalTimeout time.Duration
	AsyncPortScanner
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

func NewAsyncChanScanner(bDNs string, start int, end int) *AsyncChanPortScanner {
	return &AsyncChanPortScanner{
		time.Second,
		AsyncPortScanner{PortScanner{BaseDNS: bDNs,
			PortStart: start,
			PortEnd:   end,
			Timeout:   time.Second,
		},
		},
	}
}

func (aps *AsyncChanPortScanner) Execute() {
	pchan := make(chan string)
	for i := aps.PortStart; i <= aps.PortEnd; i++ {
		fmt.Println(i)
		go func(ch chan string, port int) {
			if port == aps.PortEnd {
				defer close(ch)
			}
			addr := fmt.Sprintf(aps.BaseDNS+":%d", port)
			//TODO remove that log
			fmt.Println(addr)
			conn, err := net.DialTimeout("tcp", addr, aps.Timeout)
			if err != nil {
				return
			}
			ch <- addr
			conn.Close()
		}(pchan, i)
	}
	for p := range pchan {
		aps.SucceededScans = append(aps.SucceededScans, p)
	}
}

func (aps *AsyncPortScanner) Execute() {
	var wg sync.WaitGroup

	for i := aps.PortStart; i <= aps.PortEnd; i++ {
		wg.Add(1)
		fmt.Println(i)
		go func(group *sync.WaitGroup, port int) {
			defer group.Done()
			addr := fmt.Sprintf(aps.BaseDNS+":%d", port)
			//TODO remove that log
			fmt.Println(addr)
			conn, err := net.DialTimeout("tcp", addr, aps.Timeout)
			if err != nil {
				return
			}
			aps.Lock()
			aps.SucceededScans = append(aps.SucceededScans, addr)
			aps.Unlock()
			conn.Close()
		}(&wg, i)
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
