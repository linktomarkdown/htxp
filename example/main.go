package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/linktomarkdown/htxp"
)

func main() {
	// 示例：生成各种随机字符串
	fmt.Println("生成订单号:", htxp.GenerateOrderNo())
	fmt.Println("生成订单ID:", htxp.GenerateOrderID("wechat"))
	fmt.Println("生成随机名称:", htxp.GenerateName(10))
	fmt.Println("生成随机密码:", htxp.GenerateRandomPassword(12, true, true, true))
	fmt.Println("MD5加密:", htxp.Md5V("password123"))

	// 示例：Minio操作
	// 首先需要初始化Minio客户端
	// err := htxp.InitMinio("localhost:9000", "minioadmin", "minioadmin", false)
	// if err != nil {
	//     log.Fatal(err)
	// }
	//
	// // 创建bucket
	// err = htxp.CreateBucket("mybucket")
	// if err != nil {
	//     log.Fatal(err)
	// }

	// 示例：RabbitMQ操作
	// 首先需要初始化RabbitMQ连接
	// rabbitConf := htxp.RabbitConf{
	//     Username: "guest",
	//     Password: "guest",
	//     Host:     "localhost",
	//     Port:     5672,
	//     VHost:    "/",
	// }
	// err := htxp.InitRabbitMQ(rabbitConf)
	// if err != nil {
	//     log.Fatal(err)
	// }
	//
	// // 初始化交换机和队列
	// err = htxp.InitRabbitMQExchange("test_exchange", "direct")
	// if err != nil {
	//     log.Fatal(err)
	// }
	//
	// // 发送消息
	// err = htxp.SendMessage("test_exchange", "test_key", []byte("Hello RabbitMQ"), "text/plain")
	// if err != nil {
	//     log.Fatal(err)
	// }

	// 示例：MDM操作
	// 首先需要初始化MDM客户端
	// err := htxp.InitMDM("localhost:9000", "minioadmin", "minioadmin", false)
	// if err != nil {
	//     log.Fatal(err)
	// }
	//
	// // 创建用户
	// err = htxp.CreateUser("newuser", "newpassword")
	// if err != nil {
	//     log.Fatal(err)
	// }

	// HTTP服务器示例
	http.HandleFunc("/success", func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"name": "张三",
			"age":  30,
		}
		htxp.Success(w, data)
	})

	http.HandleFunc("/error", func(w http.ResponseWriter, r *http.Request) {
		err := errors.New("发生了一个错误")
		htxp.Error(w, err)
	})

	http.HandleFunc("/generate", func(w http.ResponseWriter, r *http.Request) {
		// 生成随机数据
		randomData := map[string]interface{}{
			"orderNo":    htxp.GenerateOrderNo(),
			"orderID":    htxp.GenerateOrderID("alipay"),
			"randomName": htxp.GenerateName(8),
			// "randomUUID": htxp.GenerateRandomUUID(),
			"md5Hash":    htxp.Md5V("test123"),
		}
		htxp.Success(w, randomData)
	})

	fmt.Println("服务器启动在 :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
