package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
)

type EmptyObj struct{}

func (a *ArgoApi) postJson(url string, reqRes ...interface{}) error {
	baseUrl := os.Getenv("ARGOCD_SERVER") + "/api/v1"
	fullUrl := baseUrl + url

	var reqData interface{} = EmptyObj{}
	if len(reqRes) > 0 && reqRes[0] != nil {
		reqData = reqRes[0]
	}
	data, err := json.Marshal(reqData)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodPost, fullUrl, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	if &a.token != nil {
		req.AddCookie(&http.Cookie{
			Name:  "argocd.token",
			Value: a.token,
		})
	}

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if len(reqRes) > 1 {
		return json.NewDecoder(r.Body).Decode(reqRes[1])
	}

	return nil
}
