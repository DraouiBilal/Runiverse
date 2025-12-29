package git

import (
	"bytes"
	"fmt"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/packfile"
	"github.com/go-git/go-git/v5/plumbing/object"
)

func (s *Server) Fetch(repoPath string, req FetchRequest) (io.Reader, error) {
	fullPath := fmt.Sprintf("%s/%s", s.RootLocation, repoPath)
	repo, err := git.PlainOpen(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repo: %w", err)
	}

	// Collect all objects we need to send
	// Convert the "wants" into a list of objects to pack
	objectsToPack := make(map[plumbing.Hash]bool)

	for _, wantHash := range req.Wants {
		hash := plumbing.NewHash(wantHash)

		if err := s.collectObjects(repo, hash, objectsToPack, req.Haves); err != nil {
			return nil, fmt.Errorf("failed to collect objects: %w", err)
		}
	}

	// Generate a packfile containing these objects
	var buf bytes.Buffer
	encoder := packfile.NewEncoder(&buf, repo.Storer, false) // false = use OFS deltas

	// Convert map to slice of hashes
	var hashes []plumbing.Hash
	for hash := range objectsToPack {
		hashes = append(hashes, hash)
	}

	// Encode the packfile
	_, err = encoder.Encode(hashes, 10) // 10 is the packWindow size (compression)
	if err != nil {
		return nil, fmt.Errorf("failed to encode packfile: %w", err)
	}

	return bytes.NewReader(buf.Bytes()), nil
}

// Helper function to collect all objects reachable from a commit
func (s *Server) collectObjects(repo *git.Repository, startHash plumbing.Hash,
	collected map[plumbing.Hash]bool, haves []string) error {

	// Convert "haves" to a map for quick lookup
	haveMap := make(map[plumbing.Hash]bool)
	for _, have := range haves {
		haveMap[plumbing.NewHash(have)] = true
	}

	// If the client already has this object, skip it
	if haveMap[startHash] {
		return nil
	}

	// If we've already collected this object, skip it
	if collected[startHash] {
		return nil
	}

	// We don't have this object, so let's mark it as collected
	collected[startHash] = true

	// Try to interpret as a commit (most common case)
	commit, err := repo.CommitObject(startHash)
	if err == nil {
		// Add the tree
		collected[commit.TreeHash] = true

		// Walk the tree to get all blobs
		tree, err := commit.Tree()
		if err == nil {
			tree.Files().ForEach(func(f *object.File) error {
				collected[f.Hash] = true
				return nil
			})
		}

		// Recursively process parent commits
		for _, parentHash := range commit.ParentHashes {
			if !haveMap[parentHash] {
				s.collectObjects(repo, parentHash, collected, haves)
			}
		}

		return nil
	}

	// If it's not a commit, we already added the object to collected
	return nil
}
