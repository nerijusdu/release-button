package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	netUrl "net/url"
)

type EmptyObj struct{}

func (a *ArgoApi) postJson(url string, reqRes ...interface{}) error {
	var reqData interface{} = EmptyObj{}
	if len(reqRes) > 0 && reqRes[0] != nil {
		reqData = reqRes[0]
	}
	data, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, a.baseUrl+url, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")

	a.addAuthCookies(req)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	fmt.Printf("POST %s - %s\n", url, r.Status)
	if r.StatusCode > 299 {
		return fmt.Errorf("Request failed with status %s", r.Status)
	}

	if len(reqRes) > 1 {
		return parseJson(r, reqRes[1])
	}

	return nil
}

func (a *ArgoApi) getJson(url string, reqRes ...interface{}) error {
	req, err := http.NewRequest(http.MethodGet, a.baseUrl+url, nil)
	if err != nil {
		return err
	}

	a.addAuthCookies(req)

	if len(reqRes) > 0 {
		req.URL.RawQuery = reqRes[0].(netUrl.Values).Encode()
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	if len(reqRes) > 1 {
		return parseJson(r, reqRes[1])
	}

	return nil
}

func (a *ArgoApi) addAuthCookies(req *http.Request) {
	if a.gToken != "" {
		req.Header.Add("Application", a.gToken)
		req.Header.Add("X-User", a.gUser)
	}
	if a.token != "" {
		req.AddCookie(&http.Cookie{
			Name:  "argocd.token",
			Value: a.token,
		})
	}
}

func parseJson(r *http.Response, obj interface{}) error {
	defer r.Body.Close()
	bodyBytes, err2 := io.ReadAll(r.Body)
	if err2 != nil {
		return err2
	}

	if r.StatusCode >= 400 {
		bodyString := string(bodyBytes)
		fmt.Println(bodyString)
	}

	err := json.Unmarshal(bodyBytes, obj)
	if err == nil {
		return nil
	}

	return err
}
