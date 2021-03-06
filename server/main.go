package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"server/service"
)

type Response struct {
	Code int         	`json:"code"`
	Msg  string      	`json:"msg"`
	Data interface{} 	`json:"data"`
}
type Data struct {
	Token string `json:"token"`
}

const (
	ERROR   = 2007
	SUCCESS = 2000
)

type User struct {
	Username string		`json:"username" form:"username"`
	Password string 	`json:"password" form:"password"`
}


func main() {

	//ummessage := ummessage.Init("晚上好","晚上好晚上好晚上好晚上好晚上好晚上好晚上好晚" +
	//	"晚上好晚上好晚上好晚上好晚上好晚上好晚上好晚上好晚上好" +
	//	"晚上好晚上好晚上好晚上好晚上好上好晚上好晚上好晚上好晚上好晚上好晚上好晚上好晚上好晚上好","")
	//ummessage.ProductionMode = false
	//isOK,err := ummessage.BroadCast()
	//if isOK {
	//	fmt.Println("广播推送成功")
	//}else  {
	//	fmt.Println("广播推送失败",err)
	//}

	service.RunSSQ()

	//return

	r := gin.Default()

	r.POST("/user/login", func(context *gin.Context) {
		//username := context.Param("username")
		//buf := make([]byte, 1024)
		//n, _ := context.Request.Body.Read(buf)
		//fmt.Println(string(buf[0:n]))
		var user User
		var err error
		err = context.ShouldBind(&user)

		if err != nil {
			fmt.Println(err)
			log.Fatal(err)
		}

		username := user.Username
		password := user.Password
		fmt.Printf("username:%v password:%v \n",username,password)
		code := ERROR
		msg := "登录失败"
		var data Data
		if username == "admin" && password == "123456" {
			code = SUCCESS
			msg = "登录成功"
			data = Data{Token:"admin"}
		}
		context.JSON(http.StatusOK, Response{
			Code: code,
			Msg:  msg,
			Data: data,
		})
	})

	r.GET("user/info/:token", func(context *gin.Context) {
		token := context.Param("token")
		code := ERROR
		msg := "没有该用户"
		var data  map[string]interface{}
		fmt.Println("token:", token)
		if token == "admin" {
			code = SUCCESS
			msg = "获取用户数据成功"
			data = map[string]interface{}{
				"id|1-10000": 1,
				"name": "@cname",
				"roles": []string{"manager"},
			}
		}
		context.JSON(http.StatusOK, Response{
			Code: code,
			Msg:  msg,
			Data: data,
		})
	})

	r.Run()

}