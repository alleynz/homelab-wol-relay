package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/example/wol-relay/internal/config"
	"github.com/example/wol-relay/internal/wol"
)

func main() {
	configPath := flag.String("config", "/config/config.json", "path to configuration file")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("configuration error: %v", err)
	}

	listenAddr := &net.UDPAddr{IP: net.ParseIP(cfg.Listen), Port: cfg.Port}
	conn, err := net.ListenUDP("udp4", listenAddr)
	if err != nil {
		log.Fatalf("listen on %s:%d: %v", cfg.Listen, cfg.Port, err)
	}
	defer conn.Close()

	log.Printf("WoL relay listening on %s:%d with %d configured device(s)", cfg.Listen, cfg.Port, len(cfg.Devices))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-stop
		_ = conn.Close()
	}()

	buffer := make([]byte, 2048)
	for {
		n, source, err := conn.ReadFromUDP(buffer)
		if err != nil {
			if opErr, ok := err.(*net.OpError); ok && !opErr.Temporary() {
				log.Printf("listener stopped")
				return
			}
			log.Printf("receive error: %v", err)
			continue
		}

		packet := append([]byte(nil), buffer[:n]...)
		mac, err := wol.ParseMagicPacket(packet)
		if err != nil {
			log.Printf("ignored UDP packet from %s: %v", source, err)
			continue
		}

		device, found := cfg.DeviceByMAC(mac)
		if !found {
			log.Printf("ignored WoL for unknown MAC %s from %s", mac, source)
			continue
		}

		outbound, err := wol.BuildMagicPacket(mac)
		if err != nil {
			log.Printf("cannot build WoL for %s: %v", mac, err)
			continue
		}

		if err := wol.SendUnicast(device.IP, cfg.Port, outbound); err != nil {
			log.Printf("failed relaying WoL for %s (%s) to %s:%d: %v", device.Name, mac, device.IP, cfg.Port, err)
			continue
		}

		name := device.Name
		if name == "" {
			name = fmt.Sprintf("device %s", mac)
		}
		log.Printf("relayed WoL for %s to %s:%d (requested by %s)", name, device.IP, cfg.Port, source)
	}
}
