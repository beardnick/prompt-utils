package main

import "github.com/c-bata/go-prompt"

var (
	FileSuggests = []prompt.Suggest{}
	CmdSuggests  = []prompt.Suggest{}
)

//Available commands:
//bye                                Quit sftp
//cd path                            Change remote directory to 'path'
//chgrp grp path                     Change group of file 'path' to 'grp'
//chmod mode path                    Change permissions of file 'path' to 'mode'
//chown own path                     Change owner of file 'path' to 'own'
//df [-hi] [path]                    Display statistics for current directory or
//                                   filesystem containing 'path'
//exit                               Quit sftp
//get [-afPpRr] remote [local]       Download file
//reget [-fPpRr] remote [local]      Resume download file
//reput [-fPpRr] [local] remote      Resume upload file
//help                               Display this help text
//lcd path                           Change local directory to 'path'
//lls [ls-options [path]]            Display local directory listing
//lmkdir path                        Create local directory
//ln [-s] oldpath newpath            Link remote file (-s for symlink)
//lpwd                               Print local working directory
//ls [-1afhlnrSt] [path]             Display remote directory listing
//lumask umask                       Set local umask to 'umask'
//mkdir path                         Create remote directory
//progress                           Toggle display of progress meter
//put [-afPpRr] local [remote]       Upload file
//pwd                                Display remote working directory
//quit                               Quit sftp
//rename oldpath newpath             Rename remote file
//rm path                            Delete remote file
//rmdir path                         Remove remote directory
//symlink oldpath newpath            Symlink remote file
//version                            Show SFTP version
//!command                           Execute 'command' in local shell
//!                                  Escape to local shell
//?                                  Synonym for help
