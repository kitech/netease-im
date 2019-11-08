package netease

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/url"
	"path/filepath"
	"strings"
)

const (
	upfileBase64Point    = neteaseBaseURL + "/msg/upload.action"
	upfileMultipartPoint = neteaseBaseURL + "/msg/fileUpload.action"
	delfilePoint         = neteaseBaseURL + "/job/nos/del.action"
	MaxFileSize          = 15 * 1024 * 1024
)

func (c *ImClient) UpfileBase64(fpath string, expiresec int) (string, error) {
	ftag := "" // md5(basename(fpath))
	tagbin := md5.Sum([]byte(filepath.Base(fpath)))
	ftag = hex.EncodeToString(tagbin[:])
	fcc, err := ioutil.ReadFile(fpath)
	if err != nil {
		return "", err
	}
	fccb64 := base64.StdEncoding.EncodeToString(fcc)
	prmval := url.Values{
		"type":    {""},
		"ishttps": {"true"},
		"tag":     {ftag},
		// "expireSec": {fmt.Sprintf("%d", math.MaxInt32)}, // 最大允许的值
	}
	if expiresec > 0 {
		prmval.Set("expireSec", fmt.Sprintf("%d", expiresec))
	}
	prmval.Add("content", fccb64)

	client := c.client.R()
	c.setCommonHead(client)
	client.SetFormDataFromValues(prmval)

	resp, err := client.Post(upfileBase64Point)
	if err != nil {
		return "", err
	}

	var jsonRes map[string]*json.RawMessage
	jsonTool.Unmarshal(resp.Body(), &jsonRes)

	var code int
	err = json.Unmarshal(*jsonRes["code"], &code)
	if err != nil {
		return "", err
	}
	if code != 200 {
		return "", errors.New(string(resp.Body()))
	}

	furl := string(*jsonRes["url"])
	furl = strings.Trim(furl, "\"")
	return furl, nil
}

func (c *ImClient) UpfileMultipart(fpath string, expiresec int) (string, error) {
	ftag := "" // md5(basename(fpath))
	tagbin := md5.Sum([]byte(filepath.Base(fpath)))
	ftag = hex.EncodeToString(tagbin[:])
	prmval := url.Values{
		"type":    {""},
		"ishttps": {"true"},
		"tag":     {ftag},
		// "expireSec": {fmt.Sprintf("%d", math.MaxInt32)}, // 最大允许的值
	}
	if expiresec > 0 {
		prmval.Set("expireSec", fmt.Sprintf("%d", expiresec))
	}

	client := c.client.R()
	c.setCommonHead(client)
	client.SetFormDataFromValues(prmval)
	client.SetFile("content", fpath)

	resp, err := client.Post(upfileMultipartPoint)
	if err != nil {
		return "", err
	}

	var jsonRes map[string]*json.RawMessage
	jsonTool.Unmarshal(resp.Body(), &jsonRes)

	var code int
	err = json.Unmarshal(*jsonRes["code"], &code)
	if err != nil {
		return "", err
	}
	if code != 200 {
		return "", errors.New(string(resp.Body()))
	}

	furl := string(*jsonRes["url"])
	furl = strings.Trim(furl, "\"")
	return furl, nil
}

func (c *ImClient) Delfile(fpath string, filets int) error {
	ftag := "" // md5(basename(fpath))
	tagbin := md5.Sum([]byte(filepath.Base(fpath)))
	ftag = hex.EncodeToString(tagbin[:])
	prmval := url.Values{
		// "contentType": {""},
		"tag": {ftag},
	}

	prmval.Set("startTime", fmt.Sprintf("%d000", filets-3*24*3600))
	prmval.Set("endTime", fmt.Sprintf("%d000", filets+3*24*3600))

	client := c.client.R()
	c.setCommonHead(client)
	client.SetFormDataFromValues(prmval)

	resp, err := client.Post(delfilePoint)
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
