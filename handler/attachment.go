package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/idtoken"
)

type dto struct {
	Attachment *discordgo.MessageAttachment
	CreatedAt  time.Time
	ChannelID  string
	UserID     string
	UserName   string
}

type payload struct {
	ContentType string `json:"content_type"`
	CreatedAt   int64  `json:"created_at"`
	FileName    string `json:"file_name"`
	ProxyURL    string `json:"proxy_url"`
	URL         string `json:"url"`
	ChannelID   string `json:"channel_id"`
	UserID      string `json:"user_id"`
	UserName    string `json:"user_name"`
}

var (
	url = fmt.Sprintf("https://us-central1-%v.cloudfunctions.net/SaveItem", os.Getenv("PROJECT_ID"))
)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func OnRecieveAttachments(s *discordgo.Session, m *discordgo.MessageCreate) {
	ctx := context.Background()
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if len(m.Attachments) == 0 {
		return
	}
	saveAttachments(ctx, m)
}

func saveAttachments(ctx context.Context, m *discordgo.MessageCreate) {
	for _, a := range m.Attachments {
		a := a
		t, err := m.Timestamp.Parse()
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot parse timestamp: %v", err)
			continue
		}
		dto := &dto{
			Attachment: a,
			CreatedAt:  t,
			ChannelID:  m.ChannelID,
			UserID:     m.Author.ID,
			UserName:   m.Author.Username,
		}
		err = saveAttachment(ctx, dto)
		if err != nil {
			fmt.Fprintf(os.Stderr, "cannot save attachment %v", err)
			continue
		}
	}
}

func saveAttachment(ctx context.Context, dto *dto) error {
	p := payload{
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
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("cannnot read status code %v body %v", res.StatusCode, err)
		}
		return fmt.Errorf("status: %v msg: %v", res.StatusCode, string(body))
	}
	return nil
}
