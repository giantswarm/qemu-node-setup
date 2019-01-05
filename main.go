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
	dnsServersDelimiter = ","
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

	if err != nil {
		return microerror.Mask(err)
	}

	bridgeIP := flag.String("bridge-ip", "", "IP address of the bridge (used to retrieve ip and gateway).")
	dnsServers := flag.String("dns-servers", "", "Colon separated list of DNS servers.")
	hostname := flag.String("hostname", "", "Hostname of the tenant node.")
	mainConfig := flag.String("main-config", "", "Path to main ignition config (appended to small).")
	out := flag.String("out", "", "Path to save resulting ignition config.")
	flag.Parse()

	dnsServersList := strings.Split(*dnsServers, dnsServersDelimiter)
	if len(dnsServersList) == 0 {
		return microerror.New("dns servers list can not be empty")
	}

	ip := net.ParseIP(*bridgeIP)
	ip = ip.To4()
	if ip == nil {
		return microerror.New("bridge-ip should be a valid IP address")
	}

	ifaceIP := dupIP(ip)
	ifaceIP[3]++

	mainConfigData, err := ioutil.ReadFile(*mainConfig)
	if err != nil {
		return microerror.Mask(err)
	}

	nodeSetup := NodeSetup{
		DNSServers: dnsServersList,
		Gateway:    ip.String(),
		Hostname:   *hostname,
		IfaceIP:    ifaceIP.String(),
		MainConfig: base64.StdEncoding.EncodeToString(mainConfigData),
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
