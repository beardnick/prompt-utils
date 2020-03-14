package main

import (
	"fmt"
	"os"
	"path"
)

// cmd name, description, option
// option name, description

type ICmd interface {
	Init(sftp *Sftp)
	Sftp() *Sftp
	Name() string
	Description() string
	Options() []Option
	Execute(args []string) error
}

type BaseCmd struct {
	context *Sftp
}

func (b *BaseCmd) Init(sftp *Sftp) {
	b.context = sftp
}

func (b *BaseCmd) Sftp() *Sftp {
	return b.context
}

func (b *BaseCmd) Name() string {
	panic("implement me")
}

func (b *BaseCmd) Description() string {
	panic("implement me")
}

func (b *BaseCmd) Options() []Option {
	panic("implement me")
}

func (b *BaseCmd) Execute(args []string) error {
	panic("implement me")
}

type Ls struct {
	BaseCmd
}

func (c *Ls) Name() string {
	return "ls"
}

func (c *Ls) Description() string {
	return "Display remote directory listing"
}

func (c *Ls) Options() []Option {
	panic("implement me")
}

func (c *Ls) Execute(args []string) error {
	path := c.Sftp().Client.Join(c.Sftp().RemotePath(), args[0])
	// #TODO: 20-03-14 这里的log怎么总是打不出来 //
	//log.Println("path:", path)
	files := c.Sftp().RemotePathFiles(path)
	for _, v := range files {
		fmt.Println(v.Name())
	}
	return nil
}

type Exit struct {
	BaseCmd
}

func (c *Exit) Name() string {
	return "exit"
}

func (c *Exit) Description() string {
	return "Quit sftp"
}

func (c *Exit) Options() []Option {
	panic("not implemented")
}

func (c *Exit) Execute(args []string) error {
	os.Exit(0)
	return nil
}

type Quit struct {
	Exit
}

func (c *Quit) Name() string {
	return "quit"
}

type Bye struct {
	Exit
}

func (c *Bye) Name() string {
	return "bye"
}

type Option struct {
}

type Get struct {
	BaseCmd
}

func (c *Get) Name() string {
	return "get"
}

func (c *Get) Description() string {
	return "Download file"
}

func (c *Get) Options() []Option {
	panic("not implemented")
}

// #TODO: 25-02-20 实现文件下载 //
func (c *Get) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("file name needed")
	}
	//return GetFile(args[0])
	fmt.Println("get ", args[0])
	return nil
}

func GetFile(file string) error {
	remote, err := sftpCtx.Client.Open(file)
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

type Cd struct {
	BaseCmd
}

func (c *Cd) Name() string {
	return "cd"
}

func (c *Cd) Description() string {
	return "Change remote directory to 'path'"
}

func (c *Cd) Options() Option {
	panic("not implemented")
}

func (c *Cd) Execute(args []string) error {
	panic("not implemented")
}
