package api

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

type Applications struct {
	Items []Application `json:"items"`
}

type Application struct {
	Metadata AppMeta   `json:"metadata"`
	Status   AppStatus `json:"status"`
}

type AppMeta struct {
	Name      string            `json:"name"`
	Namespace string            `json:"namespace"`
	Labels    map[string]string `json:"labels"`
}

type AppStatus struct {
	Sync   AppStatusSync `json:"sync"`
	Health struct {
		Status string `json:"status"`
	} `json:"health"`
}

type AppStatusSync struct {
	Status   string `json:"status"`
	Revision string `json:"revision"`
}

type IArgoApi interface {
	GetApps(selectors map[string]string, refresh bool) (*Applications, error)
	Sync(name string) error
}
