package storage

import (
	"io"
)

// Client is the client interacting with external storage system
type Client interface {
	Put(string, io.Reader) error
}
