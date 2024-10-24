package sshClose

import "golang.org/x/crypto/ssh"

func CloseSsh(conn *ssh.Client) {
	if conn != nil {
		_ = conn.Close()
	}
}
