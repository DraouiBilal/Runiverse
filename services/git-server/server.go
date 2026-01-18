package main

import (
	"context"
	"fmt"
	"log"
	"io"
	cri "github.com/DraouiBilal/Runiverse-cri/git/v1"
	"github.com/DraouiBilal/Runiverse/git-server/services/git"
)

type Server struct {
	cri.UnimplementedGitServiceServer
	GitServer git.GitServer
}

func (s *Server) List(ctx context.Context, req *cri.ListRequest) (*cri.ListResponse, error) {
	refs, err := s.GitServer.List(req.RepoPath)
	if err != nil {
		return nil, err
	}

	return &cri.ListResponse{Refs: refs}, nil
}

func (s *Server) Push(stream cri.GitService_PushServer) error {
	// Receive first message (should be metadata)
	firstReq, err := stream.Recv()
	if err != nil {
		return err
	}

	metadata := firstReq.GetMetadata()
	if metadata == nil {
		return fmt.Errorf("first message must contain metadata")
	}

	// Convert proto RefUpdates to internal type
	var refUpdates []git.RefUpdate
	for _, ru := range metadata.RefUpdates {
		refUpdates = append(refUpdates, git.RefUpdate{
			RefName: ru.RefName,
			OldHash: ru.OldHash,
			NewHash: ru.NewHash,
		})
	}

	// Create pipe for streaming packfile data
	reader, writer := io.Pipe()

	// Goroutine to receive remaining packfile chunks
	errChan := make(chan error, 1)
	go func() {
		defer writer.Close()
		for {
			req, err := stream.Recv()
			if err == io.EOF {
				errChan <- nil
				return
			}
			if err != nil {
				errChan <- err
				return
			}

			// Write packfile chunks to pipe
			if chunk := req.GetPackfileChunk(); len(chunk) > 0 {
				if _, err := writer.Write(chunk); err != nil {
					errChan <- err
					return
				}
			}
		}
	}()

	// Call git Push (blocks until reader is fully consumed)
	pushErr := s.GitServer.Push(metadata.RepoPath, reader, refUpdates)

	// Wait for stream to finish
	streamErr := <-errChan

	// Handle errors
	if pushErr != nil {
		return stream.SendAndClose(&cri.PushResponse{
			Status: "ERROR",
			Err:    pushErr.Error(),
		})
	}
	if streamErr != nil {
		return streamErr
	}

	return stream.SendAndClose(&cri.PushResponse{Status: "OK"})
}

func (s *Server) Fetch(req *cri.FetchRequest, stream cri.GitService_FetchServer) error {
    fetchReq := git.FetchRequest{
        Wants: req.Wants,
        Haves: req.Haves,
    }
    
    packfileReader, err := s.GitServer.Fetch(req.RepoPath, fetchReq)
    if err != nil {
        return err
    }
    
    buffer := make([]byte, 32*1024) // 32KB chunks
	totalBytes := 0
    
    for {
        n, err := packfileReader.Read(buffer)
        
        // Send data if we read any
        if n > 0 {
            if sendErr := stream.Send(&cri.FetchResponse{
                PackfileChunk: buffer[:n],
            }); sendErr != nil {
                return sendErr
            }
        }
        
        // Check for EOF or errors AFTER processing data
        if err == io.EOF {
            log.Printf("Fetch complete: sent %d bytes total\n", totalBytes)
            break
        }
        
        if err != nil {
            return fmt.Errorf("read error: %w", err)
        }
    }
    
    return nil
}
