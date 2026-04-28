package proxman

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/luthermonson/go-proxmox"
)

// Creates a new Proxmox API client and initializes all the things related to it
func NewClient(apiBaseURL, apiTokenID, apiTokenSecret string) (client *Client, err error) {
	var (
		tlsConf    *tls.Config     = &tls.Config{InsecureSkipVerify: true}
		transport  *http.Transport = &http.Transport{TLSClientConfig: tlsConf}
		httpClient *http.Client    = &http.Client{
			Transport: transport,
			Timeout:   time.Second * 30,
		}
	)

	client = &Client{
		client: proxmox.NewClient(
			apiBaseURL,
			proxmox.WithHTTPClient(httpClient),
			proxmox.WithAPIToken(apiTokenID, apiTokenSecret),
		),
		lockPool: &sync.Pool{
			New: func() (mu any) {
				mu = &sync.Mutex{}
				return
			},
		},
		bg:         context.Background(),
		nodes:      make([]*proxmox.Node, 0),
		httpClient: httpClient,
	}

	if client.cluster, err = client.client.Cluster(client.bg); err != nil {
		err = fmt.Errorf(errGetClusterInfo, err)
		return
	}

	if err = client.refreshNodes(); err != nil {
		return
	}

	return
}

// Cleans up the client and all the things related to it
func (c *Client) Close() (err error) {
	if c.taskQueue != nil {
		if err = c.taskQueue.Close(); err != nil {
			return
		}
	}

	return
}

// Query nodes from the Proxmox API and update the client's node list
func (c *Client) refreshNodes() (err error) {
	var nodeStatuses []*proxmox.NodeStatus
	if nodeStatuses, err = c.client.Nodes(c.bg); err != nil {
		err = fmt.Errorf(errGetNodeStatuses, err)
		return
	}

	for _, status := range nodeStatuses {
		if status.Status == "online" {
			var node *proxmox.Node
			if node, err = c.client.Node(c.bg, status.Node); err != nil {
				err = fmt.Errorf(errGetNodeNamed, status.Node, err)
				return
			}

			c.nodes = append(c.nodes, node)
		}
	}

	return
}
