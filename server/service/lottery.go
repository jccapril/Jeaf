package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/goinggo/mapstructure"
	"github.com/robfig/cron"
	"io/ioutil"
	"net/http"
	"server/model"
	"server/ummessage"
	"time"
)

const(
	ssq_url = "http://www.cwl.gov.cn/cwl_admin/kjxx/findDrawNotice?name=ssq&issueCount=1"
	ssq_duration = 30
)

func StartTimer(F func()) {
	go func() {
		for {

			now := time.Now()
			var next time.Time
			if now.Hour() > 18 {
				next = now.Add(time.Hour * 24)
			}else  {
				next = now
			}
			next = time.Date(next.Year(), next.Month(), next.Day(), 17, 6, 0, 0, next.Location())
			t := time.NewTimer(next.Sub(time.Now()))
			<-t.C
			fmt.Printf("开始每天定时，结算完成: %v\n",time.Now())
			F()
		}
	}()
}


func RunSSQ(){

	//i := 0
	c := cron.New()

	// 每天21点
	spec := "0 0 21 * * ?"
	c.AddFunc(spec, func() {
		ticker := time.NewTicker(time.Second * ssq_duration)
		//ch := make(chan int)
		go func() {
			for {
				select {
				case <-ticker.C:
					ssq,err := GetLatestSSQByRemote()
					if err != nil {
						fmt.Println("error:",err)
					}else {
						if ssq.Code == "2021022" {
							ticker.Stop()
							ummessage := ummessage.Init(ssq.Name+"开奖啦","",ssq.Red+"+"+ssq.Blue)
							ummessage.ProductionMode = false
							isOK,err := ummessage.BroadCast()
							if isOK {
								fmt.Println("广播推送成功")
							}else  {
								fmt.Println("广播推送失败",err)
							}
						}

					}
				}
			}
		}()
	})

	c.Start()

	//关闭着计划任务, 但是不能关闭已经在执行中的任务.
	//defer c.Stop()





}
func GetLatestSSQByRemote()(ssq model.SSQ, err error){
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet,ssq_url,nil)
	req.Header.Set("Referer", "http://www.cwl.gov.cn/")
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return
	}
	resultList,isOK := result["result"].([]interface{})
	if !isOK {
		return ssq,errors.New("format error")
	}
	if len(resultList) > 0 {
		lastResult,isOK := resultList[0].(map[string]interface{})
		if !isOK {
			return ssq,errors.New("format error")
		}
		err = mapstructure.Decode(lastResult, &ssq)
		if err != nil {
			return
		}
	}else {
		err = errors.New("result is nil")
	}
	return
}

//func httpDo() {
//	client := &http.Client{}
//	req, err := http.NewRequest(
//		"POST",
//		"http://www.01happy.com/demo/accept.php",
//		strings.NewReader("name=cjb"),
//		)
//	if err != nil {
//		// handle error
//	}
//	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
//	req.Header.Set("Cookie", "name=anny")
//	resp, err := client.Do(req)
//	defer resp.Body.Close()
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		// handle error
//	}
//	fmt.Println(string(body))
//}

//func MyRequest(remoteUrl string, queryValues url.Values) {
//	// client := &http.Client{}
//	uri, err := url.Parse(remoteUrl)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//	if queryValues != nil {
//		values := uri.Query()
//		if values != nil {
//			for k, v := range values {
//				queryValues[k] = v
//			}
//		}
//		uri.RawQuery = queryValues.Encode()
//	}
//	fmt.Println(uri.String())
//	req, err := http.NewRequest("GET", uri.String(), nil)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//	fmt.Println(req.URL)
//}



