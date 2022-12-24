package main

import (
	"fmt"
	"log"
	"os"
	"net"
	"strconv"

	"github.com/miekg/dns"
	"gopkg.in/yaml.v3"
)

type Config struct {
	PORT int `yaml:"port"`
	ADDR string `yaml:"esxi_addr"`
	Domain string `yaml:"domain_name"`
}

var config = &Config{}

func parseQuery(m *dns.Msg, records map[string]string) {
	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			log.Printf("Query for %s\n", q.Name)
			ip := records[q.Name]

			if ip == "" {
				addr, err := net.ResolveIPAddr("ip", q.Name)
				if err != nil {
					log.Printf("EROR", err)
				}
				if addr.IP != nil {
					ip = string(addr.IP)
				}
			}

			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}
		}
	}
}

func parseConfig() {
	b, err := os.ReadFile("config.yml")
	if err != nil { 
		log.Fatalln(err)		
	}

	yaml.Unmarshal(b, config)
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	
	records := map[string]string{
		config.Domain:config.ADDR,
	}

	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m, records)
	}

	w.WriteMsg(m)
}

func main() {
	parseConfig()

	// attach request handler func
	dns.HandleFunc("local.", handleDnsRequest)

	// start server
	server := &dns.Server{Addr: ":" + strconv.Itoa(config.PORT), Net: "udp"}
	log.Printf("Starting at %d\n", config.PORT)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}