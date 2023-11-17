package argoApi

import (
	"encoding/json"
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
		return fmt.Errorf("dont make infinite loops pls")
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

type SyncResponseDetails struct {
	TypeUrl string `json:"type_url"`
	Value   string `json:"value"`
}

type SyncResponse struct {
	Code     int                 `json:"code"`
	Error    string              `json:"error"`
	Message  string              `json:"message"`
	Details  SyncResponseDetails `json:"details"`
	Metadata AppMeta             `json:"metadata"`
	Status   AppStatus           `json:"status"`
}

func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, "", "\t")
	return string(s)
}

func (a *ArgoApi) Sync(name string) error {
	var res SyncResponse
	err := a.postJson("/applications/"+name+"/sync", nil, &res)
	fmt.Println(prettyPrint(res))

	return err
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
