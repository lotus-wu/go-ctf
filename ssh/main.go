package main

import (
	"bytes"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
)

func main() {
	fmt.Println("helloworld")
	// An SSH client is represented with a ClientConn. Currently only
	// the "password" authentication method is supported.
	//
	// To authenticate with the remote server you must pass at least one
	// implementation of AuthMethod via the Auth field in ClientConfig.
	config := &ssh.ClientConfig{
		User: "git",
		Auth: []ssh.AuthMethod{
			ssh.Password("12345687"),
		},
		//需要验证服务端，不做验证返回nil就可以，点击HostKeyCallback看源码就知道了
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
	client, err := ssh.Dial("tcp", "127.0.0.1:2222", config)
	if err != nil {
		panic("Failed to dial: " + err.Error())
	}

	// Each ClientConn can support multiple interactive sessions,
	// represented by a Session.
	session, err := client.NewSession()
	if err != nil {
		panic("Failed to create session: " + err.Error())
	}
	defer session.Close()

	// Once a Session is created, you can execute a single command on
	// the remote side using the Run method.
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run("/usr/bin/whoami"); err != nil {
		panic("Failed to run: " + err.Error())
	}
	fmt.Println(b.String())

}
