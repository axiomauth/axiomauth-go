# axiomauth-go

Official Go SDK for the [AxiomAuth](https://axiomauth.com) identity platform.

[![Go Reference](https://pkg.go.dev/badge/github.com/axiomauth/axiomauth-go.svg)](https://pkg.go.dev/github.com/axiomauth/axiomauth-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

## Installation

```bash
go get github.com/axiomauth/axiomauth-go
```

## Quick Start

```go
package main

import (
    "fmt"
    "os"
    "github.com/axiomauth/axiomauth-go"
)

func main() {
    client := axiomauth.New(os.Getenv("AXIOMAUTH_API_KEY"))

    users, err := client.Users.List(nil)
    if err != nil {
        panic(err)
    }
    fmt.Printf("%d users\n", users.Total)

    config, err := client.Config.Get()
    if err != nil {
        panic(err)
    }
    fmt.Println(config.SAMLMetadataURL)
}
```

## API Reference

### Users
```go
users, _ := client.Users.List(&axiomauth.ListParams{Page: 1, PerPage: 20})
user,  _ := client.Users.Get("usr_abc123")
_,     _ =  client.Users.Deprovision("usr_abc123")
```

### Sessions
```go
sessions, _ := client.Sessions.List(nil)
_,         _ =  client.Sessions.Revoke("sess_abc123")
```

### Config
```go
cfg, _ := client.Config.Get()
_,   _ =  client.Config.Update(&axiomauth.ConfigUpdate{SessionDurationHours: 12})
```

### Audit
```go
events, _ := client.Audit.List(&axiomauth.AuditParams{Action: "login.success"})
```

## Support

- Docs: [axiomauth.com/docs](https://axiomauth.com/docs)
- Email: support@axiomauth.com

## License

[MIT](LICENSE) — © Axiom Identity Services Ltd.
