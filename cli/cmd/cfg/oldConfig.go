// Copyright © 2018 NAME HERE <jbonds@jbvm.io>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cfg

import (
	"io/ioutil"
	"log"
	"os"

	"github.com/jbvmio/kafkactl/cli/cx"

	"github.com/jbvmio/kafkactl/cli/x/out"
	"github.com/spf13/viper"
	yaml "gopkg.in/yaml.v2"
)

const homeConfigName = `.kafkactl.yaml`

// OldConfig contains a collection of cluster entries
type OldConfig struct {
	Current string  `json:"current" yaml:"current"`
	Entries []Entry `json:"entries" yaml:"entries"`
}

// Entry contains kafka and burrow node details for a cluster
type Entry struct {
	Name      string   `json:"name" yaml:"name"`
	Kafka     []string `json:"kafka yaml:"kafka"`
	Burrow    []string `json:"burrow" yaml:"burrow"`
	Zookeeper []string `json:"zookeeper" yaml:"zookeeper"`
}

func testForOldConfig(filePath ...string) {
	var configFilePath string
	defaultFilePath := homeDir() + `/` + homeConfigName
	switch {
	case len(filePath) > 1:
		out.Infof("Too Many Paths Specified.")
		os.Exit(0)
	case len(filePath) == 1 && filePath[0] != "":
		configFilePath = filePath[0]
	default:
		switch {
		case !fileExists(defaultFilePath):
			out.Infof("No default config file found.\n  Run kafkactl config --sample to display a sample config file.\n  Save your config in ~/.kafkactl.yaml")
			os.Exit(0)
		case fileExists(defaultFilePath):
			configFilePath = defaultFilePath
		}
	}
	v := viper.New()
	v.SetConfigFile(configFilePath)
	v.ReadInConfig()
	switch {
	case !v.InConfig("config-version") && !v.InConfig("current-context"):
		if v.InConfig("current") {
			var oldConfig OldConfig
			v.Unmarshal(&oldConfig)
			contexts := make(map[string]cx.Context, len(oldConfig.Entries))
			for _, entry := range oldConfig.Entries {
				ctx := cx.Context{
					Name:      entry.Name,
					Brokers:   entry.Kafka,
					Burrow:    entry.Burrow,
					Zookeeper: entry.Zookeeper,
				}
				contexts[entry.Name] = ctx
			}
			newConfig := Config{
				CurrentContext: oldConfig.Current,
				Contexts:       contexts,
				ConfigVersion:  configVersion,
			}
			out.Marshal(newConfig, "yaml")
			out.Infof("config [%v] has been converted to Latest.", configFilePath)
		}
	default:
		out.Infof("config [%v] at Latest.", configFilePath)
	}
}

func returnConfig(config []byte) OldConfig {
	conf := OldConfig{}
	err := yaml.Unmarshal(config, &conf)
	if err != nil {
		log.Fatalf("Error returning config: %v\n", err)
	}
	return conf
}

func readConfig(path string) []byte {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("Error reading config file: %v\n", err)
	}
	return file
}

func writeConfig(path string, config []byte) {
	err := ioutil.WriteFile(path, config, 0644)
	if err != nil {
		log.Fatalf("Error writing config: %v\n", err)
	}
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}
