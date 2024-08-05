package cmd

import (
	"fmt"

	"github.com/vimek-go/server-faker/internal/pkg/api"
	"github.com/vimek-go/server-faker/internal/pkg/plugins"
	"github.com/vimek-go/server-faker/internal/pkg/transformer"

	"github.com/vimek-go/server-faker/internal/pkg/logger"
	"github.com/vimek-go/server-faker/internal/pkg/parser"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var (
	filePath     string
	serverPort   int
	url          string
	responseType string
)

const defaultPort = 8080

var rootCmd = &cobra.Command{
	Use:   "server-faker",
	Short: "Creates a fake server with ease",
	Long: `This application creates a fake server based on provided json.
Use run to start the server. 
Use parser to prepare the json file.

The default port is 8080.
Examlpe:
server-faker run --file=../test-api.json --port=8080`,
}

var serverCmd = &cobra.Command{
	Use:   "run",
	Short: "Server commands",
	Run: func(*cobra.Command, []string) {
		startServer()
	},
}

var parserCmd = &cobra.Command{
	Use:   "parse",
	Short: "Parser commands",
	Run: func(*cobra.Command, []string) {
		transformer := transformer.New()
		output, err := transformer.Transform(filePath, url, responseType)
		if err != nil {
			fmt.Println("error transforming file", err)
			return
		}
		fmt.Println(output)
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func PrepareCommand() error {
	serverCmd.Flags().StringVarP(&filePath, "file", "f", "", "[required] The file path to the json file")
	err := serverCmd.MarkFlagRequired("file")
	if err != nil {
		fmt.Println("error marking flag required")
		return err
	}
	serverCmd.PersistentFlags().IntVarP(&serverPort, "port", "p", defaultPort, "The port to run the server on")

	parserCmd.Flags().StringVarP(&filePath, "file", "f", "", "[required] The file path to the json file")
	err = parserCmd.MarkFlagRequired("file")
	if err != nil {
		fmt.Println("error marking flag required")
		return err
	}
	parserCmd.PersistentFlags().StringVarP(&url, "url", "u", "", "The url to generate the dynamic endpoint")
	parserCmd.PersistentFlags().
		StringVarP(&responseType, "type", "t", "static", "The response type of the endpoint to generate")

	rootCmd.AddCommand(serverCmd, parserCmd)
	return nil
}

func startServer() {
	logger, err := logger.NewExternalLogger("debug")
	if err != nil {
		fmt.Printf("error creating a logger %v\n", err)
		return
	}
	pluginLoader := plugins.NewPlugingLoader(logger)
	factory := parser.NewFactory(pluginLoader, logger)
	parser := parser.NewLoader(factory, logger)
	handlers, err := parser.LoadConfig(filePath)
	if err != nil {
		fmt.Println("error creating a logger")
		return
	}

	e := gin.New()
	api := api.NewBaseAPI(e, logger)
	api.AddEndpoints(handlers)

	if err := api.Run(fmt.Sprintf(":%d", serverPort)); err != nil {
		fmt.Printf("error runnig server %v\n", err)
		return
	}
}
