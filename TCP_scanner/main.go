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

type PortScanner struct{
	baseDNS string
	portStart int
	portEnd int
	succeededScans []string
	timeout time.Duration
}

func NewScanner(bDNs string,start int, end int) *PortScanner{
	return &PortScanner{
		baseDNS:        bDNs,
		portStart:      start,
		portEnd:        end,
		timeout:time.Second,
	}

}
func (ps *PortScanner) newPortRange(start int, end int){
	ps.portStart=start
	ps.portEnd=end
}

func (ps *PortScanner) newDNS(dns string){
	ps.baseDNS=dns
}

// using decorator to extend the PortScanner adding async scan
type AsyncPortScanner struct{
	PortScanner
}

func NewAsyncScanner(bDNs string,start int, end int) *AsyncPortScanner{
	return &AsyncPortScanner{
		PortScanner{baseDNS: bDNs,
			portStart: start,
			portEnd:   end,
			timeout:   time.Second,
		},
	}
}

func (aps *AsyncPortScanner) Execute(){
	var wg sync.WaitGroup
	for i:=aps.portStart;i<=aps.portEnd;i++ {
		wg.Add(1)
		go func(){
			defer wg.Done()
			addr := fmt.Sprintf(aps.baseDNS+":%d", i)
			//TODO remove that log
			fmt.Println(addr)
			conn, err := net.DialTimeout("tcp", addr,aps.timeout)
			if err != nil {
				return
			}
			aps.succeededScans=append(aps.succeededScans,addr)
			conn.Close()
		}()
	}
	wg.Wait()

}

func (ps *PortScanner) Execute() {
	for i:=ps.portStart;i<=ps.portEnd;i++ {
		addr := fmt.Sprintf(ps.baseDNS+":%d", i)
		//TODO remove that log
		fmt.Println(addr)
		conn, err := net.DialTimeout("tcp", addr,ps.timeout)
		if err != nil {
			continue
		}
		ps.succeededScans=append(ps.succeededScans,addr)
		conn.Close()
	}
}

const(
	IP_ADDR="scanner.nmap.org"
)

func main(){

}
