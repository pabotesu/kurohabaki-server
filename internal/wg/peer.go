package wg

import (
	"net"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// AddPeer adds a WireGuard peer to the interface
func AddPeer(interfaceName, pubKey, ip string) error {
	key, err := wgtypes.ParseKey(pubKey)
	if err != nil {
		return err
	}

	allowedIP := net.ParseIP(ip)
	if allowedIP == nil {
		return err
	}

	peer := wgtypes.PeerConfig{
		PublicKey: key,
		AllowedIPs: []net.IPNet{
			{
				IP:   allowedIP,
				Mask: net.CIDRMask(32, 32),
			},
		},
	}

	client, err := wgctrl.New()
	if err != nil {
		return err
	}
	defer client.Close()

	return client.ConfigureDevice(interfaceName, wgtypes.Config{
		ReplacePeers: false,
		Peers:        []wgtypes.PeerConfig{peer},
	})
}
