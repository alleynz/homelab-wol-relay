package config

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"
)

type Device struct {
	Name string `json:"name"`
	MAC  string `json:"mac"`
	IP   string `json:"ip"`
}

type Config struct {
	Listen  string   `json:"listen"`
	Port    int      `json:"port"`
	Devices []Device `json:"devices"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if cfg.Listen == "" {
		cfg.Listen = "0.0.0.0"
	}
	if cfg.Port == 0 {
		cfg.Port = 9
	}
	if cfg.Port < 1 || cfg.Port > 65535 {
		return nil, fmt.Errorf("invalid port %d", cfg.Port)
	}

	seen := map[string]struct{}{}
	for i := range cfg.Devices {
		d := &cfg.Devices[i]
		hw, err := net.ParseMAC(d.MAC)
		if err != nil {
			return nil, fmt.Errorf("device %q has invalid MAC %q: %w", d.Name, d.MAC, err)
		}
		if len(hw) != 6 {
			return nil, fmt.Errorf("device %q MAC must be 6 bytes", d.Name)
		}
		d.MAC = strings.ToLower(hw.String())

		ip := net.ParseIP(d.IP)
		if ip == nil || ip.To4() == nil {
			return nil, fmt.Errorf("device %q has invalid IPv4 address %q", d.Name, d.IP)
		}
		d.IP = ip.To4().String()

		if _, exists := seen[d.MAC]; exists {
			return nil, fmt.Errorf("duplicate MAC mapping %s", d.MAC)
		}
		seen[d.MAC] = struct{}{}
	}

	return &cfg, nil
}

func (c *Config) DeviceByMAC(mac net.HardwareAddr) (Device, bool) {
	key := strings.ToLower(mac.String())
	for _, d := range c.Devices {
		if d.MAC == key {
			return d, true
		}
	}
	return Device{}, false
}
