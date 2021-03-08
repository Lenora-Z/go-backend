package utils

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"
)

//Created by Goland
//@User: lenora
//@Date: 2021/1/15
//@Time: 10:11 上午

func GetRandomStringSec(lenght int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	bytesLen := len(bytes)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < lenght; i++ {
		result = append(result, bytes[r.Intn(bytesLen)])
	}
	return string(result)

}

func GetUUid() (string, error) {
	var err error
	u1 := uuid.Must(uuid.NewV4(), err).String()
	if err != nil {
		logrus.Error("failed to parse uuid: ", err)
		return "", err
	}
	return u1, nil

}

func IsContain(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}

func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}

type HeaderRequest struct {
	Method string
	Url    string
	Header map[string]string
	Body   string
}

type CommonResponse struct {
	Code    uint   `json:"code"`
	Message string `json:"message"`
}

func SendRequest(headerRequest *HeaderRequest) (string, error) {
	var req *http.Request
	var err error
	switch headerRequest.Method {
	case "GET":
		req, err = http.NewRequest("GET", headerRequest.Url, nil)
	case "POST":
		jsonBody := []byte(headerRequest.Body)
		req, err = http.NewRequest("POST", headerRequest.Url, bytes.NewBuffer(jsonBody))
	default:
		req = nil
		err = errors.New("wrong method")
	}

	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	header := headerRequest.Header
	for i := range header {
		req.Header.Add(i, header[i])
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("request => GetWithParam: ", err)
		return "", err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), nil
}

func ReadFile(url string) (string, error) {
	file, err := os.Open(url)
	if err != nil {
		return "", err
	}

	defer file.Close()
	fileinfo, err := file.Stat()
	if err != nil {
		return "", err
	}

	fileSize := fileinfo.Size()
	buffer := make([]byte, fileSize)

	bytesread, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	fmt.Println("bytes read:", bytesread)
	//fmt.Println("bytestream to string:", string(buffer))
	return string(buffer), nil
}
