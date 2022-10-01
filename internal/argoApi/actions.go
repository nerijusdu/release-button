package argoApi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	netUrl "net/url"
	"strings"
)

type EmptyObj struct{}

func (a *ArgoApi) postJson(url string, reqRes ...any) error {
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

	return a.doRequest(req, "POST", url, reqRes)
}

func (a *ArgoApi) getJson(url string, reqRes ...any) error {
	req, err := http.NewRequest(http.MethodGet, a.baseUrl+url, nil)
	if err != nil {
		return err
	}

	if len(reqRes) > 0 {
		req.URL.RawQuery = reqRes[0].(netUrl.Values).Encode()
	}

	return a.doRequest(req, "GET", url, reqRes)
}

func (a *ArgoApi) doRequest(req *http.Request, m, url string, reqRes []any) error {
	a.addAuthCookies(req)

	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	fmt.Printf("%s %s - %s\n", m, url, r.Status)
	if r.StatusCode > 299 {
		if r.StatusCode == 401 && !strings.Contains(url, "/session") {
			fmt.Println("Unauthorized. Re-fetching token.")
			err := a.LoadToken()
			if err != nil {
				return fmt.Errorf("Failed to fetch argo token. %v \n", err)
			}

			a.isNewToken = true
			defer func() {
				a.isNewToken = false
			}()

			if m == "GET" {
				return a.getJson(url, reqRes...)
			}
			return a.postJson(url, reqRes...)
		}

		return fmt.Errorf("Request failed with status %s", r.Status)
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

	bodyString := string(bodyBytes)
	if r.StatusCode >= 400 || strings.HasPrefix(bodyString, "<") {
		fmt.Println(bodyString)
	}

	err := json.Unmarshal(bodyBytes, obj)
	if err == nil {
		return nil
	}

	return err
}
