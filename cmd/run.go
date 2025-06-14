package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"github.com/pabotesu/kurohabaki-server/config"
	"github.com/pabotesu/kurohabaki-server/internal/etcd"
	"github.com/pabotesu/kurohabaki-server/internal/monitor"
	"github.com/pabotesu/kurohabaki-server/internal/wg"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Start kurohabaki-server as a background coordination service",
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		cfg, err := config.LoadConfig(GetConfigPath())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
			os.Exit(1)
		}

		log.Println("kurohabaki-server starting...")
		log.Printf("Server will listen on %s:%d\n", cfg.ServerSettings.PublicIP, cfg.ServerSettings.Port)
		log.Printf("WireGuard interface IP: %s\n", cfg.ServerSettings.KhIP)
		log.Printf("Allocated subnet range: %s\n", cfg.ServerSettings.KhNwRange)

		iface := cfg.ServerSettings.InterfaceName
		if err := wg.SetupInterface(iface, cfg.ServerSettings.PrivateKey, cfg.ServerSettings.KhIP, cfg.ServerSettings.Port, cfg.ServerSettings.KhNwRange); err != nil {
			log.Fatalf("WireGuard setup failed: %v", err)
		}
		log.Printf("WireGuard interface '%s' configured\n", iface)

		// Start etcd
		if err := etcd.StartEtcd(cfg.ServerSettings.Etcd); err != nil {
			log.Fatalf("Failed to start etcd: %v", err)
		}
		log.Printf("etcd endpoint: %s:%d\n", cfg.ServerSettings.Etcd.EndpointIP, cfg.ServerSettings.Etcd.EndpointPort)

		// Derive server public key from private key
		privKey, err := wgtypes.ParseKey(cfg.ServerSettings.PrivateKey)
		if err != nil {
			log.Fatalf("Invalid private key: %v", err)
		}
		serverPubKey := privKey.PublicKey().String()

		// Register self in etcd
		etcdClient, err := etcd.New(cfg.ServerSettings)
		if err != nil {
			log.Fatalf("Failed to connect to etcd: %v", err)
		}
		err = etcdClient.RegisterNode(serverPubKey, cfg.ServerSettings.KhIP)
		if err != nil {
			log.Fatalf("Failed to register self in etcd: %v", err)
		}
		log.Printf("Registered self (%s) with IP %s in etcd\n", serverPubKey, cfg.ServerSettings.KhIP)

		// - etcd sync
		nodeList, err := etcdClient.GetAllNodeIPs()
		if err != nil {
			log.Fatalf("Failed to get node list: %v", err)
		}
		for pubKey, ip := range nodeList {
			if pubKey == serverPubKey {
				continue // skip self
			}

			if err := wg.AddPeer(iface, pubKey, ip); err != nil {
				log.Printf("Failed to add peer %s: %v", pubKey, err)
				continue
			}
			log.Printf("Added peer %s with IP %s\n", pubKey, ip)
		}
		// - peer update loop
		go monitor.ObserveAndSyncEndpoints(iface, serverPubKey, etcdClient, 1*time.Second)
		log.Println("Started peer endpoint observation loop")

		// create sugnal channel to handle graceful shutdown
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

		// 終了時の後始末をゴルーチンで待機
		go func() {
			sig := <-sigCh
			log.Printf("Received signal: %s", sig)

			log.Println("Stopping embedded etcd...")
			if err := etcd.StopEtcd(); err != nil {
				log.Printf("Failed to stop etcd cleanly: %v", err)
			} else {
				log.Println("etcd stopped successfully.")
			}

			log.Println("leaning up WireGuard interface...")
			if err := wg.TeardownInterface(iface); err != nil {
				log.Printf("⚠ Failed to teardown interface: %v", err)
			} else {
				log.Printf("Interface %s deleted", iface)
			}

			log.Println("Exiting...")
			os.Exit(0)
		}()

		// Placeholder for future background loop
		select {}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
