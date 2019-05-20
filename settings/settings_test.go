package settings

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshalConf(t *testing.T) {
	path := "./conf.yaml"

	Load(path)

	assert.True(t, len(Settings.Paths.Kubeconfig) > 0)
	fmt.Printf("%+v\n", Settings)
}
