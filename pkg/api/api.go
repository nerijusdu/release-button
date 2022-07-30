package api

type ArgoApi struct {
	token string
}

func NewArgoApi() *ArgoApi {
	return &ArgoApi{}
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
	return a.postJson("/applications/" + name + "/sync")
}
