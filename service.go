package storage

import "io"

// Service provides interface for bridging between API and external system client
type Service struct {
	client Client
}

// NewService instantiates a new service
func NewService(client Client) Service {
	return Service{
		client: client,
	}
}

// Upload performs the upload
func (s *Service) Upload(folder, key string, file io.ReadSeeker) (string, error) {
	return s.client.Put(folder, key, file)
}
