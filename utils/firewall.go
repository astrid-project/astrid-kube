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

	return changeDefaultForward(ip)
}

func changeDefaultForward(ip string) bool {
	jsonStr := []byte(`"forward"`)
	directions := []string{"ingress", "egress"}
	client := http.Client{}

	for _, direction := range directions {
		req, err := http.NewRequest("PATCH", "http://"+ip+":9000"+polycubePath+firewallPath+"fw/chain/"+direction+"/default", bytes.NewBuffer(jsonStr))
		if err != nil {
			log.Infoln("Could not change default action in", direction, err, req)
			return false
		}
		req.Header.Set("Content-Type", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			log.Infoln("Could not change default action in", direction, err, resp)
			return false
		}
		defer resp.Body.Close()
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
