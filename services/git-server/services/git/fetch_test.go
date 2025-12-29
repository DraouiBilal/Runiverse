package git

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/packfile"
)

func TestFetch(t *testing.T) {
	// Setup: Create a server with a populated repo
	server := &Server{RootLocation: "/home/drale/work/trash"}

	// Get the commit hash from your test repo
	repo, _ := git.PlainOpen("/home/drale/work/trash/test_repo/.git")
	ref, _ := repo.Head()
	wantHash := ref.Hash().String()

	t.Logf("Fetching commit: %s", wantHash)

	// Execute Fetch
	req := FetchRequest{
		Wants: []string{wantHash},
		Haves: []string{}, // Client has nothing
	}

	packfileReader, err := server.Fetch("test_repo/.git", req)
	if err != nil {
		t.Fatalf("Fetch failed: %v", err)
	}

	// Verify: Write to a new repo and check we can clone it
	targetPath := "/tmp/fetch-test-target"
	os.RemoveAll(targetPath)
	os.MkdirAll(targetPath, 0755)

	// Initialize empty repo
	targetRepo, _ := git.PlainInit(targetPath, true)

	// Read the packfile into the target
	packBytes, _ := io.ReadAll(packfileReader)
	t.Logf("Generated packfile size: %d bytes", len(packBytes))

	// Unpack it
	err = packfile.UpdateObjectStorage(targetRepo.Storer, bytes.NewReader(packBytes))
	if err != nil && err != io.EOF {
		t.Fatalf("Failed to unpack: %v", err)
	}

	// Update the ref
	newRef := plumbing.NewHashReference("refs/heads/main", plumbing.NewHash(wantHash))
	targetRepo.Storer.SetReference(newRef)

	// Add this: Set HEAD to point to main
	headRef := plumbing.NewSymbolicReference(plumbing.HEAD, "refs/heads/main")
	targetRepo.Storer.SetReference(headRef)

	// Verify we can read the commit
	commit, err := targetRepo.CommitObject(plumbing.NewHash(wantHash))
	if err != nil {
		t.Fatalf("Commit not found after fetch: %v", err)
	}

	t.Logf("✓ Fetch successful! Commit message: %s", commit.Message)
}
