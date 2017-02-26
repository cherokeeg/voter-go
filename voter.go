package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"bytes"
	"net"
	"time"
	"compress/gzip"
	"strconv"
	"strings"
	"math/rand"
	"os"
)

func seedAndReturnRandom(n int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(n)
}

func IpV4Address() string {
	blocks := []string{}
	for i := 0; i < 4; i++ {
		number := seedAndReturnRandom(255)
		blocks = append(blocks, strconv.Itoa(number))
	}

	return strings.Join(blocks, ".")
}

func main() {

	arg_num := len(os.Args)
	if arg_num != 2 {
		fmt.Println("Useage: tp loop_number")
		return
	}

	loop_number, _ := strconv.Atoi(os.Args[1])

	s := 0
	f := 0
	for i := 0; i < loop_number; i++ {
		fmt.Println(i, " / ", s, " / ", f)
		time.Sleep(2 * time.Second)
		localAddr := IpV4Address()
		localAddress, _ := net.ResolveTCPAddr("tcp", localAddr)

		transport := &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
				LocalAddr: localAddress, }).Dial, TLSHandshakeTimeout: 10 * time.Second, }

		client := &http.Client{
			Transport: transport,
		}

		body := bytes.NewBuffer([]byte("pagecount=1&page=1&key=1&eVotesForm=1&eVotesNum=5"))
		req, _ := http.NewRequest("POST", "http://xxx", body)

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //这个一定要加，不加form的值post不过去，被坑了两小时
		req.Header.Set("Cookie", "xxx")
		req.Header.Set("Host", "xxx")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Content-Length", "50")
		req.Header.Set("Origin", "xxx")
		req.Header.Set("X-Requested-With", "XMLHttpRequest")
		req.Header.Set("User-Agent", "xxx")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Referer", "xxx")
		req.Header.Set("Accept-Encoding", "gzip, deflate")
		req.Header.Set("Accept-Language", "zh-CN,en-US;q=0.8")
		req.Header.Set("Http_Client_Ip", localAddr)
		req.Header.Set("Client_Ip", localAddr)
		req.Header.Set("http_X-Forwarded-For", localAddr)
		req.Header.Set("X-Forwarded-For", localAddr)
		req.Header.Set("Remote_Addr", localAddr)
		req.RemoteAddr = localAddr
		fmt.Printf("%+v\n", req) //看下发送的结构

		resp, err := client.Do(req) //发送

		var resBody []byte
		defer resp.Body.Close()
		if resp.Header.Get("Content-Encoding") == "gzip" {
			reader, err := gzip.NewReader(resp.Body)
			if err != nil {
				fmt.Println(err.Error())
			}
			resBody, err = ioutil.ReadAll(reader)
		} else {
			resBody, err = ioutil.ReadAll(resp.Body)
		}
		fmt.Println(string(resBody), err)
		if strings.Contains(string(resBody), "Error") {
			f++
		} else {
			s++
		}
	}
}
