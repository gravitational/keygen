package lib

import (
	"bytes"
	"encoding/base64"

	"golang.org/x/crypto/ssh"
)

// MarshalAuthorizedKey serializes key for inclusion in an OpenSSH
// authorized_keys file. The return value ends with newline.
func MarshalAuthorizedKey(key ssh.PublicKey, comment string) []byte {
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
