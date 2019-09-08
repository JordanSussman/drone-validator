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

func TestPlugin(t *testing.T) {
	req := &validator.Request{
		Build: drone.Build{
			After: "3d21ec53a331a6f037a91c368710b99387d012c1",
		},
		Repo: drone.Repo{
			Slug:   "octocat/hello-world",
			Config: ".drone.yml",
		},
	}

	sensitiveImages := map[string][]string{
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

	goodConfig, err := ioutil.ReadFile("testdata/good.yml")
	if err != nil {
		t.Error(err)
		return
	}

	badConfig, err := ioutil.ReadFile("testdata/bad.yml")
	if err != nil {
		t.Error(err)
		return
	}

	req.Config.Data = string(goodConfig)
	err = plugin.Validate(noContext, req)
	if err != nil {
		t.Error(err)
		return
	}

	req.Config.Data = string(badConfig)
	expected := "validator: step test utilizing unauthorized secret of foo"
	actual := plugin.Validate(noContext, req)
	if actual.Error() != expected {
		t.Errorf("Error actual = %v, and Expected = %v.", actual, expected)
	}
}
