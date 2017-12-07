package addon

type ProvisionRequest struct {
	UUID         string            `json:"uuid"`
	AppId        string            `json:"heroku_id"`
	Plan         string            `json:"plan"`
	LogplexToken string            `json:"logplex_token"`
	Region       string            `json:"region"`
	Options      map[string]string `json:"options"`
	CallbackURL  string            `json:"callback_url"`
}

type ProvisionResponse struct {
	Id      string            `json:"id"`
	Config  map[string]string `json:"config"`
	Message string            `json:"message,omitempty"`
}

type ModifyRequest struct {
	UUID  string `json:"uuid"`
	AppId string `json:"heroku_id"`
	Plan  string `json:"plan"`
}

type ModifyResponse struct {
	Config  map[string]string `json:"config"`
	Message string            `json:"message,omitempty"`
}

type DeleteRequest struct {
	UUID string `json:"uuid"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

type SSORequest struct {
	Id         string `json:"app"`
	ResourceId string `json:"resource_id"`
	UserId     string `json:"user_id"`
	App        string `json:"app"`
	ContextApp string `json:"context_app"`
	Email      string `json:"email"`
	NavData    string `json:"nav_data"`
}
