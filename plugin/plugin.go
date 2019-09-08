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
func New(sensitiveImages map[string][]string) validator.Plugin {
	return &plugin{
		sensitiveImages: sensitiveImages,
	}
}

type plugin struct {
	sensitiveImages map[string][]string
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

// helper function to parse the steps and determine if secret isn't allowed
func parseSteps(steps []*yaml.Container, secretPath string, secretNames []string, sensitiveImages map[string][]string) error {
	// loop over all steps
	for _, step := range steps {
		// loop over all environment variables
		for _, env := range step.Environment {
			// check if the secret is one we care about
			if stringInSlice(env.Secret, secretNames) {
				// check if the image is allowed for the image
				if !stringInSlice(step.Image, sensitiveImages[secretPath]) {
					return fmt.Errorf("validator: step %s utilizing unauthorized secret of %s", step.Name, env.Secret)
				}
				// check if the step utilizes commands
				if len(step.Commands) != 0 || len(step.Command) != 0 {
					return fmt.Errorf("validator: step %s not authorized to utilize commands", step.Name)
				}
			}
		}
		// TODO: update to utilize function instead of duplication
		// loop over all settings
		for _, setting := range step.Settings {
			// check if the secret is one we care about
			if stringInSlice(setting.Secret, secretNames) {
				// check if the image is allowed for the image
				if !stringInSlice(step.Image, sensitiveImages[secretPath]) {
					return fmt.Errorf("validator: step %s utilizing unauthorized secret of %s", step.Name, setting.Secret)
				}
				// check if the step utilizes commands
				if len(step.Commands) != 0 || len(step.Command) != 0 {
					return fmt.Errorf("validator: step %s not authorized to utilize commands", step.Name)
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

	var steps []*yaml.Container
	secrets := make(map[string][]string)
	for _, r := range m.Resources {
		switch v := r.(type) {
		case *yaml.Pipeline:
			for _, step := range v.Steps {
				steps = append(steps, step)
			}
		case *yaml.Secret:
			// get the path for the secret
			fullSecret := fmt.Sprintf("%s/%s", v.Get.Path, v.Get.Name)
			// check if the map already contains the secret
			if _, ok := secrets[fullSecret]; ok {
				// append the values to safeguard incase the secret is referenced
				// multiple times
				secrets[fullSecret] = append(secrets[fullSecret], v.Name)
			} else {
				secrets[fullSecret] = []string{v.Name}
			}
		}
	}

	// further optimize this by comparing the slice and map from the start
	for secretPath, secretNames := range secrets {
		// validate if the secret referenced is one we care about
		if _, ok := p.sensitiveImages[secretPath]; ok {
			err = parseSteps(steps, secretPath, secretNames, p.sensitiveImages)
			if err != nil {
				return err
			}
		}
	}

	// a nil error indicates the configuration is valid.
	return nil
}
