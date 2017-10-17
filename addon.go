package addon

func New(file string, manager Manager) (*Server, error) {
	manifest, err := readManifest(file)
	if err != nil {
		return nil, err
	}

	server := &Server{
		manifest: manifest,
		manager:  manager,
	}

	server.configure()

	return server, nil
}
