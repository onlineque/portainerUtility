package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"portainerUtility/api"
	"strconv"
)

// createStackCmd represents the createStack command
var createStackCmd = &cobra.Command{
	Use:   "createStack",
	Short: "creates Portainer stack from supplied arguments",
	Long:  `Creates Portainer stack from supplied arguments.`,

	Run: func(cmd *cobra.Command, args []string) {
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			log.Fatalf("can't get name parameter")
		}
		log.Printf("creating stack %s", name)

		env, err := cmd.Flags().GetStringToString("env")
		if err != nil {
			log.Fatalf("can't get env parameter")
		}

		additionalFiles, err := cmd.Flags().GetStringArray("additional-files")
		if err != nil {
			log.Fatalf("can't get additional-files parameter")
		}

		composeFile, err := cmd.Flags().GetString("compose-file")
		if err != nil {
			log.Fatalf("can't get compose-file parameter")
		}

		repositoryUrl, err := cmd.Flags().GetString("repository-url")
		if err != nil {
			log.Fatalf("can't get repository-url parameter")
		}

		repositoryReferenceName, err := cmd.Flags().GetString("repository-reference-name")
		if err != nil {
			log.Fatalf("can't get repository-reference-name parameter")
		}

		repositoryUsername, err := cmd.Flags().GetString("repository-username")
		if err != nil {
			log.Fatalf("can't get repository-username parameter")
		}

		repositoryPassword, err := cmd.Flags().GetString("repository-password")
		if err != nil {
			log.Fatalf("can't get repository-password parameter")
		}

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

		autoUpdateInterval, err := cmd.Flags().GetString("auto-update-interval")
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

		_, result := pApi.GetStackByName(name)
		if result {
			log.Printf("skipping the stack creation, stack %s already exists", name)
			return
		}

		autoUpdate := api.PortainerStackAutoUpdate{
			Interval: autoUpdateInterval,
		}

		envVariables := []api.PortainerStackEnv{}

		for key, val := range env {
			envVariables = append(envVariables, api.PortainerStackEnv{
				Name:  key,
				Value: val,
			})
		}

		err = pApi.CreateStack(api.PortainerStack{
			Name:                     name,
			AdditionalFiles:          additionalFiles,
			ComposeFile:              composeFile,
			RepositoryUrl:            repositoryUrl,
			RepositoryReferenceName:  repositoryReferenceName,
			RepositoryUsername:       repositoryUsername,
			RepositoryPassword:       repositoryPassword,
			RepositoryAuthentication: true,
			EndpointId:               float64(endpointId),
			SwarmId:                  swarmId,
			FromAppTemplate:          false,
			TlsSkipVerify:            tlsSkipVerify,
			AutoUpdate:               autoUpdate,
			Env:                      envVariables,
		})
		if err != nil {
			log.Fatalf("failed creating the stack %s: %s", name, err.Error())
		}

		log.Printf("stack %s has been successfully created", name)
	},
}

func init() {
	rootCmd.AddCommand(createStackCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createStackCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createStackCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	var envFlag map[string]string

	createStackCmd.Flags().String("name", "", "Stack name")
	createStackCmd.MarkFlagRequired("name")
	createStackCmd.Flags().String("compose-file", "docker-compose.yaml", "Docker compose filename")
	createStackCmd.MarkFlagRequired("compose-file")
	createStackCmd.Flags().StringArray("additional-files", []string{}, "Additional files needed from the repository")
	createStackCmd.Flags().String("repository-url", "", "GIT repository URL to pull docker-compose from")
	createStackCmd.MarkFlagRequired("repository-url")
	createStackCmd.Flags().String("repository-reference-name", "refs/heads/main", "GIT repository reference name, e.g.: refs/heads/main")
	createStackCmd.MarkFlagRequired("repository-reference-name")
	createStackCmd.Flags().String("repository-username", "", "GIT repository username")
	createStackCmd.MarkFlagRequired("repository-username")
	createStackCmd.Flags().String("repository-password", "", "GIT repository password")
	createStackCmd.MarkFlagRequired("repository-password")
	createStackCmd.Flags().StringToStringVar(&envFlag, "env", map[string]string{}, "environment variables")
	createStackCmd.Flags().String("auto-update-interval", "5m", "auto-update interval - default is 5m")
}
