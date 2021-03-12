package tcp

import (
	"bufio"
	"github.com/zucchinidev/hands-on-concurrency-go/tcp-memcache/cache"
	"github.com/zucchinidev/hands-on-concurrency-go/tcp-memcache/printer"
	"io"
	"net"
	"strings"
)

type server struct {
	address  string
	listener net.Listener
}

func New(address string) *server {
	return &server{address: address}
}

func (s *server) Listen() error {
	listener, errDialing := net.Listen("tcp", s.address)
	if errDialing != nil {
		return errDialing
	}
	s.listener = listener
	return nil
}

func (s *server) Close() error {
	return s.listener.Close()
}

func (s *server) Accept() (io.ReadWriteCloser, error) {
	return s.listener.Accept()
}

func (s *server) Invoke(rw io.ReadWriteCloser, ch *cache.Cache) {
	scanner := bufio.NewScanner(rw)
	for scanner.Scan() {
		line := strings.ToLower(strings.Trim(scanner.Text(), " "))
		values := strings.Split(line, " ")
		switch {
		case len(values) == 3 && values[0] == "set":
			ch.Set(values[1], values[2])
			printer.Print(rw, "OK\n-> ")
		case len(values) == 2 && values[0] == "get":
			printer.Print(rw, "%v\n-> ", ch.Get(values[1]))
		case len(values) == 1 && values[0] == "exit":
			_ = rw.Close()
		default:
			printer.Print(rw, "UNKNOWN: %v\n", line)
		}
	}
}
