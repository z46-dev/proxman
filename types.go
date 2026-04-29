package proxman

import (
	"context"
	"net/http"
	"sync"

	"github.com/luthermonson/go-proxmox"
)

type Client struct {
	lockPool   *sync.Pool
	bg         context.Context
	httpClient *http.Client

	client  *proxmox.Client
	nodes   []*proxmox.Node
	cluster *proxmox.Cluster
}
