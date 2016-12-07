// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/terminal"
)

const timeout = 10 * time.Second

const (
	endpointInfo = "http://plugin/gpu/info/json"
	endpointCLI  = "http://plugin/docker/cli"
)

func GenerateRemoteArgs(image string, vols []string) ([]string, error) {
	var info struct {
		Version struct{ CUDA string }
	}

	c := httpClient(Host)

	r, err := c.Get(endpointInfo)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&info); err != nil {
		return nil, err
	}
	if err := cudaSupported(image, info.Version.CUDA); err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("%s?vol=%s&dev=%s", endpointCLI,
		strings.Join(vols, "+"),
		strings.Join(GPU, "+"),
	)
	r2, err := c.Get(uri)
	if err != nil {
		return nil, err
	}
	defer r2.Body.Close()

	b, err := ioutil.ReadAll(r2.Body)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(b), " "), nil
}

func httpClient(addr *url.URL) *http.Client {
	dial := func(string, string) (net.Conn, error) {
		if addr.Scheme == "ssh" {
			c, err := ssh.Dial("tcp", addr.Host, &ssh.ClientConfig{
				User: addr.User.Username(),
				Auth: sshAuths(addr),
			})
			if err != nil {
				return nil, err
			}
			return c.Dial("tcp", addr.Opaque)
		}
		return net.Dial("tcp", addr.Host)
	}

	return &http.Client{
		Timeout:   timeout,
		Transport: &http.Transport{Dial: dial},
	}
}

func sshAuths(addr *url.URL) (methods []ssh.AuthMethod) {
	if sock := os.Getenv("SSH_AUTH_SOCK"); sock != "" {
		c, err := net.Dial("unix", sock)
		if err != nil {
			log.Println("Warning: failed to contact the local SSH agent")
		} else {
			auth := ssh.PublicKeysCallback(agent.NewClient(c).Signers)
			methods = append(methods, auth)
		}
	}
	auth := ssh.PasswordCallback(func() (string, error) {
		fmt.Printf("%s@%s password: ", addr.User.Username(), addr.Host)
		b, err := terminal.ReadPassword(int(syscall.Stdin))
		fmt.Print("\n")
		return string(b), err
	})
	methods = append(methods, auth)
	return
}
