package servers

import (
    "fmt"
    "errors"
)

type Packet struct {
    client interface{}
    Message string
}

type Config struct {
    Address string
}

type Server interface {
    Run (cfg Config) (inbound <-chan Packet, outbound chan<- Packet, err error)
    Shutdown()
}

var ErrNoSuchServer = errors.New("No Server by that name has been registered")

func GetServer(name string) (Server, error) {
    b, ok := servers[name]

    if !ok {
        return nil, ErrNoSuchServer
    }

    return b(), nil
}

// Private

type serverBuild func () Server

var servers = make(map[string]serverBuild)

func registerServer (name string, builder serverBuild) {
    if _, ok := servers[name]; ok {
        panic(fmt.Sprintf("Server %s is already registered", name))
    }
    servers[name] = builder
}

