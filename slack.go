package audiy_importer

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/slack-go/slack"
)

type Handler struct {
	conf SlackConf
}

func NewHandler(conf SlackConf) *Handler {
	return &Handler{conf: conf}
}

func (sh *Handler) Interactive(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Receive HandleSlackInteractive")
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to read body %v", err), http.StatusInternalServerError)
		return
	}
	msg, err := interactionCallbackParse(body)
	if err != nil {
		log.Printf("failed to read request body: %v", err)
		http.Error(w, fmt.Sprintf("failed to parse slack interactive payload %v", err), http.StatusInternalServerError)
		return
	}

	switch msg.Type {
	case slack.InteractionTypeMessageAction:
		payload := &IssueDialogPayload{
			TriggerID:   msg.TriggerID,
			ChannelID:   msg.Channel.ID,
			ChannelName: msg.Channel.Name,
			TeamID:      msg.Team.ID,
			TeamDomain:  msg.Team.Domain,
			MessageTs:   msg.MessageTs,
			UserID:      msg.User.ID,
			MessageText: msg.Message.Text,
			File: File{
				ID:                 msg.Message.Msg.Files[0].ID,
				Title:              msg.Message.Msg.Files[0].Title,
				Name:               msg.Message.Msg.Files[0].Name,
				URLPrivateDownload: msg.Message.Msg.Files[0].URLPrivateDownload,
				Created: int64(msg.Message.Msg.Files[0].Created),
				Mimetype:           msg.Message.Msg.Files[0].Mimetype,
			},
		}
		if err = sh.run(payload); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		log.Println("Nothing to do")
	}

	w.WriteHeader(http.StatusOK)
}

func interactionCallbackParse(reqBody []byte) (*slack.InteractionCallback, error) {
	var req slack.InteractionCallback
	jsonStr, err := url.QueryUnescape(string(reqBody)[8:])
	if err != nil {
		return nil, fmt.Errorf("error query unescape interaction callback: Body: %s | Err: %s", reqBody, err)
	}

	err = json.Unmarshal([]byte(jsonStr), &req)
	log.Printf("req is %+v", req)
	if err != nil {
		return nil, fmt.Errorf("error parsing interaction callback: Body: %s | Err: %s", reqBody, err)
	}
	return &req, nil
}

func (sh *Handler) run(payload *IssueDialogPayload) error {
	ctx := context.Background()
	//slackSvc := slacksvc.NewClient(sh.conf.SlackBotToken)
	////fileName := filepath.Base(payload.File.URLPrivateDownload)
	//b := bytes.Buffer{}
	//err := slackSvc.Download(payload.File.URLPrivateDownload, &b)
	//if err != nil {
	//	log.Printf("failed to get a slack file err=%v", err)
	//	return err
	//}
	//gcsClient, err := gcssvc.NewGCSClient(ctx, sh.conf.GCSInputAudioBucket)
	//if err != nil {
	//	log.Printf("failed to read gcs client")
	//	return err
	//}

	//ext := filepath.Ext(payload.File.Name)
	//if err := gcsClient.Put(ctx, fmt.Sprintf("%s%s",payload.File.ID, ext), b.Bytes()); err != nil {
	//	log.Printf("failed to put gcs client")
	//	return err
	//}

	message := &AudioEnqueueMessage{
		ID: payload.File.ID,
		Title: payload.File.Title,
		Name: payload.File.Name,
		Created: payload.File.Created,
		Mimetype: payload.File.Mimetype,
		URLPrivateDownload: payload.File.URLPrivateDownload,
	}

	pubsub, err := NewClient(os.Getenv("GCP_PROJECT"), os.Getenv("TOPIC_NAME"))
	if err != nil {
		return err
	}

	log.Printf("message %+v", message)
	if err := pubsub.PublishAudioCreateMessage(ctx, message); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
