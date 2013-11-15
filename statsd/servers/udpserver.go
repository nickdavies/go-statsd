package servers

import (
    "fmt"
    "net"
    "sync"
)

const MAXPACKET = 8192

type udpServer struct {
    sync.Mutex
    shutdown bool
    conn *net.UDPConn
}

func (s *udpServer) Run(cfg Config) (inbound <-chan Packet, outbound chan<- Packet, err error) {
    addr, err := net.ResolveUDPAddr("udp", cfg.Address)
    if err != nil {
        return nil, nil, err
    }

    s.conn, err = net.ListenUDP("udp", addr)
    if err != nil {
        return nil, nil, err
    }

    var raw_inbound = make(chan Packet)
    var raw_outbound = make(chan Packet)

    // Send replies
    go func() {
        for p := range raw_outbound {
            // Ok to panic because this should never be messed with
            sendAddr := p.client.(*net.UDPAddr)

            _, err := s.conn.WriteTo([]byte(p.Message), sendAddr)
            if err != nil {
                // TODO: handle error properly
                fmt.Println("SEND ERROR!", err)
            }
        }
    }()

    go func (){
        defer close(raw_inbound)
        for {
            var rawPacket = make([]byte, MAXPACKET)

            n, retaddr, err := s.conn.ReadFromUDP(rawPacket)
            if err != nil {
                if s.shutdown == true {
                    return
                }
                // TODO: handle error properly
                fmt.Println("RECV ERROR!", err)
            }

            raw_inbound <- Packet{
                client: retaddr,
                Message: string(rawPacket[:n]),
            }
        }
    }()

    return raw_inbound, raw_outbound, nil
}

func (s *udpServer) Shutdown() {
    s.shutdown = true
    s.conn.Close()
}

func init () {
    registerServer("udp", func() Server {
        return &udpServer{}
    })
}
