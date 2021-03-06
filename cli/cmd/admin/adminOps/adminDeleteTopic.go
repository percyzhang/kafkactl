package adminops

import (
	"github.com/jbvmio/kafkactl/cli/kafka"
	"github.com/jbvmio/kafkactl/cli/x/out"

	"github.com/spf13/cobra"
)

var cmdAdminDeleteTopic = &cobra.Command{
	Use:     "topic",
	Aliases: []string{"topics"},
	Short:   "Delete Kafka Topics",
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch true {
		case cmd.Flags().Changed("out"):
			out.Warnf("Error: Cannot use --out when deleting topics.")
			return
		default:
			kafka.DeleteTopics(args...)
		}
	},
}

func init() {
}
