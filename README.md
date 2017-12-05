# heroku-addon

Heroku Addon SDK for Go application

## Install

Install the package:

```
go get -u github.com/sosedoff/heroku-addon
```

## Usage

Minimalistic example:

```golang
package main

import (
  "log"

  "github.com/sosedoff/heroku-addon"
)

// Our resource manager
type Manager struct {
}

func (m *Manager) Provision(req *addon.ProvisionRequest) (*addon.Resource, error) {
  log.Println("Provision request:", req)
  return nil, nil
}

func (m *Manager) Modify(req *addon.ModifyRequest) (*addon.Resource, error) {
  log.Println("Modify request:", req)
  return nil, nil
}

func (m *Manager) Delete(req *addon.DeleteRequest) (*addon.Resource, error) {
  log.Println("Delete request:", req)
  return nil, nil
}

func main() {
  manager := &Manager{}

  server, err := addon.New("./addon-manifest.json", manager)
  if err != nil {
    log.Fatal(err)
  }

  if err := server.Start(":4567"); err != nil {
    log.Fatal(err)
  }
}
```

Then run the demo:

```
go run main.go
```

## License

MIT