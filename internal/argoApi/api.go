package argoApi

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type ArgoApi struct {
	token      string
	baseUrl    string
	gToken     string
	gUser      string
	isNewToken bool
}

func NewArgoApi() *ArgoApi {
	return &ArgoApi{
		baseUrl: os.Getenv("ARGOCD_SERVER") + "/api/v1",
		gToken:  os.Getenv("AUTH_TOKEN"),
		gUser:   os.Getenv("AUTH_USER"),
	}
}

func (a *ArgoApi) LoadToken() error {
	if a.isNewToken {
		return fmt.Errorf("Dont make infinite loops pls")
	}

	res := new(AuthResponse)
	req := AuthRequest{
		Username: os.Getenv("ARGOCD_USERNAME"),
		Password: os.Getenv("ARGOCD_PASSWORD"),
	}

	err := a.postJson("/session", req, res)
	if err != nil {
		return err
	}

	a.token = res.Token
	return nil
}

func (a *ArgoApi) Sync(name string) error {
	return a.postJson("/applications/" + name + "/sync")
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
