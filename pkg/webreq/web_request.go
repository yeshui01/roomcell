/* ====================================================================
 * Author           : tianyh(mknight)
 * Email            : 824338670@qq.com
 * Last modified    : 2021-12-18 14:16
 * Filename         : web_request.go
 * Description      :
 * ====================================================================*/
package webreq

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	ReqTimeOut = 5
)

/*
   发送Get请求
   url：         请求地址
*/
func Get(url string) (error, []byte) {
	client := &http.Client{Timeout: ReqTimeOut * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	var buffer [1024]byte
	var errResult error
	var nRead int
	result := bytes.NewBuffer(nil)
	for {
		nRead, errResult = resp.Body.Read(buffer[0:])
		result.Write(buffer[0:nRead])
		if errResult != nil && errResult == io.EOF {
			break
		} else if errResult != nil {
			break
		}
	}
	if errResult != nil && errResult != io.EOF {
		return errResult, nil
	}
	return nil, result.Bytes()
}

/*
   发送POST请求
   url：         请求地址
   data：        POST请求提交的数据
   contentType： 请求体格式，如：application/json
*/
func Post(url string, data []byte, contentType string) (error, []byte) {
	client := &http.Client{Timeout: ReqTimeOut * time.Second}
	resp, err := client.Post(url, contentType, bytes.NewBuffer(data))
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	result, errRead := ioutil.ReadAll(resp.Body)
	return errRead, result
}

/*
   发送POST请求,请求参数内容为json格式
   url：         请求地址
   data：        POST请求提交的数据,内部会自动序列化成json格式
*/
func PostJson(url string, data interface{}) (error, []byte) {
	client := &http.Client{Timeout: ReqTimeOut * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonStr))
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)
	return nil, result
}

/*
   发送POST请求, 请求参数内容为form表单
   url：         请求地址
   params        Form表单数据
*/
func PostForm(url string, params url.Values) (error, []byte) {
	// resp, err := http.PostForm(url, params)
	client := &http.Client{Timeout: ReqTimeOut * time.Second}
	resp, err := client.PostForm(url, params)
	if err != nil {
		logrus.Error("PostForm error:%s", err.Error())
		return err, nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Error("read body error:%s", err.Error())
		return err, nil
	}
	return nil, body
}
