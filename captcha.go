package main

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/xkeyideal/captcha/pool"

	"github.com/gin-gonic/gin"
)

/*
* demo
 */
const HTML_Content string = `
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title></title>
</head>

<body>
<form method="post" action="/v1/captcha">
    <p><img src="{{.imgbase64}}" /></p>
    <p><input name="captcha" placeholder="请输入验证码" type="text" /></p>
    <input name="captcha_id" type="hidden" value="{{.CaptchaId}}" />
	 <input name="captcha_id" type="hidden" value="{{.Val}}" />
    <input type="submit" />
</form>
</body>
</html>`

var CaptchaPool *pool.CaptchaPool

func GetCaptcha(c *gin.Context) {
	cacheBuffer := CaptchaPool.GetImage()

	base_url := base64.StdEncoding.EncodeToString(cacheBuffer.Data.Bytes())
	html := strings.Replace(HTML_Content, "{{.CaptchaId}}", cacheBuffer.Id, -1)
	html = strings.Replace(html, "{{.Val}}", string(cacheBuffer.Val), -1)
	final_html := strings.Replace(html, "{{.imgbase64}}", "data:image/png;base64,"+base_url, -1)
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, final_html)
}

func main() {
	CaptchaPool = pool.NewCaptchaPool(240, 80, 6, 2, 2, 2)

	var router *gin.Engine

	router = gin.Default()

	router.GET("/v1/captcha", GetCaptcha)

	if err := router.Run(fmt.Sprintf("0.0.0.0:%s", "7312")); err != nil {
		fmt.Println("Eagle Eye Server Run Failed: ", err.Error())
		os.Exit(-1)
	}

	CaptchaPool.Stop()
}
