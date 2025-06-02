package handlers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type WebhookPayload struct {
	Entry map[string]interface{} `json:"entry"`
	Feed  map[string]interface{} `json:"feed"`
}

func HandleWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var payload WebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}

	dateFolder := fmt.Sprintf("data/%s", time.Now().Format("20060102"))
	_ = os.MkdirAll(dateFolder, os.ModePerm)
	fileName := fmt.Sprintf("%s/webhook_%s.md", dateFolder, time.Now().Format("150405"))
	fileContent := fmt.Sprintf("# Webhook Data\n\n## Entry\n\n%v\n\n## Feed\n\n%v\n", payload.Entry, payload.Feed)
	ioutil.WriteFile(fileName, []byte(fileContent), 0644)

	fmt.Fprintf(w, "Data saved to %s", fileName)
}
