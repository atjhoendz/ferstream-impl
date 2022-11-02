package console

import "github.com/spf13/cobra"

var consumerCmd = &cobra.Command{
	Use:   "consumer",
	Short: "start consumer",
	Run:   consumer,
}

func init() {
	rootCmd.AddCommand(consumerCmd)
}

func consumer(_ *cobra.Command, _ []string) {

}
