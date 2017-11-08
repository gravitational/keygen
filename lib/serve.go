package lib

import (
	"net/http"

	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
)

func Serve(certPath, keyPath string, hostport string) error {
	if certPath == "" {
		return trace.BadParameter("missing parameter certPath")
	}
	if keyPath == "" {
		return trace.BadParameter("missing parameter keyPath")
	}
	if hostport == "" {
		return trace.BadParameter("missing parameter hostport")
	}

	log.Infof("[DISK] Start serving on %v, certsPath: %v, keyPath: %v", hostport, certPath, keyPath)

	handler, err := NewAPIServer()
	if err != nil {
		return trace.Wrap(err)
	}

	server := &http.Server{
		Addr:    hostport,
		Handler: handler,
	}

	return server.ListenAndServeTLS(certPath, keyPath)
}
