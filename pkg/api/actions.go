package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
)

func (a *ArgoApi) postJson(url string, reqRes ...interface{}) error {
	baseUrl := os.Getenv("ARGOCD_SERVER") + "/api/v1"
	fullUrl := baseUrl + url

	var dataReader io.Reader = nil
	if len(reqRes) > 0 && reqRes[0] != nil {
		data, err := json.Marshal(reqRes[0])
		if err != nil {
			return err
		}
		dataReader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(http.MethodPost, fullUrl, dataReader)
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
