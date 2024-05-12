package request

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	"github.com/sreioi/jkbd/log"
)

type QueryParams map[string]string

func Get(url string, params QueryParams) map[string]interface{} {
	client := resty.New()
	resp, err := client.R().SetQueryParams(params).Get(url)

	if err != nil {
		log.Log.Printf("[GET][%s]请求数据出错：%s \n", url, err.Error())
	}

	// 解析响应数据获取总页数
	var respData map[string]interface{}
	err = json.Unmarshal(resp.Body(), &respData)
	if err != nil {
		log.Log.Printf("[GET][%s]解析响应数据出错:%s \n", url, err.Error())
	}

	return respData
}
