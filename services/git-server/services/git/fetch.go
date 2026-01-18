package git

import (
	"bytes"
	"fmt"
	"io"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/packfile"
	"github.com/go-git/go-git/v5/plumbing/filemode"
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

	// Mark this object as collected
	collected[startHash] = true

	// Try to interpret as a commit
	commit, err := repo.CommitObject(startHash)
	if err == nil {
		// Recursively collect the tree and all its contents
		if err := s.collectTreeObjects(repo, commit.TreeHash, collected); err != nil {
			return err
		}

		// Recursively process parent commits
		for _, parentHash := range commit.ParentHashes {
			if !haveMap[parentHash] && !collected[parentHash] {
				if err := s.collectObjects(repo, parentHash, collected, haves); err != nil {
					return err
				}
			}
		}

		return nil
	}

	// Try as a tree
	tree, err := repo.TreeObject(startHash)
	if err == nil {
		return s.collectTreeObjects(repo, tree.Hash, collected)
	}

	// Try as a blob - already added to collected above
	_, err = repo.BlobObject(startHash)
	if err == nil {
		return nil
	}

	// Try as a tag
	tag, err := repo.TagObject(startHash)
	if err == nil {
		collected[tag.Target] = true
		return s.collectObjects(repo, tag.Target, collected, haves)
	}

	return fmt.Errorf("unknown object type for %s", startHash)
}

// New helper function to recursively collect tree objects
func (s *Server) collectTreeObjects(repo *git.Repository, treeHash plumbing.Hash,
	collected map[plumbing.Hash]bool) error {

	// Skip if already collected
	if collected[treeHash] {
		return nil
	}

	collected[treeHash] = true

	tree, err := repo.TreeObject(treeHash)
	if err != nil {
		return fmt.Errorf("failed to get tree %s: %w", treeHash, err)
	}

	// Iterate over all entries in the tree
	for _, entry := range tree.Entries {
		collected[entry.Hash] = true

		// If it's a tree (subdirectory), recursively collect it
		if entry.Mode == filemode.Dir {
			if err := s.collectTreeObjects(repo, entry.Hash, collected); err != nil {
				return err
			}
		}
		// Blobs (files) are already added above
	}

	return nil
}

