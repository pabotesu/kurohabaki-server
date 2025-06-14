package ipalloc

import (
    "fmt"
    "net"
    "strings"
)

// Allocate finds an unused IP address from the given CIDR block
// usedIPs: map[publicKey]ip (e.g. "abc123" -> "100.100.96.2")
func Allocate(cidr string, usedIPs map[string]string) (string, error) {
    _, ipnet, err := net.ParseCIDR(cidr)
    if err != nil {
        return "", fmt.Errorf("invalid CIDR block: %w", err)
    }

    used := map[string]bool{}
    for _, ip := range usedIPs {
        used[strings.TrimSpace(ip)] = true
    }

    ip := ipnet.IP.To4()
    if ip == nil {
        return "", fmt.Errorf("only IPv4 is supported")
    }

    // Skip .0 (network), .1 (reserved for server), start from .2
    for i := 2; i < 255; i++ {
        candidate := net.IPv4(ip[0], ip[1], ip[2], byte(i)).String()
        if ipnet.Contains(net.ParseIP(candidate)) && !used[candidate] {
            return candidate, nil
        }
    }

    return "", fmt.Errorf("no available IPs in range %s", cidr)
}
