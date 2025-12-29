package git

import (
	"io"
)

type Server struct {
	RootLocation string
}

type RefUpdate struct {
    RefName string  
    OldHash string  
    NewHash string  
}

type FetchRequest struct {
    Wants []string 
    Haves []string
}

type GitServer interface {
	List(repoPath string) (string, error)
	Push(repoPath string, packfileReader io.Reader, refUpdates []RefUpdate) error
	Fetch(repoPath string, req FetchRequest) (io.Reader, error)
}
