package addon

import (
	"errors"
)

type TestManager struct {
}

func (m *TestManager) Provision(req *ProvisionRequest) (*Resource, error) {
	if req.UUID == "bad" {
		return nil, errors.New("Something is wrong")
	}
	return &Resource{}, nil
}

func (m *TestManager) Modify(req *ModifyRequest) (*Resource, error) {
	if req.UUID == "bad" {
		return nil, errors.New("Something is wrong")
	}
	return &Resource{}, nil
}

func (m *TestManager) Delete(req *DeleteRequest) (*Resource, error) {
	if req.UUID == "bad" {
		return nil, errors.New("Something is wrong")
	}
	return &Resource{}, nil
}
