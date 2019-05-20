package settings

import (
	"io/ioutil"
	osUser "os/user"

	"github.com/SunSince90/ASTRID-kube/types"
	"gopkg.in/yaml.v2"

	log "github.com/sirupsen/logrus"
)

var (
	Settings types.Settings
)

func Load(path string) {
	conf, err := ioutil.ReadFile(path)
	if err != nil {
		log.Panic("Could not load configuration:", err)
	}

	settings := types.Settings{}

	//	Parse the YAML
	err = yaml.Unmarshal(conf, &settings)
	if err != nil {
		log.Panic("Could not parse configuration file:", err)
	}

	//	kubeconfig is empty?
	if len(settings.Paths.Kubeconfig) < 1 {
		settings.Paths.Kubeconfig = loadDefaultKubeconfigPath()
	}

	Settings = settings
}

func loadDefaultKubeconfigPath() string {
	currentUser, err := osUser.Current()
	if err != nil {
		log.Panic("Could not get current user's name")
	}

	return "/home/" + currentUser.Name + "/.kube/config"
}
