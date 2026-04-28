package proxman

import (
	"context"
	"net/http"
	"sync"

	"github.com/hibiken/asynq"
	"github.com/luthermonson/go-proxmox"
)

type Client struct {
	lockPool   *sync.Pool
	taskQueue  *asynq.Client
	bg         context.Context
	httpClient *http.Client

	client  *proxmox.Client
	nodes   []*proxmox.Node
	cluster *proxmox.Cluster
}
