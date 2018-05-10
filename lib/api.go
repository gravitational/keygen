package lib

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/gravitational/teleport/lib/httplib"

	"github.com/gravitational/teleport/lib/limiter"
	"github.com/gravitational/teleport/lib/sshutils"
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

	srv.POST("/v1/parsecert", srv.wrap(srv.parseCert))
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
	err = t.Execute(w, templateParams{Header: template.HTML(headerHTML)})
	if err != nil {
		trace.WriteError(w, err)
	}
}

func (s *APIServer) wrap(handler HandlerFunc) httprouter.Handle {
	return httplib.MakeHandler(func(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
		return handler(w, r, p)
	})
}

type parseCertReq struct {
	Cert string `json:"cert"`
}

func (p parseCertReq) Check() error {
	if p.Cert == "" {
		return trace.BadParameter("please paste a valid SSH certificate")
	}
	return nil
}

type parseCertResponse struct {
	Info string `json:"info"`
}

// HumanDateFormatSeconds is a human readable date formatting with seconds
const HumanDateFormatSeconds = "Jan _2 15:04:05 UTC"

func (a *APIServer) parseCert(w http.ResponseWriter, r *http.Request, p httprouter.Params) (interface{}, error) {
	var req parseCertReq
	if err := httplib.ReadJSON(r, &req); err != nil {
		return nil, trace.Wrap(err)
	}
	if err := req.Check(); err != nil {
		return nil, trace.Wrap(err)
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(req.Cert))
	if err != nil {
		return nil, trace.BadParameter("please paste a valid SSH certificate")
	}

	cert, ok := pubKey.(*ssh.Certificate)
	if !ok {
		return nil, trace.BadParameter("please paste a valid SSH certificate, not a public SSH key")
	}

	buf := &bytes.Buffer{}
	tab := tabwriter.NewWriter(buf, 0, 0, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintf(tab, "Certificate Type:\t%v\n", cert.Type())
	fmt.Fprintf(tab, "Public Key:\t%v\n", sshutils.Fingerprint(pubKey))
	fmt.Fprintf(tab, "Signing CA:\t%v\n", sshutils.Fingerprint(cert.SignatureKey))
	fmt.Fprintf(tab, "Key ID:\t%v\n", cert.KeyId)
	fmt.Fprintf(tab, "Principals:\t%v\n", strings.Join(cert.ValidPrincipals, ","))

	if cert.ValidAfter == 0 {
		fmt.Fprintf(tab, "Valid After:\teffective immediatelly\n")
	} else {
		fmt.Fprintf(tab, "Valid After:\t%v\n", time.Unix(int64(cert.ValidAfter), 0).Format(HumanDateFormatSeconds))
	}

	if cert.ValidBefore == 0 {
		fmt.Fprintf(tab, "Valid Before:\tdoes not expire\n")
	} else {
		fmt.Fprintf(tab, "Valid Before:\t%v\n", time.Unix(int64(cert.ValidBefore), 0).Format(HumanDateFormatSeconds))
	}

	if len(cert.CriticalOptions) == 0 {
		fmt.Fprintf(tab, "Critical Options: none\n")
	} else {
		fmt.Fprintf(tab, "Critical Options:\n")
		for key, val := range cert.CriticalOptions {
			if val == "" {
				fmt.Fprintf(tab, "    %v\n", key)
			} else {
				fmt.Fprintf(tab, "    %v:\t%v\n", key, val)
			}
		}
	}

	if len(cert.Extensions) == 0 {
		fmt.Fprintf(tab, "Extensions: none\n")
	} else {
		fmt.Fprintf(tab, "Extensions:\n")
		for key, val := range cert.Extensions {
			if val == "" {
				fmt.Fprintf(tab, "    %v\n", key)
			} else {
				fmt.Fprintf(tab, "    %v:\t%v\n", key, val)
			}
		}
	}

	tab.Flush()

	return &parseCertResponse{
		Info: buf.String(),
	}, nil
}
