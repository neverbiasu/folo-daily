package handlers

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func GenerateReport() {
	files, err := ioutil.ReadDir("data")
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	reportContent := fmt.Sprintf("# Daily Report - %s\n\n", time.Now().Format("2006-01-02"))
	for _, file := range files {
		if file.IsDir() || (file.Name()[0:7] != "webhook" && file.Name()[0:5] != "data_") {
			continue
		}

		content, err := ioutil.ReadFile(fmt.Sprintf("data/%s", file.Name()))
		if err != nil {
			fmt.Println("Error reading file:", err)
			continue
		}

		reportContent += fmt.Sprintf("## %s\n\n%s\n\n", file.Name(), string(content))
	}

	dateFolder := fmt.Sprintf("data/%s", time.Now().Format("20060102"))
	_ = os.MkdirAll(dateFolder, os.ModePerm)
	reportFileName := fmt.Sprintf("%s/daily_report.md", dateFolder)
	ioutil.WriteFile(reportFileName, []byte(reportContent), 0644)
	fmt.Printf("Daily report generated: %s\n", reportFileName)
}
