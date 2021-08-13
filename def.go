package audiy_importer

type AudioEnqueueMessage struct {
	ID                 string `json:"id"`
	Title              string `json:"title"`
	Name               string `json:"name"`
	URLPrivateDownload string `json:"url_private_download"`
	Created            int64 `json:"created"`
	Mimetype           string `json:"mimetype"`
}

func (c *pubsubClient) PublishAudioCreateMessage(ctx context.Context, m *AudioEnqueueMessage) error {
	serialized, err := json.Marshal(&m)
	if err != nil {
		return err
	}
	if err := c.Publish(ctx, serialized); err != nil {
		return err
	}
	return nil
}


type IssueDialogPayload struct {
	TriggerID   string
	ChannelID   string
	ChannelName string
	TeamID      string
	TeamDomain  string
	MessageTs   string
	UserID      string
	MessageText string
	File        File
	Files       []File
}

type File struct {
	ID                 string
	Title              string
	Name               string
	URLPrivateDownload string
	Created            int64
	Mimetype           string
}
