/*
Copyright © 2023 Patrick Hermann patrick.hermann@sva.de
*/

package cli

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/atc0005/go-teams-notify/v2/messagecard"
)

type MsTeamsWebhook struct {
	Title string `yaml:"title"`
	Text  string `yaml:"text"`
	Color string `yaml:"color"`
	Url   string `yaml:"url"`
}

func SendWebhookToTeams(webhook MsTeamsWebhook) bool {

	// Initialize a new Microsoft Teams client.
	mstClient := goteamsnotify.NewTeamsClient()

	// Set webhook url.
	webhookUrl := webhook.Url

	// Setup message card.
	msgCard := messagecard.NewMessageCard()
	msgCard.Title = webhook.Title
	msgCard.Text = webhook.Text
	msgCard.ThemeColor = webhook.Color

	// Send the message with default timeout/retry settings.
	if err := mstClient.Send(webhookUrl, msgCard); err != nil {
		log.Printf("failed to send message: %v", err)
		os.Exit(1)
	}

	return true
}

func DownloadFileWithProgressBar(downloadUrl, targetDir string) {

	req, _ := http.NewRequest("GET", downloadUrl, nil)
	resp, _ := http.DefaultClient.Do(req)
	defer resp.Body.Close()

	fileName := filepath.Base(downloadUrl)

	f, _ := os.OpenFile(targetDir+"/"+fileName, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	bar := progressbar.DefaultBytes(
		resp.ContentLength,
		"downloading "+fileName,
	)
	io.Copy(io.MultiWriter(f, bar), resp.Body)

}
