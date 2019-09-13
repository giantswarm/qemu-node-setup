package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"text/template"

	"github.com/giantswarm/microerror"
)

const (
	serversDelimiter = ","
	gateway          = "169.254.1.1"
)

func dupIP(ip net.IP) net.IP {
	dup := make(net.IP, len(ip))
	copy(dup, ip)
	return dup
}

func main() {
	err := mainError()
	if err != nil {
		panic(fmt.Sprintf("%#v", err))
	}
}

func mainError() error {
	var err error

	nodeIP := flag.String("node-ip", "", "IP address of the node")
	dnsServers := flag.String("dns-servers", "", "Colon separated list of DNS servers.")
	hostname := flag.String("hostname", "", "Hostname of the tenant node.")
	mainConfig := flag.String("main-config", "", "Path to main ignition config (appended to small).")
	ntpServers := flag.String("ntp-servers", "", "Colon separated list of NTP servers.")
	out := flag.String("out", "", "Path to save resulting ignition config.")
	flag.Parse()

	var dnsServersList []string
	{
		if len(*dnsServers) == 0 {
			return microerror.New("dns servers list can not be empty")
		}
		for _, x := range strings.Split(*dnsServers, serversDelimiter) {
			dnsServersList = append(dnsServersList, strings.TrimSpace(x))
		}
	}

	var ntpServersList []string
	{
		if len(*ntpServers) > 0 {
			for _, x := range strings.Split(*ntpServers, serversDelimiter) {
				ntpServersList = append(ntpServersList, strings.TrimSpace(x))
			}
		}
	}

	mainConfigData, err := ioutil.ReadFile(*mainConfig)
	if err != nil {
		return microerror.Mask(err)
	}

	nodeSetup := NodeSetup{
		DNSServers: dnsServersList,
		Gateway:    gateway,
		Hostname:   *hostname,
		IfaceIP:    *nodeIP,
		MainConfig: base64.StdEncoding.EncodeToString(mainConfigData),
		NTPServers: ntpServersList,
	}

	f, err := os.Create(*out)
	if err != nil {
		return microerror.Mask(err)
	}
	defer f.Close()

	t := template.Must(template.New("nodeSetup").Parse(smallIgnition))
	err = t.Execute(f, nodeSetup)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
