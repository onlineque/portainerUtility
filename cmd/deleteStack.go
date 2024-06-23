package cmd

import (
	"fmt"
	"log"
	"portainerUtility/api"
	"strconv"

	"github.com/spf13/cobra"
)

// deleteStackCmd represents the deleteStack command
var deleteStackCmd = &cobra.Command{
	Use:   "deleteStack",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatalf("can't get name parameter")
		}
		log.Printf("deleting stack %s", name)

		portainerUrl, err := cmd.Flags().GetString("portainer-url")
		if err != nil {
			log.Fatalf("can't get portainer-url parameter")
		}
		portainerApiKey, err := cmd.Flags().GetString("portainer-api-key")
		if err != nil {
			log.Fatalf("can't get portainer-api-key parameter")
		}
		endpointId, err := strconv.Atoi(cmd.Flag("endpoint-id").Value.String())
		if err != nil {
			log.Fatalf("can't get endpoint-id parameter")
		}

		swarmId, err := cmd.Flags().GetString("swarm-id")
		if err != nil {
			log.Fatalf("can't get swarm-id parameter")
		}
		tlsSkipVerify, err := cmd.Flags().GetBool("tls-skip-verify")
		if err != nil {
			log.Fatalf("can't get tls-skip-verify parameter")
		}

		baseUrl := fmt.Sprintf("%s/api", portainerUrl)
		pApi := api.PortainerAPI{
			ApiBaseUrl:         baseUrl,
			ApiKey:             portainerApiKey,
			EndpointId:         endpointId,
			SwarmId:            swarmId,
			InsecureSkipVerify: tlsSkipVerify,
		}
		foundStack, result := pApi.GetStackByName(name)
		if !result {
			log.Fatalf("can't find the stack %s by name", name)
		}
		err = pApi.DeleteStack(foundStack)
		if err != nil {
			log.Fatalf("can't delete the stack %s: %s", name, err.Error())
		}
		log.Printf("stack %s has been successfully deleted", name)
	},
}

func init() {
	rootCmd.AddCommand(deleteStackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteStackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteStackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	deleteStackCmd.Flags().String("name", "", "Stack name to be deleted")
	deleteStackCmd.MarkFlagRequired("name")
}
