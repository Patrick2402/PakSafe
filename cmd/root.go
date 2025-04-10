package cmd

import (
	"fmt"
	"log"
	"os"

	"strings"

	"github.com/fatih/color"

	table "PakSafe/format"
	"PakSafe/internal/dependencies"
	"PakSafe/internal/registry"
	"PakSafe/utils"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "paksafe",
	Short: "PakSafe - A tool for scanning package dependencies",
	Long:  `A tool that scans package dependencies for vulnerabilities and other issues.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
func init() {
    var onlyList bool
    var scan bool
    var format string
    var outputFile string  // Changed from defaultPath for clarity

    var npmCmd = &cobra.Command{
        Use:   "npm [package.json path]",
        Short: "Scan npm dependencies",
        Long:  `Scan npm dependencies for vulnerabilities and other issues.`,
        Run: func(cmd *cobra.Command, args []string) {
            defaultPath := "package.json"
            if len(args) > 0 {
                defaultPath = args[0]
            }

            deps, err := dependencies.GetDependencies(defaultPath)
            if err != nil {
                log.Fatalf("Error during downloading dependencies: %v", err)
            }

            if len(deps) == 0 {
                fmt.Println("Haven't found any dependencies")
                return
            }

            if onlyList {
                utils.ListDependencies(deps)
                return
            }

            if scan || (!onlyList && !scan) {
                color.Cyan("Scanning dependencies...")
                registryStatuses := registry.CheckNpmRegistry(deps)
                jsonResult := dependencies.BuildJson(deps, registryStatuses)

                switch strings.ToLower(format) {
                case "json":
                    resultPath := "./result.json" 
                    if outputFile != "" {
                        resultPath = outputFile
                    }
                    
                    err = dependencies.SaveJsonToFile(jsonResult, resultPath)
                    if err != nil {
                        log.Fatalf("Error during saving to JSON: %v", err)
                    }
                    fmt.Printf("Results saved to: %s\n", resultPath)

                case "table", "terminal", "console":
                    table.CreateTableFromPackages(jsonResult)
                default:
                    table.CreateTableFromPackages(jsonResult)
                }
            }
        },
    }

    npmCmd.Flags().BoolVar(&onlyList, "only-list", false, "Only list dependencies without scanning")
    npmCmd.Flags().BoolVar(&scan, "scan", false, "Scan dependencies for vulnerabilities") 
    npmCmd.Flags().StringVar(&format, "format", "table", "Output format: json or table")
    npmCmd.Flags().StringVar(&outputFile, "output-file", "", "Output file path (default: ./result.json)")

    rootCmd.AddCommand(npmCmd)
}