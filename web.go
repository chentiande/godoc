package main

import (
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.StaticFS("/assets", http.Dir("assets"))
	r.StaticFS("/monacoeditorwork", http.Dir("monacoeditorwork"))
	r.StaticFS("/static", http.Dir("static"))
	r.StaticFile("/", "index.html")

	r.Any("/api/*any", func(c *gin.Context) {
		// 获取来自地址 A 的请求信息
		req := c.Request
		// 创建一个新的 HTTP 客户端
		client := &http.Client{}

		// 构建请求，转发给地址 B
		newReq, err := http.NewRequest(req.Method, "http://127.0.0.1:8888/"+c.Request.URL.Path[5:]+"?"+c.Request.URL.RawQuery, req.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 复制地址 A 的请求头到新的请求
		newReq.Header = make(http.Header)
		for h, val := range req.Header {
			newReq.Header[h] = val
		}

		// 发送请求到地址 B
		resp, err := client.Do(newReq)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer resp.Body.Close()

		// 读取地址 B 返回的响应内容
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		// 将地址 B 返回的响应内容作为地址 A 的响应内容
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	})

	r.Run(":8080") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
