package handler

import (
	"github.com/fernandoocampo/grpcgwgokit/internal/domain"
	pb "github.com/fernandoocampo/grpcgwgokit/pkg/proto/grpcgwgokit/pb"
)

// NewEcho new echo parameter data.
type NewEcho struct {
	echoPB *pb.StringMessage
}

// NewNewEcho creates a new echo
func NewNewEcho(echoPB *pb.StringMessage) *NewEcho {
	return &NewEcho{
		echoPB: echoPB,
	}
}

func (n *NewEcho) toEcho() *domain.Echo {
	if n == nil || n.echoPB == nil {
		return nil
	}
	return &domain.Echo{
		Value: n.echoPB.Value,
	}
}
