package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/volcengine/vefaas-golang-runtime/events"
	"github.com/volcengine/vefaas-golang-runtime/vefaas"
	"github.com/volcengine/vefaas-golang-runtime/vefaascontext"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

func main() {
	// Start your vefaas function =D.
	vefaas.Start(handler)
}

// Define your handler function.
func handler(ctx context.Context, r *events.HTTPRequest) (*events.EventResponse, error) {
	fmt.Printf("received new request: %s %s, request id: %s\n", r.HTTPMethod, r.Path, vefaascontext.RequestIdFromContext(ctx))
	fmt.Printf("debug request: header:%v, body:%s\n", r.Headers, r.Body)
	ret := make(map[string]interface{})
	ret["content"] = "Hello From QA YTR Test For gateway!"
	ret["http_method"] = r.HTTPMethod
	ret["http_path"] = r.Path
	query := make(map[string]interface{})
	header := make(map[string]interface{})
	for k, v := range r.QueryStringParameters {
		query[k] = v
	}
	for k, v := range r.Headers {
		header[k] = v
	}
	ret["http_query"] = query
	ret["http_header"] = header
	ret["http_body"] = string(r.Body)
	openapi_url, openapi_body := AccessOpenApi()
	//openapi_url, openapi_body := AccessOpenApi2()
	ret["openapi_body"] = openapi_body
	ret["openapi_url"] = openapi_url
	openapi_url4, openapi_body4 := AccessOpenApi4()
	ret["openapi_body3"] = openapi_body4
	ret["openapi_url3"] = openapi_url4
	retBody, _ := json.Marshal(ret)
	return &events.EventResponse{
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: retBody,
	}, nil
}

func AccessOpenApi() (string, string) {

	//url := "https://developer.toutiao.com/api/apps/qrcode"
	//host := "dev.douyincloud.gateway.egress.ivolces.com"
	//host := "douyincloud.gateway.egress.ivolces.com"
	host := "developer.toutiao.com"

	url := fmt.Sprintf("http://%s/api/v2/tags/text/antidirt", host)
	method := "POST"

	//payload := strings.NewReader(`{"access_token": "0801121846765a5a4d2f6b385a68307237534d43397a667865513d3d","appname": "douyin"}`)
	payloadWithoutToken := strings.NewReader(`{"tasks": [{"content": "要检测的文本"}]}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadWithoutToken)

	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	//req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	resheader := res.Header
	log.Printf("ytr test raw request:%+v\n\n", req)
	log.Printf("ytr test resp from openapi:%+v,%+v\n\n", string(body), resheader)
	return url, string(body)
}

func AccessOpenApi2() (string, string) {

	//url := "https://developer.toutiao.com/api/apps/qrcode"
	//host := "dev.douyincloud.gateway.egress.ivolces.com"
	//host := "douyincloud.gateway.egress.ivolces.com"
	host := "developer.toutiao.com"

	url := fmt.Sprintf("http://%s/api/apps/qrcode", host)
	method := "POST"

	//payload := strings.NewReader(`{"access_token": "0801121846765a5a4d2f6b385a68307237534d43397a667865513d3d","appname": "douyin"}`)
	payloadWithoutToken := strings.NewReader(`{"appname": "douyin"}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadWithoutToken)

	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	//req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	resheader := res.Header
	log.Printf("ytr test raw request:%+v\n\n", req)
	log.Printf("ytr test resp from openapi:%+v,%+v\n\n", string(body), resheader)
	return url, string(body)
}

func AccessOpenApi3() (string, string) {

	//url := "https://developer.toutiao.com/api/apps/qrcode"
	//host := "dev.douyincloud.gateway.egress.ivolces.com"
	//host := "douyincloud.gateway.egress.ivolces.com"
	host := "developer.toutiao.com"
	//url=http://developer-boe.toutiao.com/api/apps/customer_service/url?appid=ttded7a86a41b127b6&access_token=080112184676786c51723461725a346f4e536d342f5242544b673d3d&openid=76f096e2-47c8-4473-8621-71ee839f6378&type=1128&scene=1

	url := fmt.Sprintf("http://%s/api/apps/customer_service/url?appid=ttd01ba6a64f25f03901&openid=76f096e2-47c8-4473-8621-71ee839f6378&type=1128&scene=1", host)
	method := "GET"

	//payload := strings.NewReader(`{"access_token": "0801121846765a5a4d2f6b385a68307237534d43397a667865513d3d","appname": "douyin"}`)
	//payloadWithoutToken := strings.NewReader(`{"appname": "douyin"}`)
	//payloadWithoutToken := strings.NewReader(`{"appname": "douyin"}`)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	//req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	resheader := res.Header
	log.Printf("ytr test raw request:%+v\n\n", req)
	log.Printf("ytr test resp from openapi:%+v,%+v\n\n", string(body), resheader)
	return url, string(body)
}

func AccessOpenApi4() (string, string) {

	//url := "https://developer.toutiao.com/api/apps/qrcode"
	//host := "dev.douyincloud.gateway.egress.ivolces.com"
	//host := "douyincloud.gateway.egress.ivolces.com"
	host := "developer.toutiao.com"

	url := fmt.Sprintf("http://%s/api/apps/convert_video_id/video_id_to_open_item_id", host)
	method := "POST"

	//payload := strings.NewReader(`{"access_token": "0801121846765a5a4d2f6b385a68307237534d43397a667865513d3d","appname": "douyin"}`)
	payloadWithoutToken := strings.NewReader(`{"video_ids":["111"],"app_id":"ttaa3adc873504973d01","access_key":"key"}`)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payloadWithoutToken)

	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	//req.Header.Add("Content-Type", "application/json")
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", ""
	}
	resheader := res.Header
	log.Printf("ytr test raw request:%+v\n\n", req)
	log.Printf("ytr test resp from openapi:%+v,%+v\n\n", string(body), resheader)
	return url, string(body)
}
