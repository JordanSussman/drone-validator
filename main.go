// Copyright 2019 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"net/http"

	"github.com/JordanSussman/drone-validator/plugin"
	"github.com/drone/drone-go/plugin/validator"
	"gopkg.in/yaml.v2"

	_ "github.com/joho/godotenv/autoload"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"
)

// spec provides the plugin settings.
type spec struct {
	Bind       string `envconfig:"DRONE_BIND"`
	Debug      bool   `envconfig:"DRONE_DEBUG"`
	Secret     string `envconfig:"DRONE_SECRET"`
	ConfigFile string `envconfig:"DRONE_CONFIG_FILE"`
}

// Config contains the secrets and allowed images
type Config struct {
	Secrets map[string][]string
}

func main() {
	spec := new(spec)
	err := envconfig.Process("", spec)
	if err != nil {
		logrus.Fatal(err)
	}

	if spec.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
	if spec.Secret == "" {
		logrus.Fatalln("missing secret key")
	}
	if spec.Bind == "" {
		spec.Bind = ":3000"
	}

	yamlFile, err := ioutil.ReadFile(spec.ConfigFile)
	if err != nil {
		logrus.Fatal(err)
	}

	var config Config
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		logrus.Fatal(err)
	}

	handler := validator.Handler(
		spec.Secret,
		plugin.New(
			config.Secrets,
		),
		logrus.StandardLogger(),
	)

	logrus.Infof("server listening on address %s", spec.Bind)

	http.Handle("/", handler)
	logrus.Fatal(http.ListenAndServe(spec.Bind, nil))
}
