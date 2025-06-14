package monitor

import (
	"log"
	"sync"
	"time"

	"github.com/pabotesu/kurohabaki-server/internal/etcd"
	"golang.zx2c4.com/wireguard/wgctrl"
)

// ObserveAndSyncEndpoints periodically fetches peer endpoints and stores them in etcd if changed
func ObserveAndSyncEndpoints(interfaceName string, serverPubKey string, etcdClient *etcd.EtcdClient, interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// 前回のエンドポイント情報を保持
	prevEndpoints := make(map[string]string)
	var mu sync.Mutex

	for {
		<-ticker.C

		client, err := wgctrl.New()
		if err != nil {
			log.Printf("Failed to open wgctrl: %v", err)
			continue
		}

		device, err := client.Device(interfaceName)
		client.Close()
		if err != nil {
			log.Printf("Failed to get device info: %v", err)
			continue
		}

		mu.Lock()
		for _, peer := range device.Peers {
			pubKey := peer.PublicKey.String()

			// Skip self
			if pubKey == serverPubKey {
				continue
			}

			endpoint := ""
			if peer.Endpoint != nil {
				endpoint = peer.Endpoint.String()
			}

			// 差分がある場合のみ etcd に反映
			if prev, ok := prevEndpoints[pubKey]; !ok || prev != endpoint {
				if err := etcdClient.UpdatePeerEndpoint(pubKey, endpoint); err != nil {
					log.Printf("Failed to update endpoint for %s: %v", pubKey, err)
				} else {
					prevEndpoints[pubKey] = endpoint
				}
			}
		}
		mu.Unlock()
	}
}
