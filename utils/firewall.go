package utils

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	firewall_path string = "/polycube/v1/firewall/"
)

func CreateFirewall(ip string) bool {
	resp, err := http.Post("http://"+ip+":9000"+firewall_path+"fw", "application/json", nil)
	if err != nil {
		log.Infoln("Could not create firewall:", err, resp)
		return false
	}

	fmt.Println("ok, response for:"+ip, resp)
	return true
}
