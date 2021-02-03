package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/c-bata/go-prompt"
)

var (
	sftpCtx *Sftp
)

func executor(in string) {
	in = strings.TrimSpace(in)
	args := strings.Split(in, " ")
	// 没有指令
	if len(args) < 1 {
		return
	}
	if sftpCtx.Cmds[args[0]] == nil {
		if strings.TrimSpace(args[0]) != "" {
			fmt.Println("no such command:", args[0])
		}
		return
	}
	if err := sftpCtx.Cmds[args[0]].Execute(args[1:]); err != nil {
		fmt.Println("err:", err)
	}
	return
}

func completer(in prompt.Document) []prompt.Suggest {
	line := in.CurrentLine()
	fileCompleter := Completer{Source: &FileSource{Connection: sftpCtx}}.Of("^(ls\\s+|get\\s+)")
	if fileCompleter.Match(line) {
		fileCompleter.Source.Refresh()
		return prompt.FilterHasPrefix(fileCompleter.Source.Get(), in.GetWordBeforeCursor(), true)
	}
	return nil
}

func main() {
	// 设置日志头
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	f, err := os.OpenFile(time.Now().Format("2006-01-02")+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()
	log.SetOutput(f)
	//var err error
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
