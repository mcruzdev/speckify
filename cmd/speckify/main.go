package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mcruzdev/speckify/internal/generator"
	"github.com/mcruzdev/speckify/internal/model"
	"github.com/mcruzdev/speckify/internal/oapiwrapper"

	"github.com/spf13/cobra"
)

var (
	openApiFile string
)

func main() {
	if err := CreateRootCmd().Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func CreateRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "speckify",
		Short: "Generate Wiremock StubRules from OpenAPI specification",
	}

	rootCmd.AddCommand(CreateGenerateCmd())

	return rootCmd
}

func CreateGenerateCmd() *cobra.Command {

	cmd := cobra.Command{
		Use:     "generate",
		Aliases: []string{"s"},
		Short:   "GenerateStubRules the Wiremock definition from OpenAPI specification file",
		Run: func(cmd *cobra.Command, args []string) {
			openapi3 := oapiwrapper.Get(openApiFile)
			stubRules := generator.GenerateStubRules(openapi3)
			mappings := model.NewT(stubRules)
			marshal, err := json.Marshal(mappings)
			if err != nil {
				panic(any(err))
			}

			err = os.Mkdir("mappings", 0777)
			if err != nil && os.IsNotExist(err) {
				panic(any(err))
			}

			err = os.WriteFile("mappings/wiremock-mappings.json", marshal, 0644)
			if err != nil {
				panic(any(err))
			}
		},
	}

	cmd.PersistentFlags().StringVar(&openApiFile, "oapi", "", `The OpenAPI specification file.
Usage: chorumemock generate --oapi=my-openapi-spec.yaml`)

	return &cmd
}
