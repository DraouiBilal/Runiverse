package git

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/packfile"
	"io"
)

func (s *Server) Push(repoPath string, packfileReader io.Reader, refUpdates []RefUpdate) error {

	fullPath := fmt.Sprintf("%s/%s", s.RootLocation, repoPath)

	repo, err := git.PlainOpen(fullPath)

	if err != nil {
		return fmt.Errorf("failed to open repo: %w", err)
	}

	if err := packfile.UpdateObjectStorage(repo.Storer, packfileReader); err != nil && err != io.EOF {
		return fmt.Errorf("failed to process packfile: %w", err)
	}

	for _, update := range refUpdates {
		refName := plumbing.ReferenceName(update.RefName)
		newHash := plumbing.NewHash(update.NewHash)

		// Optional: Verify the old hash matches (optimistic locking)
		if update.OldHash != "" {
			currentRef, err := repo.Reference(refName, true)
			if err == nil && currentRef.Hash().String() != update.OldHash {
				return fmt.Errorf("ref %s has changed (expected %s, got %s)",
					update.RefName, update.OldHash, currentRef.Hash().String())
			}
		}

		newRef := plumbing.NewHashReference(refName, newHash)
		if err := repo.Storer.SetReference(newRef); err != nil {
			return fmt.Errorf("failed to update ref %s: %w", update.RefName, err)
		}
	}

	return nil
}
