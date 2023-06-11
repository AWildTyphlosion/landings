package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
)

var (
	unixFlag   = flag.String("unix", "", "Unix socket for http connection")
	hostFlag   = flag.String("host", "", "Host network bindings")
	justOk     = flag.Bool("justok", false, "Only return 200")
	removeSock = flag.Bool("removesock", false, "Remove Sock if it already exists")
)

func run() (err error) {
	flag.Parse()
	if isSet(unixFlag) == isSet(hostFlag) {
		return fmt.Errorf("Must include one of host or unix flag")
	}

	var (
		s     = &http.Server{}
		l     net.Listener
		start func() error
	)

	if isSet(unixFlag) {
		if *removeSock {
			os.Remove(*unixFlag)
		}
		l, err = net.Listen("unix", *unixFlag)
		if err != nil {
			return err
		}
		start = func() error { return s.Serve(l) }
	} else {
		s.Addr = *hostFlag
		start = s.ListenAndServe
	}

	s.Handler = newBasicRouter()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, os.Kill)
		<-sig
		if err := s.Close(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}()

	return start()
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func isSet(s *string) bool { return s != nil && *s != "" }

func deref[T any](t *T) T {
	if t != nil {
		return *t
	}

	return *new(T)
}

func newBasicRouter() *http.ServeMux {
	r := http.NewServeMux()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if *justOk {
			w.WriteHeader(http.StatusOK)
			return
		}
		byts, err := json.MarshalIndent(map[string]string{"hello": "world"}, "", "  ")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("{}"))
			return
		}

		w.Write(byts)
	})

	return r
}
