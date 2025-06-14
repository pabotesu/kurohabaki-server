package wg

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// ensureCleanInterface deletes an existing WireGuard interface with the same name
func ensureCleanInterface(name string) error {
	cmd := exec.Command("ip", "link", "del", name)
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "Cannot find device") {
		return fmt.Errorf("failed to delete existing interface: %s (%v)", string(output), err)
	}
	return nil
}

// SetupInterface creates and configures the WireGuard interface
func SetupInterface(interfaceName, privateKey, khIP string, port int, routeCIDR string) error {

	// Step 0: Ensure the interface name is valid
	if err := ensureCleanInterface(interfaceName); err != nil {
		return fmt.Errorf("interface cleanup failed: %v", err)
	}

	// Step 1: Create interface
	if err := exec.Command("ip", "link", "add", "dev", interfaceName, "type", "wireguard").Run(); err != nil {
		return fmt.Errorf("failed to create WireGuard interface: %w", err)
	}

	// Step 2: Parse private key
	key, err := wgtypes.ParseKey(privateKey)
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}

	// Step 3: Configure wg device
	client, err := wgctrl.New()
	if err != nil {
		return fmt.Errorf("failed to open wgctrl: %w", err)
	}
	defer client.Close()

	cfg := wgtypes.Config{
		PrivateKey:   &key,
		ListenPort:   &port,
		ReplacePeers: true,
	}

	if err := client.ConfigureDevice(interfaceName, cfg); err != nil {
		return fmt.Errorf("failed to configure WireGuard interface: %w", err)
	}

	// Step 4: Assign IP address to the interface
	ipWithCIDR := khIP
	if !strings.Contains(khIP, "/") {
		ipWithCIDR = khIP + "/32"
	}
	if err := exec.Command("ip", "address", "add", ipWithCIDR, "dev", interfaceName).Run(); err != nil {
		return fmt.Errorf("failed to assign IP address: %w", err)
	}

	// Step 5: Bring interface up
	if err := exec.Command("ip", "link", "set", "up", "dev", interfaceName).Run(); err != nil {
		return fmt.Errorf("failed to bring up interface: %w", err)
	}

	// Step 6: Add route to allow traffic to kh_nw_range (e.g., 100.100.96.0/24)
	if err := exec.Command("ip", "route", "add", routeCIDR, "dev", interfaceName).Run(); err != nil {
		return fmt.Errorf("failed to add route to %s: %w", routeCIDR, err)
	}

	return nil
}

func TeardownInterface(interfaceName string) error {
	cmd := exec.Command("ip", "link", "del", interfaceName)
	output, err := cmd.CombinedOutput()
	if err != nil && !strings.Contains(string(output), "Cannot find device") {
		return fmt.Errorf("failed to delete interface %s: %s (%v)", interfaceName, string(output), err)
	}
	return nil
}
