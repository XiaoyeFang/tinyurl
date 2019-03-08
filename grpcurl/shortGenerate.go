package grpcurl

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang/glog"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"url-shortener/config"
	"url-shortener/protos"
)

func ShortGenerate(pcUrl, mobileUrl, mobileClientUrl string) (reply *protos.GenerateReply, err error) {
	reply = &protos.GenerateReply{}
	client := &http.Client{}
	var respJson struct {
		Success  bool
		Response string
	}
	var clientType string
	switch {

	case pcUrl != "":

		clientType = "PC"

		//包装url并请求
		request, err := MakeCheckUrl(pcUrl, clientType)
		if err != nil {
			glog.V(0).Infoln(err)
		}
		if request.Header.Get("User-Agent") == "" {
			return reply, errors.New("URL is not valid link")

		}
		resp, err := client.Do(request)

		if err != nil {
			fmt.Println(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &respJson)
		if err != nil {
			return reply, err
		}
		reply.Shortener = respJson.Response
		reply.ClientType = clientType

		return reply, err

	case mobileUrl != "" && pcUrl == "":
		clientType = "MOBILE"

		//包装url并请求
		request, err := MakeCheckUrl(mobileUrl, clientType)
		if err != nil {
			glog.V(0).Infoln(err)
		}
		if request.Header.Get("User-Agent") == "" {
			return reply, errors.New("URL is not valid link")

		}
		resp, err := client.Do(request)

		if err != nil {
			fmt.Println(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &respJson)
		if err != nil {
			return reply, err
		}
		reply.Shortener = respJson.Response
		reply.ClientType = clientType

		return reply, err

	case mobileClientUrl != "" && pcUrl == "" && mobileUrl == "":
		clientType = "MOBILECLIENT"
		//包装url并请求
		request, err := MakeCheckUrl(mobileClientUrl, clientType)
		if err != nil {
			glog.V(0).Infoln(err)
		}
		if request.Header.Get("User-Agent") == "" {
			return reply, errors.New("URL is not valid link")

		}
		resp, err := client.Do(request)

		if err != nil {
			fmt.Println(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		err = json.Unmarshal(body, &respJson)
		if err != nil {
			return reply, err
		}
		reply.Shortener = respJson.Response
		reply.ClientType = clientType
		return reply, err

	}

	return reply, err
}

func MakeCheckUrl(link, clientType string) (*http.Request, error) {
	if !IsValidUrl(link) {
		return &http.Request{}, nil
	}

	r, err := url.Parse(config.UrlConfig.Options.Prefix + "encode/")
	if err != nil {

		return nil,err
	}

	url := `{"Url":"` + link + `"
	       }`

	req, err := http.NewRequest(http.MethodPost, r.String(), strings.NewReader(url))

	//defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", config.HTTP_ACCEPT)
	req.Header.Set("User-Agent", clientType)

	return req, err

}

func IsValidUrl(link string) bool {
	var match bool
	if !strings.HasPrefix(link, "http://") && !strings.HasPrefix(link, "https://") && !strings.HasPrefix(link, "ftp://") {
		return false
	}

	strRegex := "(http://|ftp://|https://|www){0,1}[^\u4e00-\u9fa5\\s]*?\\.(com|net|cn|me|tw|fr)[^\u4e00-\u9fa5\\s]*"
	match, _ = regexp.MatchString(strRegex, link)
	return match
}
