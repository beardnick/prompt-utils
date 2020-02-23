package main

import (
	"fmt"
	"os"
	"path"
)

// cmd name, description, option
// option name, description

type ICmd interface {
	Name() string
	Description() string
	Options() []Option
	Execute(args []string) error
}

type IFlagCmd interface {
	ICmd
	Init() ICmd
}

type Ls struct {
}

func (c Ls) Name() string {
	return "ls"
}

func (c Ls) Description() string {
	return "Display remote directory listing"
}

func (c Ls) Options() []Option {
	panic("implement me")
}

func (c Ls) Execute(args []string) error {
	fmt.Println("exec Ls")
	return nil
}

type Exit struct {
}

func (c Exit) Name() string {
	return "exit"
}

func (c Exit) Description() string {
	return "Quit sftp"
}

func (c Exit) Options() []Option {
	panic("not implemented")
}

func (c Exit) Execute(args []string) error {
	os.Exit(0)
	return nil
}

type Quit struct {
	Exit
}

func (c Quit) Name() string {
	return "quit"
}

type Bye struct {
	Exit
}

func (c Bye) Name() string {
	return "bye"
}

type Option struct {
}

type Get struct {
}

func (c Get) Init() ICmd {
	// #TODO: 24-02-20 初始化flag //
	return c
}

func (c Get) Name() string {
	return "get"
}

func (c Get) Description() string {
	return "Download file"
}

func (c Get) Options() []Option {
	panic("not implemented")
}

func (c Get) Execute(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("file name needed")
	}
	return GetFile(args[0])
}

func GetFile(file string) error {
	remote, err := sftpInstance.Client.Open(file)
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
