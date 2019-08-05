package main

//go:generate ./download-vibs.sh
//go:generate drbundler content content.go
//go:generate drbundler content content.yaml
//go:generate sh -c "drpcli contents document content.yaml > vmware.rst"
//go:generate rm content.yaml
//go:generate go-bindata -pkg main -o embed.go -prefix embedded embedded/...

import (
	"os"

	"github.com/digitalrebar/logger"
	"github.com/digitalrebar/provision/v4/api"
	"github.com/digitalrebar/provision/v4/models"
	"github.com/digitalrebar/provision/v4/plugin"
	"github.com/digitalrebar/provision-plugins/v4"
)

var (
	version = v4.RS_VERSION
	def     = models.PluginProvider{
		Name:          "vmware",
		Version:       version,
		PluginVersion: 2,
		AutoStart:     true,
		HasPublish:    false,
		Content:       contentYamlString,
	}
)

type Plugin struct {
}

func (p *Plugin) Config(l logger.Logger, session *api.Client, data map[string]interface{}) *models.Error {
	return nil
}

func (p *Plugin) Unpack(thelog logger.Logger, path string) error {
	return RestoreAssets(path, "")
}

func main() {
	plugin.InitApp("vmware", "Add VMware vSphere and vCloud capabilities.", version, &def, &Plugin{})
	err := plugin.App.Execute()
	if err != nil {
		os.Exit(1)
	}
}