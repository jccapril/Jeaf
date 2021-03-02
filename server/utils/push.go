package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	url = "https://msgapi.umeng.com/api/send"
	appkey = "603de189b8c8d45c1386aac0"
	app_master_secret = "zkfkjixnekpcayukfgrzofnlahvamxtf"
)



//appkey = '你的appkey'
//app_master_secret = '你的app_master_secret'
//timestamp = '你的timestamp'
//method = 'POST'
//url = 'http://msg.umeng.com/api/send'
//params = {
//	'appkey': appkey,
//	'timestamp': timestamp,
//	'device_tokens': device_token,
//	'type': 'unicast',
//	'payload': {
//		'body': {
//			'ticker': 'Hello World',
//			'title':'你好',
//			'text':'来自友盟推送',
//			'after_open': 'go_app'
//		},
//		'display_type': 'notification'
//	}
//}
//post_body = json.dumps(params)
//print post_body
//sign = md5('%s%s%s%s' % (method,url,post_body,app_master_secret))

type UMessage struct {
	Appkey 			string		`json:"appkey"`
	Timestamp 		string		`json:"timestamp"`
	Type			string		`json:"type"`
	Payload			Payload		`json:"payload"`
	ProductionMode  string   	`json:"production_mode"`
	Description		string  	`json:"description"`
}

type Payload struct {
	Aps 	APS `json:"aps"`
}

type APS struct {
	Alert 	Alert 	`json:"alert"`
}


type Alert struct {
	Title 		string	`json:"title"`
	Subtitle 	string	`json:"subtitle"`
	Body		string	`json:"body"`
}


func Push(){

	alert := Alert{Title:"test_title",Subtitle:"test_subtitle",Body:"test_body"}
	aps  := APS{alert}
	payload := Payload{aps}
	message := UMessage{
		Appkey:      	appkey,
		Timestamp:   	fmt.Sprintf("%d",time.Now().Unix()),
		Type:        	"broadcast",
		Payload:     	payload,
		ProductionMode: "false",
		Description: 	"iOS推送测试",
	}
	postBody, err := json.Marshal(message)
	if err!=nil{
		log.Fatal(err)
	}

	fmt.Printf("%s\n",postBody)
	//sign = md5('%s%s%s%s' % (method,url,post_body,app_master_secret))

	method := "POST"

	origin := fmt.Sprintf("%s%s%s%s",method,url,postBody,app_master_secret)

	sign := fmt.Sprintf("%x",md5.Sum([]byte(origin)))
	fmt.Println("origin:",origin)
	fmt.Println("signed:",sign)
	signedUrl := fmt.Sprintf("%s?sign=%s",url,sign)

	response, err := http.Post(signedUrl,"application/x-www-form-urlencoded",bytes.NewBuffer([]byte(postBody)))
	if err != nil {
		fmt.Println("err",err)
	}
	b1, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("err",err)
	}
	fmt.Println(string(b1))

}

