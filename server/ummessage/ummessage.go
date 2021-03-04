package ummessage

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	UMBaseApiUrl            = "https://msgapi.umeng.com/api/send"
	UMAppkey            	= "603de189b8c8d45c1386aac0"
	UMAppMasterSecret 		= "zkfkjixnekpcayukfgrzofnlahvamxtf"
	UMHttpMethod			= "POST"
)

const (
	UMTypeBroadcast	=	"broadcast"
)

type UMessage struct {
	Appkey         string  `json:"appkey"`
	Timestamp      string  `json:"timestamp"`
	Type           string  `json:"type"`
	Payload        Payload `json:"payload"`
	ProductionMode bool    `json:"production_mode"`
	Description    string  `json:"description"`
}

type Payload struct {
	Aps APS `json:"aps"`
	//  "key1":"value1",       // 可选，用户自定义内容, "d","p"为友盟保留字段，
}

type APS struct {
	Alert Alert `json:"alert"`
}

type Alert struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Body     string `json:"body"`
}

func Init(title,subtitle,body string) *UMessage{
	alert := Alert{Title: title, Subtitle: subtitle, Body: body}
	aps := APS{alert}
	payload := Payload{aps}
	message := UMessage{
		Appkey:         UMAppkey,
		Timestamp:      fmt.Sprintf("%d", time.Now().Unix()),
		Type:           UMTypeBroadcast,
		Payload:        payload,
		ProductionMode: true,
		Description:    "iOS推送测试", // 可选，发送消息描述，建议填写。
	}
	return &message
}
func (message *UMessage)BroadCast()(isOK bool,err error){
	var postBody []byte
	postBody, err = json.Marshal(message)
	if err != nil {
		return
	}

	origin := fmt.Sprintf("%s%s%s%s", UMHttpMethod, UMBaseApiUrl, postBody, UMAppMasterSecret)
	sign := fmt.Sprintf("%x", md5.Sum([]byte(origin)))
	signedUrl := fmt.Sprintf("%s?sign=%s", UMBaseApiUrl, sign)

	var response *http.Response
	response, err = http.Post(
		signedUrl,
		"application/x-www-form-urlencoded",
		bytes.NewBuffer([]byte(postBody)),
	)
	if err != nil {
		return
	}
	var result []byte
	result, err = ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	var resultMap map[string]interface{}
	err = json.Unmarshal(result, &resultMap)
	fmt.Println()
	isOK = resultMap["ret"] == "SUCCESS"
	fmt.Println(resultMap)
	return
}
