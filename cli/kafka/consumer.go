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

package kafka

import (
	"fmt"
	"os"
	"os/signal"

	kafkactl "github.com/jbvmio/kafka"
)

func cgHandler(msg *kafkactl.Message) (bool, error) {
	fmt.Printf("[%v] %v > %s\n", msg.Partition, msg.Offset, msg.Value)
	return true, nil
}

func launchCG(groupID string, topics ...string) {
	consumer, err := client.NewConsumerGroup(groupID, topics...)
	if err != nil {
		closeFatal("Error creating consumer group: %v\n", err)
	}
	for _, topic := range topics {
		consumer.GET(topic, cgHandler)
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	go consumer.Consume()
	<-signals
	consumer.Close()
}
