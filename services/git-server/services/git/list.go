package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

// Returns an array of the needed refs for list command
func getRefs(repoPath string) ([]string, error) {
	r, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	refs, err := r.References()
	if err != nil {
		return nil, err
	}

	refsArray := []string{}

	refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Type() == plumbing.HashReference {
			refsArray = append(refsArray, fmt.Sprintf("%s %s", ref.Hash().String(), ref.Name()))
		}
		return nil
	})

	headRef, err := r.Reference(plumbing.HEAD, false) // false = don't resolve to hash

	if err == nil && headRef.Type() == plumbing.SymbolicReference {
		refsArray = append(refsArray, fmt.Sprintf("@%s HEAD", headRef.Target()))
	}

	refsArray = append(refsArray, fmt.Sprintln())
	return refsArray, nil
}

func (s *Server) List(repoPath string) (string, error) {
	refs, err := getRefs(s.RootLocation + "/" + repoPath)

	if err != nil {
		return "", err
	}

	result := ""

	for _, ref := range refs {
		result += ref + "\n"
	}

	return result, nil
}
