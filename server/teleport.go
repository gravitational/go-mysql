package server

import (
	"net"
	"sync/atomic"

	. "github.com/go-mysql-org/go-mysql/mysql"
	"github.com/go-mysql-org/go-mysql/packet"
)

func (c *Conn) ReadHandshakeResponse() error {
	return c.readHandshakeResponse()
}

func (c *Conn) WriteInitialHandshake() error {
	return c.writeInitialHandshake()
}

func (c *Conn) WriteOK(r *Result) error {
	return c.writeOK(r)
}

func (c *Conn) WriteError(e error) error {
	return c.writeError(e)
}

// MakeConn creates a new server side connection without performing the handshake.
func MakeConn(conn net.Conn, serverConf *Server, p CredentialProvider, h Handler) *Conn {
	var packetConn *packet.Conn
	if serverConf.tlsConfig != nil {
		packetConn = packet.NewTLSConn(conn)
	} else {
		packetConn = packet.NewConn(conn)
	}

	salt, _ := RandomBuf(20)
	c := &Conn{
		Conn:               packetConn,
		serverConf:         serverConf,
		credentialProvider: p,
		h:                  h,
		connectionID:       atomic.AddUint32(&baseConnID, 1),
		stmts:              make(map[uint32]*Stmt),
		salt:               salt,
	}
	c.closed.Set(false)

	return c
}
