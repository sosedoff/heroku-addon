package addon

type Manager interface {
	Provision(req *ProvisionRequest) (*Resource, error)
	Modify(req *ModifyRequest) (*Resource, error)
	Delete(req *DeleteRequest) (*Resource, error)
}
