package server

import (
	"bufio"
	"net"
	"sync"
	"time"

	"jj/resp"

	"strings"

	"github.com/juju/errors"
	log "github.com/ngaut/logging"
)

type cmdFunc func(r *resp.Resp, client *session) (*resp.Resp, error)

var (
	cmdFuncs = map[string]cmdFunc{
		"jdocset": cmdJdocSet,
		"jdocget": cmdJdocGet,
	}
)

type Server struct {
	db   Db
	lock sync.RWMutex
	addr string
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
		db:   NewMapDb(),
		lock: sync.RWMutex{},
	}
}

func (s *Server) Run() {
	log.Info("listening on", s.addr)
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Warning(errors.ErrorStack(err))
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(c net.Conn) {
	log.Info("new connection", c.RemoteAddr())

	client := &session{
		Conn:     c,
		srv:      s,
		r:        bufio.NewReader(c),
		CreateAt: time.Now(),
	}

	var err error

	defer func() {
		if err != nil {
			log.Infof("close connection %v, %+v", c.RemoteAddr(), client)
		}
		c.Close()
	}()

	for {
		r, err := resp.Parse(client.r)
		if err != nil {
			log.Warning(err)
			return
		}

		op, err := r.Op()
		if err != nil {
			log.Warning(err)
		}

		strOp := strings.ToLower(string(op))

		var ret *resp.Resp
		f, ok := cmdFuncs[strOp]
		if !ok {
			ret = RespNoSuchCmd
		} else {
			ret, err = f(r, client)
			if err != nil {
				log.Warning(err)
				return
			}
		}
		b, _ := ret.Bytes()
		client.Write(b)
	}
}
