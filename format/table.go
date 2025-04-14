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
    table.SetHeader([]string{"Package", "Version", "Status", "Reason", "Private Latest Version", "Public Latest Version"})
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
        reason, _ := pkg["reason"].(string)
        privateVersion, _ := pkg["privateVersion"].(string)
        publicVersion, _ := pkg["publicVersion"].(string)

        statusText := status
        if status == "secure" {
            statusText = color.GreenString(status)
        } else if status == "vulnerable" {
            statusText = color.RedString(status)
        } else if status == "suspicious" {
            statusText = color.YellowString(status)
        } else if status == "not found" {
            statusText = color.CyanString(status)
        }

        table.Append([]string{packageName, version, statusText, reason, privateVersion, publicVersion})
    }

    table.Render()
    color.Green("\nScan completed successfully.")
}