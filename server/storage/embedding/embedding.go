package embedding

import (
	"context"
	"time"

	qdrant "github.com/qdrant/go-client/qdrant"
	"github.com/vargspjut/wlog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Embedding struct {
	log wlog.Logger

	conn *grpc.ClientConn
}

func New(log wlog.Logger) *Embedding {
	return &Embedding{
		log: log,
	}
}

func (e *Embedding) Connect(addr string) error {
	// Set up gRPC connection to the server.
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		e.log.Error("failed to connect to addr: %s with error: %v", addr, err)
		return err
	}
	e.conn = conn

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// Check Qdrant version
	qdrantClient := qdrant.NewQdrantClient(conn)
	healthCheckResult, err := qdrantClient.HealthCheck(ctx, &qdrant.HealthCheckRequest{})
	if err != nil {
		e.log.Errorf("failed to health check qdrant: %v", err)
		return err
	} else {
		e.log.Infof("Qdrant version: %s", healthCheckResult.GetVersion())
	}
	return nil
}

func (e *Embedding) Disconnect() {
	if e.conn == nil {
		return
	}
	if err := e.conn.Close(); err != nil {
		e.log.Errorf("disconnect error: %v", err)
		return
	}
	e.log.Info("embedding engine disconnected")
}

func (e Embedding) Load(ctx context.Context, question string) (string, error) {
	_ = qdrant.NewCollectionsClient(e.conn)

	return "", nil
}
