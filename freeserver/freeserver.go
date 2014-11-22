package freeserver

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cenkalti/backoff"
	"github.com/facebookgo/freeport"
	"github.com/koding/logging"
	"github.com/rcrowley/go-tigertonic"
)

// Initer is for initializing the api endpoints for a module
type Initer interface {
	Init(mux *tigertonic.TrieServeMux) *tigertonic.TrieServeMux
}

type LoggerHolder interface {
	Logger() logging.Logger
}

type Server struct {
	Mux    *tigertonic.TrieServeMux
	Port   int
	Server *tigertonic.Server
	// App handles all the required information about the application
	// eg. database connections, singletons etc
	App interface{}

	// DebugEnabled holds the status of the server's debug level
	DebugEnabled bool
}

func New(app interface{}) *Server {
	s := &Server{
		Mux: tigertonic.NewTrieServeMux(),
		App: app,
	}

	return s
}

func (s *Server) InitWith(initers ...Initer) *Server {
	for _, initer := range initers {
		s.Mux = initer.Init(s.Mux)
	}

	return s
}

func (s *Server) Listen() error {
	server, err := s.crateServer()
	if err != nil {
		return err
	}

	s.Server = server

	fmt.Println("started listening on: ", s.Server.Server.Addr)
	// todo move those to a more generic place
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGTERM)

	log.Println(<-ch)
	err = s.Close()

	fmt.Println("server closed succesfully:", err == nil)

	return err
}

func (s *Server) Close() error {
	if s.Server != nil {
		return s.Server.Close()
	}

	return nil
}

func (s *Server) crateServer() (*tigertonic.Server, error) {
	ticker := backoff.NewTicker(
		backoff.NewExponentialBackOff(),
	)

	var err error
	var server *tigertonic.Server
	for _ = range ticker.C {
		if server, err = s.startListening(); err != nil {
			log.Printf("err while creating the server: %s, will try again", err.Error())
			continue
		}

		break
	}

	if err != nil {
		return nil, err
	}

	return server, nil
}

func (s *Server) startListening() (*tigertonic.Server, error) {
	var err error
	if s.Port == 0 {
		s.Port, err = freeport.Get()
		if err != nil {
			return nil, err
		}
	}

	var handler http.Handler

	handler = tigertonic.WithContext(s.Mux, s.App)

	if s.DebugEnabled {
		// TODO make second param configurable
		h := tigertonic.Logged(handler, nil)

		var l logging.Logger

		// if App doesnt have any logger, create a new one
		if logHolder, ok := s.App.(LoggerHolder); ok {
			l = logHolder.Logger()
		} else {
			l = logging.NewLogger("gene")
		}

		h.Logger = NewTigerTonicLogger(l)
		handler = h
	}

	server := tigertonic.NewServer(
		fmt.Sprintf("0.0.0.0:%d", s.Port),
		handler,
	)

	go func() {
		err := server.ListenAndServe()
		if nil != err {
			log.Fatalf("listen err:", err.Error())
		}
	}()

	return server, nil
}
