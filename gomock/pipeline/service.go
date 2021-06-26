package pipeline

import (
	"context"

	"github.com/nathanusask/backend-studies/gomock/datastore"
)

type Server struct {
	datastoreService datastore.Interface
}

func New(ds datastore.Interface) *Server {
	return &Server{datastoreService: ds}
}

func (s *Server) Run(ctx context.Context) error {

	return nil
}
