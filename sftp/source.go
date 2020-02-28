package main

import (
	"github.com/c-bata/go-prompt"
)

type ISource interface {
	Refresh()
	Get() []prompt.Suggest
}

type FileSource struct {
	source  []prompt.Suggest
	Connect *Sftp
}

func (s *FileSource) Refresh() {
	files := s.Connect.RemoteFiles()
	for _, v := range files {
		if v.IsDir() {
			s.source = append(s.source, prompt.Suggest{v.Name(), "dir"})
		} else {
			s.source = append(s.source, prompt.Suggest{v.Name(), "file"})
		}
	}
}

func (s *FileSource) Get() []prompt.Suggest {
	return s.source
}

type CmdSource struct {
}

func (s *CmdSource) Refresh() {
}

func (s *CmdSource) Get() []prompt.Suggest {
	return []prompt.Suggest{
		{"ls", "Display remote directory listing"},
		{"get", "Download file"},
	}
}
