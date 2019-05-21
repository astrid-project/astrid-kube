package utils

import (
	"bytes"
	"encoding/json"
	"net/http"

	k8sfirewall "github.com/SunSince90/polycube/src/components/k8s/utils/k8sfirewall"

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

	if !changeDefaultForward(ip) {
		return false
	}
	return setAsync(ip)
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

func setAsync(ip string) bool {
	jsonStr := []byte(`true`)
	client := http.Client{}

	req, err := http.NewRequest("PATCH", "http://"+ip+":9000"+polycubePath+firewallPath+"fw/interactive", bytes.NewBuffer(jsonStr))
	if err != nil {
		log.Infoln("Could not set firewall as asynchronous", err, req)
		return false
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Infoln("Could not set firewall as asynchronous", err, resp)
		return false
	}
	defer resp.Body.Close()
	return true
}

func DemoFakeDropAll(ips map[string]string) {
	marshal := func(rule k8sfirewall.ChainRule) ([]byte, error) {
		data, err := json.MarshalIndent(&rule, "", "   ")
		if err != nil {
			log.Errorln("Cannot marshal to json:", err)
			return nil, err
		}
		return data, nil
	}

	push := func(ip, target string) {
		endPoint := "http://" + ip + ":9000/polycube/v1/firewall/fw/chain/ingress/append/"
		rule := k8sfirewall.ChainRule{
			Action: "drop",
			Src:    target,
			Dst:    ip,
		}
		data, err := marshal(rule)
		if err == nil {
			req, err := http.NewRequest("POST", endPoint, bytes.NewBuffer(data))
			req.Header.Set("Content-Type", "application/json")

			client := &http.Client{}
			_, err = client.Do(req)
			if err != nil {
				log.Errorln("Error while trying to send request:", err)
			}
		}
	}

	apply := func(ip, name string) {
		endPoint := "http://" + ip + ":9000/polycube/v1/firewall/fw/chain/ingress/apply-rules/"
		req, err := http.NewRequest("POST", endPoint, nil)
		req.Header.Set("Content-Type", "application/json")
		client := &http.Client{}

		_, err = client.Do(req)
		if err != nil {
			log.Errorln("Error while trying to apply rules:", err)
		}

		log.Infoln("Applied drop all rule in", name)
	}

	for currentIP, currentName := range ips {
		for target := range ips {
			if currentIP != target {
				push(currentIP, target)
			}
		}

		apply(currentIP, currentName)
	}
}
