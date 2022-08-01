package api

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type ArgoApi struct {
	token   string
	baseUrl string
	fToken  string
}

func NewArgoApi() *ArgoApi {
	return &ArgoApi{
		baseUrl: os.Getenv("ARGOCD_SERVER") + "/api/v1",
		fToken:  os.Getenv("FORWARDED_TOKEN"),
	}
}

func (a *ArgoApi) LoadToken(req AuthRequest) error {
	res := new(AuthResponse)
	err := a.postJson("/session", req, res)
	if err != nil {
		return err
	}

	a.token = res.Token
	return nil
}

func (a *ArgoApi) Sync(name string) error {
	return a.postJson("/applications/" + name + "/syncaaaaa")
}

func (a *ArgoApi) GetApps(selectors map[string]string, refresh bool) (*Applications, error) {
	q := url.Values{}
	s := ""

	q.Add("refresh", strconv.FormatBool(refresh))

	for k, v := range selectors {
		s = s + fmt.Sprintf("%s=%s,", k, v)
	}
	if s != "" {
		q.Add("selector", strings.TrimRight(s, ","))
	}

	apps := new(Applications)
	err := a.getJson("/applications", q, apps)
	return apps, err
}