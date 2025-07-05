package wg

import (
	"errors"
	"net"
	"testing"

	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// テスト可能な部分のみを切り出した関数
func parseWgParams(pubKey, ip string) (wgtypes.Key, net.IP, error) {
	key, err := wgtypes.ParseKey(pubKey)
	if err != nil {
		return key, nil, err
	}

	allowedIP := net.ParseIP(ip)
	if allowedIP == nil {
		return key, nil, errors.New("invalid IP address")
	}

	return key, allowedIP, nil
}

func TestParseWgParams(t *testing.T) {
	tests := []struct {
		name    string
		pubKey  string
		ip      string
		wantErr bool
	}{
		{
			name:    "valid inputs",
			pubKey:  "iS0vVH8S3HkCeWoJ5mZNSd9mQQ0xCLmtWEMMGOfq6kQ=", // public key example
			ip:      "10.0.0.3",
			wantErr: false,
		},
		{
			name:    "invalid public key",
			pubKey:  "invalid-key",
			ip:      "10.0.0.3",
			wantErr: true,
		},
		{
			name:    "invalid IP",
			pubKey:  "iS0vVH8S3HkCeWoJ5mZNSd9mQQ0xCLmtWEMMGOfq6kQ=",
			ip:      "invalid-ip",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := parseWgParams(tt.pubKey, tt.ip)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseWgParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
