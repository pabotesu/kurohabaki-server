package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
	"gopkg.in/yaml.v3"

	"github.com/pabotesu/kurohabaki-server/config"
	"github.com/pabotesu/kurohabaki-server/internal/clientconfig"
	"github.com/pabotesu/kurohabaki-server/internal/etcd"
	"github.com/pabotesu/kurohabaki-server/internal/ipalloc"
	"github.com/pabotesu/kurohabaki-server/internal/wg"
	// ←実際の import パスに置き換えてください
	// server public keyを生成してる場所
)

var nodeAddCmd = &cobra.Command{
	Use:   "add <public-key>",
	Short: "Register a new node by its WireGuard public key",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		publicKey := args[0]

		// Load configuration
		cfg, err := config.LoadConfig(GetConfigPath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}

		// Initialize etcd client
		etcdClient, err := etcd.New(cfg.ServerSettings)
		if err != nil {
			log.Fatalf("Failed to connect to etcd: %v", err)
		}

		// Check if the public key is already registered
		exists, err := etcdClient.NodeExists(publicKey)
		if err != nil {
			log.Fatalf("Failed to check existing node: %v", err)
		}
		if exists {
			log.Fatalf("Public key %s is already registered", publicKey)
		}

		// Get already allocated IPs
		existingIPs, err := etcdClient.GetAllNodeIPs()
		if err != nil {
			log.Fatalf("Failed to fetch node list from etcd: %v", err)
		}

		// Allocate a new IP
		allocatedIP, err := ipalloc.Allocate(cfg.ServerSettings.KhNwRange, existingIPs)
		if err != nil {
			log.Fatalf("Failed to allocate IP: %v", err)
		}

		// Register node
		err = etcdClient.RegisterNode(publicKey, allocatedIP)
		if err != nil {
			log.Fatalf("Failed to register node in etcd: %v", err)
		}

		// Generate client YAML config
		var clientCfg clientconfig.ClientYAMLConfig

		clientCfg.Interface.PrivateKey = "<YOUR_PRIVATE_KEY_HERE>"
		clientCfg.Interface.Address = allocatedIP + "/32"
		clientCfg.Interface.DNS = cfg.ServerSettings.KhIP
		clientCfg.Interface.Routes = []string{cfg.ServerSettings.KhNwRange}

		// Server pubkey
		privKey, err := wgtypes.ParseKey(cfg.ServerSettings.PrivateKey)
		if err != nil {
			log.Fatalf("Failed to parse server private key: %v", err)
		}
		clientCfg.Peer.PublicKey = privKey.PublicKey().String()
		// Server endpoint and allowed IPs
		clientCfg.Peer.Endpoint = fmt.Sprintf("%s:%d",
			cfg.ServerSettings.PublicIP,
			cfg.ServerSettings.Port)
		// Allowed IPs for the peer
		// This should be the range of the Kurohabaki network
		// (e.g.,
		clientCfg.Peer.AllowedIPs = cfg.ServerSettings.KhIP + "/32"
		// Persistent keepalive for the peer
		clientCfg.Peer.PersistentKeepalive = 5 // seconds
		// Etcd endpoint
		clientCfg.Etcd.Endpoint = fmt.Sprintf("%s:%d",
			cfg.ServerSettings.Etcd.AdvertiseClientIP,
			cfg.ServerSettings.Etcd.AdvertiseClientPort)

		out, err := yaml.Marshal(&clientCfg)
		if err != nil {
			log.Fatalf("Failed to marshal client config: %v", err)
		}

		fmt.Println("\n# Client YAML configuration")
		fmt.Println(string(out))

		err = wg.AddPeer(cfg.ServerSettings.InterfaceName, publicKey, allocatedIP)
		if err != nil {
			log.Printf("Failed to add peer to WireGuard: %v", err)
		} else {
			log.Printf("Added peer %s with IP %s to WireGuard", publicKey, allocatedIP)
		}
	},
}

func init() {
	rootCmd.AddCommand(nodeAddCmd)
}
