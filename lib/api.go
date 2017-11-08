package lib

import (
	"html/template"
	"net/http"
	"time"

	"github.com/gravitational/teleport/lib/httplib"

	"github.com/gravitational/teleport/lib/auth/native"
	"github.com/gravitational/teleport/lib/limiter"
	"github.com/gravitational/trace"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/crypto/ssh"
)

// APIServer implements http API server for AuthServer interface
type APIServer struct {
	httprouter.Router
	limiter *limiter.Limiter
}

// NewAPIServer returns a new instance of APIServer HTTP handler
func NewAPIServer() (http.Handler, error) {
	srv := APIServer{}
	srv.Router = *httprouter.New()
	limiter, err := limiter.NewLimiter(limiter.LimiterConfig{
		Rates: []limiter.Rate{
			{
				Period:  time.Second,
				Average: 5,
				Burst:   10,
			},
		},
		MaxConnections:   100,
		MaxNumberOfUsers: 100,
	})
	if err != nil {
		return nil, trace.Wrap(err)
	}

	srv.POST("/v1/keypairs", srv.wrap(srv.newKeyPair))
	srv.GET("/", srv.index)
	limiter.WrapHandle(&srv.Router)
	return limiter, nil
}

// HandlerFunc is http handler
type HandlerFunc func(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error)

func (s *APIServer) index(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	t := template.New("form")
	t, err := t.Parse(indexHTML)
	if err != nil {
		trace.WriteError(w, err)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	err = t.Execute(w, nil)
	if err != nil {
		trace.WriteError(w, err)
	}
}

func (s *APIServer) wrap(handler HandlerFunc) httprouter.Handle {
	return httplib.MakeHandler(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
		return handler(w, r, p)
	})
}

type newKeyPairReq struct {
	Comment    string `json:"commment"`
	Passphrase string `json:"passphrase"`
}

type newKeyPairResponse struct {
	Pub  string `json:"pub"`
	Priv string `json:"priv"`
}

func (a *APIServer) newKeyPair(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
	var req newKeyPairReq
	if err := httplib.ReadJSON(r, &req); err != nil {
		return nil, trace.Wrap(err)
	}

	priv, pub, err := native.New().GenerateKeyPair(req.Passphrase)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pub)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	pub = MarshalAuthorizedKey(pubKey, req.Comment)
	if err != nil {
		return nil, trace.Wrap(err)
	}

	return &newKeyPairResponse{Pub: string(pub), Priv: string(priv)}, nil
}
