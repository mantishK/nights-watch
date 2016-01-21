package model

import (
	"bytes"
	"encoding/json"
	"net/http"
	"plivo/nights-watch/config"
)

//Makes external calls to send SMS
func FireSMS(src, dest, authID, text string) error {
	body := struct {
		Src    string `json:"src"`
		Dest   string `json:"dst"`
		AuthId string `json:"auth_id"`
		Text   string `json:"text"`
	}{src, dest, authID, text}

	url, _ := config.GetString("sms_url")

	bodyJson, _ := json.Marshal(body)
	buf := bytes.NewReader(bodyJson)

	resp, err := http.Post(url, "application/json", buf)
	if resp.StatusCode != 200 {
		return err
	}

	return nil
}
