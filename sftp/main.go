package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

func connect(url string) (*sftp.Client, error) {
	var (
		sshClient  *ssh.Client
		sftpClient *sftp.Client
	)
	account := url[:strings.Index(url, "@")]
	host := url[strings.Index(url, "@")+1:]

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
	port := 22
	// connet to ssh
	addr := fmt.Sprintf("%s:%d", host, port)
	if sshClient, err = ssh.Dial("tcp", addr, config); err != nil {
		return nil, err
	}
	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}
	return sftpClient, nil
}

func uploadFile(sftpClient *sftp.Client, localFilePath string, remotePath string) {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		log.Fatal(err)

	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)

	dstFile, err := sftpClient.Create(path.Join(remotePath, remoteFileName))
	if err != nil {
		fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
		log.Fatal(err)

	}
	defer dstFile.Close()

	ff, err := ioutil.ReadAll(srcFile)
	if err != nil {
		fmt.Println("ReadAll error : ", localFilePath)
		log.Fatal(err)

	}
	dstFile.Write(ff)
	fmt.Println(localFilePath + "  copy file to remote server finished!")
}

func uploadDirectory(sftpClient *sftp.Client, localPath string, remotePath string) {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Fatal("read dir list fail ", err)
	}

	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			sftpClient.Mkdir(remoteFilePath)
			uploadDirectory(sftpClient, localFilePath, remoteFilePath)
		} else {
			uploadFile(sftpClient, path.Join(localPath, backupDir.Name()), remotePath)
		}
	}

	fmt.Println(localPath + "  copy directory to remote server finished!")
}

//func DoBackup(host string, port int, userName string, password string, localPath string, remotePath string) {
//var (
//err        error
//sftpClient *sftp.Client
//)
//start := time.Now()
//sftpClient, err = connect(userName, password, host, port)
//if err != nil {
//log.Fatal(err)
//}
//defer sftpClient.Close()

//_, errStat := sftpClient.Stat(remotePath)
//if errStat != nil {
//log.Fatal(remotePath + " remote path not exists!")
//}

//backupDirs, err := ioutil.ReadDir(localPath)
//if err != nil {
//log.Fatal(localPath + " local path not exists!")
//}
//uploadDirectory(sftpClient, localPath, remotePath)
//elapsed := time.Since(start)
//fmt.Println("elapsed time : ", elapsed)
//}

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
	fmt.Println("keys:", signer.PublicKey())
	return ssh.PublicKeys(signer), nil
}

func main() {
	client, err := connect("qianz@192.168.3.227")
	if err != nil {
		fmt.Println("err:", err)
	}
	dir, _ := client.Getwd()
	fmt.Println("client:", dir)
}
