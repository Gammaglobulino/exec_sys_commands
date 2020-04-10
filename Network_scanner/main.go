package main

import (
	"context"
	"fmt"
	"github.com/Ullaakut/nmap"
	"log"
	"time"
)

func main() {
	targetIP := "192.168.33.113"
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(targetIP),
		nmap.WithPorts("80,443", "999"),
		nmap.WithContext(ctx),
	)
	if err != nil {
		log.Fatal(err)
	}

	results, warnings, errs := scanner.Run()
	if errs != nil {
		log.Fatal(errs)
	}
	if warnings != nil {
		log.Fatalf("Warnings: %s\n", warnings)
	}
	for _, host := range results.Hosts {
		if len(host.Ports) != 0 {
			for _, add := range host.Addresses {
				fmt.Println(add)
			}
			for _, port := range host.Ports {
				fmt.Println(port)
			}
			fmt.Println()
		}
	}

}
