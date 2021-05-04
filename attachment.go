package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/idtoken"
)

type DTO struct {
	Attachment *discordgo.MessageAttachment
	CreatedAt  time.Time
	ChannelID  string
	UserID     string
	UserName   string
}

type Payload struct {
	ContentType string `json:"content_type"`
	CreatedAt   int64  `json:"created_at"`
	FileName    string `json:"file_name"`
	ProxyURL    string `json:"proxy_url"`
	URL         string `json:"url"`
	ChannelID   string `json:"channel_id"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
}

func saveAttachment(ctx context.Context, dto *DTO) error {
	url := fmt.Sprintf("https://us-central1-%v.cloudfunctions.net/SaveItem", os.Getenv("PROJECT_ID"))

	p := Payload{
		ContentType: "",
		CreatedAt:   dto.CreatedAt.Unix(),
		FileName:    dto.Attachment.Filename,
		ProxyURL:    dto.Attachment.ProxyURL,
		URL:         dto.Attachment.URL,
		ChannelID:   dto.ChannelID,
		UserID:      dto.UserID,
		UserName:    dto.UserName,
	}

	// json values
	values, err := json.Marshal(p)
	if err != nil {
		return err
	}

	client, err := idtoken.NewClient(ctx, url)
	if err != nil {
		return fmt.Errorf("idtoken.NewClient: %v", err)
	}

	res, err := client.Post(url, "application/json", bytes.NewBuffer(values))
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		body, error := ioutil.ReadAll(res.Body)
		if error != nil {
			log.Fatal(error)
		}
		fmt.Printf("status: %v msg: %v\n", res.StatusCode, string(body))
	}
	return nil
}
