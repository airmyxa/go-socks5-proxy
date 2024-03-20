package proto

import (
	"fmt"
	"net"
)

type ConnListener struct {
	ln          net.Listener
	sessionPool SessionPool
}

func (cl *ConnListener) Listen() error {
	var err error
	cl.ln, err = net.Listen("tcp", ":1080")
	if err != nil {
		fmt.Printf("Error while starting listen for connections: %s", err.Error())
		return err
	}
	defer cl.ln.Close()

	for {
		conn, err := cl.ln.Accept()
		if err != nil {
			fmt.Printf("Error while connecting to client: %s", err.Error())
			continue
		}
		session := cl.sessionPool.NewSession(conn)
		go func() {
			err := session.Start(conn)
			if err != nil {
				fmt.Printf("Error while starting session: %s", err.Error())
			}
		}()
	}
}
