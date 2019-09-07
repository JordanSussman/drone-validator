// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"bytes"
	"context"
	"fmt"

	"github.com/drone/drone-go/plugin/validator"
	"github.com/drone/drone-yaml/yaml"
)

// New returns a new validator plugin.
func New(images []string) validator.Plugin {
	return &plugin{
		images: images,
	}
}

type plugin struct {
	images []string
}

// helper function to determine if string exists in slice of strings
func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func parsePipeline(pipeline *yaml.Pipeline, p *plugin) error {
	// loop over all pipeline steps
	for _, step := range pipeline.Steps {
		// validate if image being utilized has commands disallowed with secrets
		if stringInSlice(step.Image, p.images) {
			// validate if command or commands stanza utilized
			if len(step.Commands) != 0 || len(step.Command) != 0 {
				// loop over all settings
				for _, setting := range step.Settings {
					// validate if secrets are defined
					if setting.Secret != "" {
						return fmt.Errorf("failing due to step %s utilizing commands", step.Name)
					}
				}
			}
		}
	}
	return nil
}

func (p *plugin) Validate(ctx context.Context, req *validator.Request) error {
	f := bytes.NewBufferString(req.Config.Data)
	m, err := yaml.Parse(f)
	if err != nil {
		return err
	}

	for _, r := range m.Resources {
		switch v := r.(type) {
		case *yaml.Pipeline:
			err := parsePipeline(v, p)
			if err != nil {
				return err
			}
		}
	}

	// a nil error indicates the configuration is valid.
	return nil
}
