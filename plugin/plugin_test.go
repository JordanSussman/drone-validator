// Copyright 2019 Drone.IO Inc. All rights reserved.
// Use of this source code is governed by the Polyform License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/drone/drone-go/drone"
	"github.com/drone/drone-go/plugin/validator"
)

// empty context
var noContext = context.Background()

// mock github token
const mockToken = "d7c559e677ebc489d4e0193c8b97a12e"

var req = &validator.Request{
	Build: drone.Build{
		After: "3d21ec53a331a6f037a91c368710b99387d012c1",
	},
	Repo: drone.Repo{
		Slug:   "octocat/hello-world",
		Config: ".drone.yml",
	},
}

func testOutput(plugin validator.Plugin, filePath, expected string) func(*testing.T) {
	return func(t *testing.T) {
		config, err := ioutil.ReadFile(filePath)
		if err != nil {
			t.Error(err)
			return
		}
		req.Config.Data = string(config)
		actual := plugin.Validate(noContext, req)
		if actual == nil && expected == "" {
			return
		}
		if actual.Error() != expected {
			t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
		}
	}
}

func TestPlugin(t *testing.T) {
	var sensitiveImages = map[string][]string{
		"secret/data/docker/password": {
			"alpine",
			"ubuntu",
		},
		"secret/data/docker/username": {
			"alpine",
			"ubuntu",
		},
	}

	plugin := New(sensitiveImages)

	err := plugin.Validate(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	t.Run("authorized secret from env", testOutput(plugin, "testdata/authorized-env.yml", ""))
	t.Run("authorized secret from setting", testOutput(plugin, "testdata/authorized-setting.yml", ""))
	t.Run("unauthorized secret from env", testOutput(plugin, "testdata/unauthorized-env.yml", "validator: step test utilizing unauthorized secret of foo"))
	t.Run("unauthorized secret from setting", testOutput(plugin, "testdata/unauthorized-setting.yml", "validator: step test utilizing unauthorized secret of foo"))
	t.Run("authorized env secret with unauthorized command", testOutput(plugin, "testdata/unauthorized-env-command.yml", "validator: step build not authorized to utilize commands"))
	t.Run("authorized env secret with unauthorized commands", testOutput(plugin, "testdata/unauthorized-env-commands.yml", "validator: step build not authorized to utilize commands"))
	t.Run("authorized setting secret with unauthorized command", testOutput(plugin, "testdata/unauthorized-setting-command.yml", "validator: step build not authorized to utilize commands"))
	t.Run("authorized setting secret with unauthorized commands", testOutput(plugin, "testdata/unauthorized-setting-commands.yml", "validator: step build not authorized to utilize commands"))
}
