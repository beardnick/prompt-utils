package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var (
	Cmds map[string]ICmd = map[string]ICmd{}
)

type Sftp struct {
	localPath  string
	remotePath string
	Client     *sftp.Client
}

//func (s Sftp) RemotePath() string {
//    return s.remotePath
//}

func (s *Sftp) LocalPath() string {
	return s.localPath
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

//// #TODO: 23-02-20 ls local files //
//func (s Sftp) LLs() []os.FileInfo {
//    files, err := s.Client.ReadDir(s.remotePath)
//    if err != nil {
//        log.Fatal("read dir err:", err)
//    }
//    return files
//}

func (s *Sftp) Connect(url string) error {
	account := url[:strings.Index(url, "@")]
	host := url[strings.Index(url, "@")+1:]
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
	addr := fmt.Sprintf("%s:%d", host, 22)
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

func init() {
	l := []ICmd{
		Ls{},
		Exit{},
		Quit{},
		Bye{},
		Get{}.Init(),
	}
	for _, v := range l {
		Cmds[v.Name()] = v
		CmdSuggests = append(CmdSuggests, prompt.Suggest{v.Name(), v.Description()})
	}
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
