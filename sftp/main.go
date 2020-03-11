package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/c-bata/go-prompt"
)

var (
	sftpCtx *Sftp
)

//func uploadFile(localFilePath string, remotePath string) {
//    srcFile, err := os.Open(localFilePath)
//    if err != nil {
//        fmt.Println("os.Open error : ", localFilePath)
//        log.Fatal(err)
//    }
//    defer srcFile.Close()

//    var remoteFileName = path.Base(localFilePath)

//    dstFile, err := client.Create(path.Join(remotePath, remoteFileName))
//    if err != nil {
//        fmt.Println("sftpClient.Create error : ", path.Join(remotePath, remoteFileName))
//        log.Fatal(err)

//    }
//    defer dstFile.Close()

//    ff, err := ioutil.ReadAll(srcFile)
//    if err != nil {
//        fmt.Println("ReadAll error : ", localFilePath)
//        log.Fatal(err)

//    }
//    dstFile.Write(ff)
//    fmt.Println(localFilePath + "  copy file to remote server finished!")
//}

//func refresh() {
//    for {
//        time.Sleep(time.Second)
//        refreshFiles()
//    }
//}

//func uploadDirectory(localPath string, remotePath string) {
//    localFiles, err := ioutil.ReadDir(localPath)
//    if err != nil {
//        log.Fatal("read dir list fail ", err)
//    }

//    for _, backupDir := range localFiles {
//        localFilePath := path.Join(localPath, backupDir.Name())
//        remoteFilePath := path.Join(remotePath, backupDir.Name())
//        if backupDir.IsDir() {
//            client.Mkdir(remoteFilePath)
//            uploadDirectory(localFilePath, remoteFilePath)
//        } else {
//            uploadFile(path.Join(localPath, backupDir.Name()), remotePath)
//        }
//    }

//    fmt.Println(localPath + "  copy directory to remote server finished!")
//}

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

func executor(in string) {
	in = strings.TrimSpace(in)
	args := strings.Split(in, " ")
	// 没有指令
	if len(args) < 1 {
		return
	}
	if sftpCtx.Cmds[args[0]] == nil {
		fmt.Println("no such command: ", args[0])
	}
	if err := sftpCtx.Cmds[args[0]].Execute(args[1:]); err != nil {
		fmt.Println("err:", err)
	}
	return
}

//func executor(in string) {
//    in = strings.TrimSpace(in)
//    args := strings.Split(in, " ")
//    if len(args) < 1 {
//        return
//    }
//    if len(args) == 1 && args[0] == "exit" {
//        os.Exit(0)
//    }
//    if len(args) == 1 && args[0] == "pwd" {
//        fmt.Println(Cwd)
//        return
//    }
//    if len(args) == 1 && args[0] == "cd" {
//        pwd, err := client.Getwd()
//        if err != nil {
//            fmt.Println("err:", err)
//            return
//        }
//        Cwd = pwd
//        FileSet = make(map[string]bool)
//        FileSuggests = []prompt.Suggest{}
//        return
//    }
//    if len(args) == 1 {
//        return
//    }
//    switch args[0] {
//    case "get":
//        err := GetFile(client.Join(Cwd, args[1]))
//        if err != nil {
//            fmt.Println("err:", err)
//            return
//        }
//    case "put":
//    case "cd":
//        target := args[1]
//        current, err := GetCwd()
//        if err != nil {
//            fmt.Println("err:", err)
//            return
//        }
//        fmt.Println(current, " -> ", target)
//        Cwd = client.Join(current, target)
//        FileSet = make(map[string]bool)
//        FileSuggests = []prompt.Suggest{}
//    default:
//    }
//}

//func refreshFiles() {
//    pwd, err := GetCwd()
//    if err != nil {
//        return
//    }
//    walker := client.Walk(pwd)
//    walker.Step()
//    for walker.Step() {
//        dir := walker.Stat()
//        if FileSet[dir.Name()] {
//            continue
//        } else {
//            FileSet[dir.Name()] = true
//        }
//        filetype := "file"
//        if dir.IsDir() {
//            walker.SkipDir()
//            filetype = "directory"
//        }
//        FileSuggests = append(FileSuggests, prompt.Suggest{Text: dir.Name(), Description: filetype})
//        //FileSuggests = append(FileSuggests, prompt.Suggest{Text: walker.Path()})
//    }
//}

// complete source : remote file, local file, cmd name, option name, error source(? complete nothing and display error)
// line wrapper: first word, options,

func completer(in prompt.Document) []prompt.Suggest {
	line := in.CurrentLine()
	//args := strings.Split(line, " ")
	//if len(args) == 0 {
	//    return nil
	//}
	//position := in.CursorPositionCol()
	//if position == len(args[0]) {
	//    return prompt.FilterHasPrefix(CmdSuggests, in.GetWordBeforeCursor(), true)
	//}
	//cmdCompleter := Completer{Source: CmdSource{}}.Of("^(ls|get)")
	fileCompleter := Completer{Source: &FileSource{Connection: sftpCtx}}.Of("^(ls\\s+|get\\s+)")
	//if cmdCompleter.Match(line) {
	//    return prompt.FilterHasPrefix(cmdCompleter.Source.Get(), in.GetWordBeforeCursor(), true)
	//}
	if fileCompleter.Match(line) {
		fileCompleter.Source.Refresh()
		return prompt.FilterHasPrefix(fileCompleter.Source.Get(), in.GetWordBeforeCursor(), true)
	}
	return nil
	//line := in.CurrentLineBeforeCursor()
	//args := strings.Split(line, " ")
	//if len(args) <= 1 && len(line) > 0 && line[len(line)-1] != ' ' {
	//    return prompt.FilterHasPrefix(CmdSuggests, in.GetWordBeforeCursor(), true)
	//}
	//return prompt.FilterHasPrefix(FileSuggests, in.GetWordBeforeCursor(), true)
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
	var err error
	sftpCtx, err = NewSftp(args[1])
	if err != nil {
		fmt.Println("err:", err)
		return
	}
	p := prompt.New(
		executor,
		completer,
		//prompt.OptionLivePrefix(livePrefix),
		prompt.OptionPrefix(args[1]+"> "),
		prompt.OptionTitle("sftp-prompt"),
	)
	p.Run()
}
