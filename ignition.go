package main

const smallIgnition = `
{
  "ignition": {
    "version": "2.2.0",
    "config": {
      "append": [
        {
	    "source": "data:text/plain;charset=utf-8;base64,{{ .MainConfig }}"
        }
      ]
    }
   },
  "networkd": {
    "units": [
      {
        "contents": "[Match]\nName=eth0\n\n[Network]\nAddress={{ .IfaceIP }}/32\n{{range $srv := .DNSServers }}\nDNS={{ $srv }}{{ end }}{{range $srv := .NTPServers }}\nNTP={{ $srv }}{{ end }}\n[Route]Destination={{.Gateway}}/32\nScope=link\n[Route]\nGateway={{.Gateway}}",
        "name": "00-eth0.network"
      }
    ]
  },
  "storage": {
    "files": [{
      "filesystem": "root",
      "path": "/etc/hostname",
      "mode": 420,
      "contents": { "source": "data:,{{ .Hostname }}" }
    }]
  }
}
`

type NodeSetup struct {
	DNSServers []string
	Gateway    string
	Hostname   string
	IfaceIP    string
	MainConfig string
	NTPServers []string
}
