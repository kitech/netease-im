package netease

import (
	"encoding/json"
	"errors"
)

const (
	addFriendPoint = neteaseBaseURL + "/friend/add.action"
)

func (c *ImClient) AddFriend(accid string, faccid string) error {
	param := map[string]string{"accid": accid, "faccid": faccid, "type": "1", "msg": "heloaf"}

	client := c.client.R()
	c.setCommonHead(client)
	client.SetFormData(param)

	resp, err := client.Post(addFriendPoint)
	if err != nil {
		return err
	}

	var jsonRes map[string]*json.RawMessage
	jsonTool.Unmarshal(resp.Body(), &jsonRes)

	var code int
	err = json.Unmarshal(*jsonRes["code"], &code)
	if err != nil {
		return err
	}

	if code != 200 {
		return errors.New(string(resp.Body()))
	}

	return nil
}
