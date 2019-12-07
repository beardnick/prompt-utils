package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"
	"regexp"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

var client *sftp.Client

var CmdSuggests = []prompt.Suggest{
	{"ls", "list files and directories"},
	{"put", "upload files"},
	{"cd", "change working directory"},
	{"exit", "close the prompt"},
}

func GetCwd() (string, error) {
	if len(Cwd) == 0 {
		pwd, err := client.Getwd()
		if err != nil {
			fmt.Println("err:", err)
			return "", err
		}
		Cwd = pwd
	}
	return Cwd, nil
}

var FileSuggests []prompt.Suggest
var FileSet = make(map[string]bool)
var Cwd = ""

func GetFile(file string) error {
	remote, err := client.Open(file)
	defer remote.Close()
	if err != nil {
		return err
	}
	info, err := os.Stat(path.Base(file))
	if err == nil {
		return fmt.Errorf("file %s is already exists", info.Name())
	}
	local, err := os.OpenFile(path.Base(file), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	defer local.Close()
	if err != nil {
		return err
	}
	buf := make([]byte, 1024)
	if _, err := remote.Read(buf); err == nil {
		//fmt.Println("length:", length)
		//fmt.Println("file:", buf)
		_, err = local.Write(buf)
		if err != nil {
			return err
		}
	} else {
		return err
	}
	return nil
}

func connect(url string) error {
	var (
		sshClient *ssh.Client
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
		return err
	}
	// create sftp client
	if client, err = sftp.NewClient(sshClient); err != nil {
		return err
	}
	return nil
}

func uploadFile(localFilePath string, remotePath string) {
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		fmt.Println("os.Open error : ", localFilePath)
		log.Fatal(err)
	}
	defer srcFile.Close()

	var remoteFileName = path.Base(localFilePath)

	dstFile, err := client.Create(path.Join(remotePath, remoteFileName))
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

func refresh() {
	for {
		time.Sleep(time.Second)
		refreshFiles()
	}
}

func uploadDirectory(localPath string, remotePath string) {
	localFiles, err := ioutil.ReadDir(localPath)
	if err != nil {
		log.Fatal("read dir list fail ", err)
	}

	for _, backupDir := range localFiles {
		localFilePath := path.Join(localPath, backupDir.Name())
		remoteFilePath := path.Join(remotePath, backupDir.Name())
		if backupDir.IsDir() {
			client.Mkdir(remoteFilePath)
			uploadDirectory(localFilePath, remoteFilePath)
		} else {
			uploadFile(path.Join(localPath, backupDir.Name()), remotePath)
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
	//fmt.Println("keys:", signer.PublicKey())
	return ssh.PublicKeys(signer), nil
}

func executor(in string) {
	in = strings.TrimSpace(in)
	args := strings.Split(in, " ")
	if len(args) < 1 {
		return
	}
	if len(args) == 1 && args[0] == "exit" {
		os.Exit(0)
	}
	if len(args) == 1 && args[0] == "pwd" {
		fmt.Println(Cwd)
		return
	}
	if len(args) == 1 && args[0] == "cd" {
		pwd, err := client.Getwd()
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		Cwd = pwd
		FileSet = make(map[string]bool)
		FileSuggests = []prompt.Suggest{}
		return
	}
	if len(args) == 1 {
		return
	}
	switch args[0] {
	case "get":
		err := GetFile(client.Join(Cwd, args[1]))
		if err != nil {
			fmt.Println("err:", err)
			return
		}
	case "put":
	case "cd":
		target := args[1]
		current, err := GetCwd()
		if err != nil {
			fmt.Println("err:", err)
			return
		}
		fmt.Println(current, " -> ", target)
		Cwd = client.Join(current, target)
		FileSet = make(map[string]bool)
		FileSuggests = []prompt.Suggest{}
	default:
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	line := in.CurrentLineBeforeCursor()
	args := strings.Split(line, " ")
	if len(args) <= 1 && len(line) > 0 && line[len(line)-1] != ' ' {
		return prompt.FilterHasPrefix(CmdSuggests, in.GetWordBeforeCursor(), true)
	}
	return prompt.FilterHasPrefix(FileSuggests, in.GetWordBeforeCursor(), true)
}

func refreshFiles() {
	pwd, err := GetCwd()
	if err != nil {
		return
	}
	walker := client.Walk(pwd)
	walker.Step()
	for walker.Step() {
		dir := walker.Stat()
		if FileSet[dir.Name()] {
			continue
		} else {
			FileSet[dir.Name()] = true
		}
		filetype := "file"
		if dir.IsDir() {
			walker.SkipDir()
			filetype = "directory"
		}
		FileSuggests = append(FileSuggests, prompt.Suggest{Text: dir.Name(), Description: filetype})
		//FileSuggests = append(FileSuggests, prompt.Suggest{Text: walker.Path()})
	}
}

func main() {
	args := os.Args
	if len(args) < 2 {
		fmt.Println("connect url is needed !")
		return
	}
	pattern := regexp.MustCompile("[a-z]@[0-9.]")
	if !pattern.MatchString(args[1]) {
		fmt.Println("url should be user@host")
		return
	}
	err := connect(args[1])
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	go refresh()
	p := prompt.New(
		executor,
		completer,
		//prompt.OptionLivePrefix(livePrefix),
		prompt.OptionPrefix(args[1]+"> "),
		prompt.OptionTitle("sftp-prompt"),
	)
	p.Run()
}
