package utils

import (
	"bytes"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	polycubePath string = "/polycube/v1/"
	firewallPath string = "firewall/"
)

func CreateFirewall(ip string) bool {
	resp, err := http.Post("http://"+ip+":9000"+polycubePath+firewallPath+"fw", "application/json", nil)
	if err != nil {
		log.Infoln("Could not create firewall:", err, resp)
		return false
	}
	return true
}

func AttachFirewall(ip string) bool {
	var jsonStr = []byte(`{"cube":"fw", "port":"eth0"}`)
	resp, err := http.Post("http://"+ip+":9000"+polycubePath+"attach", "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Infoln("Could not attach firewall:", err, resp)
		return false
	}
	return true
}
