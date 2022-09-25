package web

type SaveConfigRequest struct {
	AllowedApps []string `json:"allowedApps"`
}
