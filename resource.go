package addon

type Resource struct {
	Id     string            `json:"id"`
	Config map[string]string `json:"config"`
}
