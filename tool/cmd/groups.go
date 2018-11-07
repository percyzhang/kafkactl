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

package cmd

import (
	"log"
	"sort"
	"strings"

	"github.com/jbvmio/kafkactl"
)

func searchGroupListMeta(groups ...string) []kafkactl.GroupListMeta {
	client, err := kafkactl.NewClient(bootStrap)
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}
	defer func() {
		if err := client.Close(); err != nil {
			log.Fatalf("Error closing client: %v\n", err)
		}
	}()
	if verbose {
		client.Logger("")
	}
	glMeta, err := client.GetGroupListMeta()
	if err != nil {
		log.Fatalf("Error getting grouplist metadata: %s\n", err)
	}
	var groupListMeta []kafkactl.GroupListMeta
	for _, g := range groups {
		for _, m := range glMeta {
			if exact {
				if m.Group == g {
					groupListMeta = append(groupListMeta, m)
				}
			} else {
				if strings.Contains(m.Group, g) {
					groupListMeta = append(groupListMeta, m)
				}
			}
		}
	}
	sort.Slice(groupListMeta, func(i, j int) bool {
		return groupListMeta[i].Group < groupListMeta[j].Group
	})
	return groupListMeta
}