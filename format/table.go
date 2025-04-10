package format

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func CreateTableFromPackages(jsonResult map[string]interface{}) {
	fmt.Println()
	fmt.Println()

	packages, ok := jsonResult["packages"].([]map[string]interface{})
	if !ok {
		fmt.Println("Error: Failed to parse packages data")
		return
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Package", "Version", "Status"})
	table.SetBorder(true)
	table.SetRowLine(true)
	table.SetColumnSeparator("|")
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetAutoWrapText(false)

	for _, pkg := range packages {
		packageName, _ := pkg["package"].(string)
		version, _ := pkg["version"].(string)
		status, _ := pkg["status"].(string)

		statusText := status
		if status == "available" {
			statusText = color.GreenString(status)
		} else if status == "vulnerable" {
			statusText = color.RedString(status)
		}

		table.Append([]string{packageName, version, statusText})
	}

	table.Render()
	color.Green("\nScan completed successfully.")
}