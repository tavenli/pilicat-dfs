package utils

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"pilicat-core/logs"
	"strings"

	"bytes"
)

func GetIP(req *http.Request) string {

	ip := req.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	ip = req.Header.Get("Remote_addr")
	if ip == "" {
		ip = req.RemoteAddr
	}

	//host, port, err = net.SplitHostPort(req.RemoteAddr)
	if Contains(ip, ":") {
		ip = Split(ip, ":")[0]
	}
	return ip
}

func HttpPostJson(url string, json string) (string, error) {
	resp, err := http.Post(url, "application/json", strings.NewReader(json))
	if err != nil {
		logs.Error("HttpPostJson error: ", err)
		return "", err
	}

	if resp == nil {
		return "", errors.New("返回对象为空")
	}

	defer resp.Body.Close()
	result := ""
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		result = string(body)
		//logs.Debug("HttpPostJson result: ", result)
	} else {
		logs.Error("HttpPostJson error: ", err)
	}

	return result, nil
}

func HttpPostJsonReturnByte(url string, json string) ([]byte, error) {
	resp, err := http.Post(url, "application/json", strings.NewReader(json))
	if err != nil {
		logs.Error("HttpPostJson error: ", err)
		return nil, err
	}

	if resp == nil {
		return nil, errors.New("返回对象为空")
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		return body, err
		//logs.Debug("HttpPostJson result: ", result)
	} else {
		logs.Error("HttpPostJson error: ", err)
		return nil, err
	}

}

func HttpPost(url string, param map[string]string) (string, error) {

	var paramBuf bytes.Buffer
	paramBuf.WriteString("curTime=" + GetCurrentTime())
	for k, v := range param {
		paramBuf.WriteString("&" + k + "=" + v)
	}

	resp, err := http.Post(url, "application/x-www-form-urlencoded", strings.NewReader(paramBuf.String()))
	if err != nil {
		logs.Error("HttpPost error: ", err)
		return "", err
	}

	if resp == nil {
		return "", errors.New("返回对象为空")
	}

	defer resp.Body.Close()
	result := ""
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		result = string(body)
		logs.Debug("HttpPost result: ", result)
	} else {
		logs.Error("HttpPost error: ", err)
	}

	return result, nil
}

func UrlEncode(input string) string {
	if IsEmpty(input) {
		return ""
	}
	return url.QueryEscape(input)
}

func UrlDecode(input string) string {
	if IsEmpty(input) {
		return ""
	}
	result, err := url.QueryUnescape(input)
	if err != nil {
		return input
	} else {
		return result
	}
}
