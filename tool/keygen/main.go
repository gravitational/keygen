package main

import (
	"fmt"
	"os"

	"github.com/gravitational/keygen/lib"

	"github.com/gravitational/kingpin"
	"github.com/gravitational/teleport/lib/auth/native"
	"github.com/gravitational/teleport/lib/utils"
	"github.com/gravitational/trace"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(trace.DebugReport(err))
	}
}

func run() error {
	app := kingpin.New("keygen", "Generate SSH keys easilly")

	cnew := app.Command("new", "Generate new SSH keypair")
	cnewComment := cnew.Flag("comment", "comment to add to the key").String()
	cnewPassphrase := cnew.Flag("pass", "optional passphrase to protect the private key").String()

	cserve := app.Command("serve", "Start keygen server")
	cserveHostPort := cserve.Flag("hostport", "hostport to serve on").Default("127.0.0.1:8080").OverrideDefaultFromEnvar("KEYGEN_HOSTPORT").String()
	cserveCert := cserve.Flag("certPath", "path to cert").OverrideDefaultFromEnvar("KEYGEN_CERT").String()
	cserveKey := cserve.Flag("keyPath", "path to key").OverrideDefaultFromEnvar("KEYGEN_KEY").String()

	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		return trace.Wrap(err)
	}

	utils.InitLogger(utils.LoggingForDaemon, log.DebugLevel)

	switch cmd {
	case cnew.FullCommand():
		return newKey(*cnewComment, *cnewPassphrase)
	case cserve.FullCommand():
		return lib.Serve(*cserveCert, *cserveKey, *cserveHostPort)
	default:
		return trace.BadParameter("unsupported command: %v", cmd)
	}
}

func newKey(comment string, passphrase string) error {
	_, pub, err := native.New().GenerateKeyPair(passphrase)
	if err != nil {
		return trace.Wrap(err)
	}

	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(pub)
	if err != nil {
		return trace.Wrap(err)
	}

	pub = lib.MarshalAuthorizedKey(pubKey, comment)
	if err != nil {
		return trace.Wrap(err)
	}

	fmt.Printf("%v", string(pub))
	return nil
}
