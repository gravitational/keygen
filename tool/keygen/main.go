package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"

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

	cmd, err := app.Parse(os.Args[1:])
	if err != nil {
		return trace.Wrap(err)
	}

	utils.InitLogger(utils.LoggingForDaemon, log.DebugLevel)

	switch cmd {
	case cnew.FullCommand():
		return newKey(*cnewComment, *cnewPassphrase)
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

	pub = marshalAuthorizedKey(pubKey, comment)
	if err != nil {
		return trace.Wrap(err)
	}

	fmt.Printf("%v", string(pub))
	return nil
}

// marshalAuthorizedKey serializes key for inclusion in an OpenSSH
// authorized_keys file. The return value ends with newline.
func marshalAuthorizedKey(key ssh.PublicKey, comment string) []byte {
	b := &bytes.Buffer{}
	b.WriteString(key.Type())
	b.WriteByte(' ')
	e := base64.NewEncoder(base64.StdEncoding, b)
	e.Write(key.Marshal())
	e.Close()
	if comment != "" {
		b.WriteString(" " + comment)
	}
	b.WriteByte('\n')
	return b.Bytes()
}
