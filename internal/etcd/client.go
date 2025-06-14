package etcd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pabotesu/kurohabaki-server/config"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdClient struct {
	cli *clientv3.Client
}

const nodePrefix = "/kurohabaki/nodes"

// New creates a new EtcdClient from config
func New(cfg config.ServerSettings) (*EtcdClient, error) {
	endpoint := fmt.Sprintf("%s:%d", cfg.Etcd.EndpointIP, cfg.Etcd.EndpointPort)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &EtcdClient{cli: cli}, nil
}

// nodeIPKey constructs the etcd key for a node's IP address
func nodeIPKey(pubKey string) string {
	return fmt.Sprintf("%s/%s/ip", nodePrefix, pubKey)
}

// nodeLastSeenKey constructs the etcd key for a node's last seen timestamp
func nodeLastSeenKey(pubKey string) string {
	return fmt.Sprintf("%s/%s/last_seen", nodePrefix, pubKey)
}

// nodeEndpointKey constructs the etcd key for a node's endpoint
func nodeEndpointKey(pubKey string) string {
	return fmt.Sprintf("%s/%s/endpoint", nodePrefix, pubKey)
}

// NodeExists checks if a node with the given public key exists
func (e *EtcdClient) NodeExists(pubKey string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resp, err := e.cli.Get(ctx, nodeIPKey(pubKey))
	if err != nil {
		return false, err
	}
	return len(resp.Kvs) > 0, nil
}

// GetAllNodeIPs returns a map of publicKey -> ip
func (e *EtcdClient) GetAllNodeIPs() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	resp, err := e.cli.Get(ctx, nodePrefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for _, kv := range resp.Kvs {
		key := string(kv.Key) // e.g. /kurohabaki/nodes/<pubkey>/ip
		val := string(kv.Value)

		if strings.HasSuffix(key, "/ip") {
			parts := strings.Split(key, "/")
			if len(parts) >= 4 {
				pubKey := parts[3]
				result[pubKey] = val
			}
		}
	}

	return result, nil
}

// RegisterNode writes publicKey -> ip and timestamp
func (e *EtcdClient) RegisterNode(pubKey, ip string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	now := time.Now().UTC().Format(time.RFC3339)

	_, err := e.cli.Txn(ctx).
		Then(
			clientv3.OpPut(nodeIPKey(pubKey), ip),
			clientv3.OpPut(nodeLastSeenKey(pubKey), now),
		).Commit()
	return err
}

// UpdatePeerEndpoint updates the endpoint and last seen timestamp for a node
func (e *EtcdClient) UpdatePeerEndpoint(pubKey, endpoint string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	now := time.Now().UTC().Format(time.RFC3339)

	_, err := e.cli.Txn(ctx).Then(
		clientv3.OpPut(nodeEndpointKey(pubKey), endpoint),
		clientv3.OpPut(nodeLastSeenKey(pubKey), now),
	).Commit()

	return err
}
