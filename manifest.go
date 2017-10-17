package addon

import (
	"encoding/json"
	"io/ioutil"
)

type Manifest struct {
	Id  string `json:"id"`
	Api struct {
		ConfigVars []string `json:"config_vars"`
		Regions    []string `json:"regions"`
		Password   string   `json:"password"`
		SsoSalt    string   `json:"sso_salt"`
	} `json:"api"`
}

func readManifest(path string) (*Manifest, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	manifest := &Manifest{}

	if err := json.Unmarshal(data, manifest); err != nil {
		return nil, err
	}
	return manifest, nil
}

func (m *Manifest) isValidAuth(auth BasicAuth) bool {
	return auth.Username == m.Id && auth.Password == m.Api.Password
}
