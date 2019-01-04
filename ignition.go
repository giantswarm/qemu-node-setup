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
        "contents": "[Match]\nName=eth0\n\n[Network]\nAddress={{ .IfaceIP }}/30\nGateway={{ .Gateway }}\nDNS=1.1.1.1\nDNS=8.8.4.4",
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
	Gateway    string
	Hostname   string
	IfaceIP    string
	MainConfig string
}
