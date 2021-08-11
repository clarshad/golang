package cmd

import (
	"github.com/clarshad/golang/terraform-service/server"
	"github.com/spf13/cobra"
)

var (
	port    int
	rootCmd = &cobra.Command{
		Use:   "terraform-service-cli",
		Short: "Starts terraform service to accept HTTP request",
		Run: func(cmd *cobra.Command, args []string) {
			server.Handle(port)
		},
	}
)

func init() {
	rootCmd.Flags().IntVarP(&port, "port", "p", 8080, "Port to allow HTTP request")
}

func Execute() error {
	cobra.CheckErr(rootCmd.Execute())
	return nil
}
