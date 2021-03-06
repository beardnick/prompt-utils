package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// sftp的上下文，对原生sftp库的封装
type Sftp struct {
	localPath  string
	remotePath string
	Client     *sftp.Client
	Cmds       map[string]ICmd
}

func (s *Sftp) LocalPath() string {
	return s.localPath
}

func (s *Sftp) RemotePath() string {
	return s.remotePath
}

func (s *Sftp) LGetCwd() string {
	return s.localPath
}

func (s *Sftp) GetCwd() string {
	return s.remotePath
}

func (s *Sftp) RemoteFiles() []os.FileInfo {
	log.Println("sftp.remotePath:", s.remotePath)
	files, err := s.Client.ReadDir(s.remotePath)
	if err != nil {
		log.Fatal("read dir err:", err)
	}
	return files
}

func (s *Sftp) RemotePathFiles(path string) []os.FileInfo {
	//log.Println("path:", path)
	files, err := s.Client.ReadDir(path)
	if err != nil {
		log.Fatal("read dir err:", err)
	}
	return files
}

//// #TODO: 23-02-20 ls local files //
//func (s Sftp) LLs() []os.FileInfo {
//    files, err := s.Client.ReadDir(s.remotePath)
//    if err != nil {
//        log.Fatal("read dir err:", err)
//    }
//    return files
//}

var sshAddr = regexp.MustCompile(`([a-zA-Z]+)@([1-9.]+)(:[0-9]+)?`)

func (s *Sftp) Connect(url string) (err error) {
	if !sshAddr.MatchString(url) {
		return errors.New("not a valid ssh address")
	}
	groups := sshAddr.FindStringSubmatch(url)
	account := groups[1]
	host := groups[2]
	p := strings.Trim(groups[3], ":")
	port, err := strconv.Atoi(p)
	if err != nil {
		port = 22
	}
	log.Printf("account:%v host:%v", account, host)
	config := &ssh.ClientConfig{
		Timeout:         time.Second * 10,
		User:            account,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	auth, err := publicKeyAuthFunc()
	if err == nil {
		config.Auth = []ssh.AuthMethod{auth}
	} else {
		// @todo: 支持密码登陆 <04-12-19> //
		log.Fatal("now public key")
	}
	// connet to ssh
	addr := fmt.Sprintf("%s:%d", host, port)
	sshClient, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return err
	}
	// create sftp client
	client, err := sftp.NewClient(sshClient)
	if err != nil {
		return err
	}
	s.Client = client
	s.remotePath, err = s.Client.Getwd()
	if err != nil {
		return err
	}
	log.Printf("remotePath: %v", s.remotePath)
	return nil
}

func NewSftp(url string) (*Sftp, error) {
	s := &Sftp{}
	err := s.Connect(url)
	if err != nil {
		return nil, err
	}
	l := []ICmd{
		&Ls{},
		&Exit{},
		&Quit{},
		&Bye{},
		&Get{},
	}
	s.Cmds = map[string]ICmd{}
	for _, v := range l {
		v.Init(s)
		s.Cmds[v.Name()] = v
	}
	return s, nil
}

func publicKeyAuthFunc() (ssh.AuthMethod, error) {
	current, err := user.Current()
	if err != nil {
		return nil, err
	}
	path := current.HomeDir + "/.ssh/id_rsa"
	key, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("ssh key file read failed:", err)
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		log.Fatal("ssh key signer failed:", err)
	}
	//fmt.Println("keys:", signer.PublicKey())
	return ssh.PublicKeys(signer), nil
}
