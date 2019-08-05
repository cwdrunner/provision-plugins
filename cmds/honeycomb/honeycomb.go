package main

//go:generate drbundler content content.go
//go:generate drbundler content content.yaml
//go:generate sh -c "drpcli contents document content.yaml > honeycomb.rst"
//go:generate rm content.yaml

import (
	"fmt"
	"os"

	"github.com/digitalrebar/logger"
	"github.com/digitalrebar/provision/v4/api"
	"github.com/digitalrebar/provision/v4/models"
	"github.com/digitalrebar/provision/v4/plugin"
	"github.com/honeycombio/libhoney-go"
	"github.com/digitalrebar/provision-plugins/v4"
)

var (
	version = v4.RS_VERSION
	def     = models.PluginProvider{
		Name:          "honeycomb",
		Version:       version,
		PluginVersion: 2,
		HasPublish:    true,
		RequiredParams: []string{
			"honeycomb/writekey",
		},
		OptionalParams: []string{
			"honeycomb/dataset",
		},
		Content: contentYamlString,
	}
)

type Plugin struct {
	session *api.Client
}

func (p *Plugin) Config(l logger.Logger, session *api.Client, config map[string]interface{}) (err *models.Error) {
	writekey, ok := config["honeycomb/writekey"].(string)
	if !ok {
		err = &models.Error{Code: 400, Model: "plugin", Key: "honeycomb", Type: "rpc", Messages: []string{"Bad write key"}}
		return
	}

	dataset, ok := config["honeycomb/dataset"].(string)
	if !ok {
		dataset = "digitalrebar"
	}

	honeyconfig := libhoney.Config{
		WriteKey: writekey,
		Dataset:  dataset,
	}
	libhoney.Init(honeyconfig)
	libhoney.UserAgentAddition = fmt.Sprintf("rackn/%s", v4.RS_VERSION)
	return
}

func (p *Plugin) Publish(l logger.Logger, e *models.Event) (err *models.Error) {

	ev := libhoney.NewEvent()
	defer ev.Send()

	// build up object by field instead of relying on generic JSON ev.add(e);
	ev.AddField("Time", e.Time)
	ev.AddField("Action", e.Action)
	ev.AddField("Type", e.Type)
	ev.AddField("Key", e.Key)
	ev.Add(e.Object)

	switch {
	case e.Type == "tbd":
		ev.AddField("foo", "bar")
	default:
		// nothing
	}

	return
}

func main() {
	plugin.InitApp("honeycomb", "Sends events as specified honeycomb event.", version, &def, &Plugin{})
	err := plugin.App.Execute()
	if err != nil {
		os.Exit(1)
	}
}