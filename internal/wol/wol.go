package wol

import (
	"bytes"
	"fmt"
	"net"
)

const MagicPacketLength = 102

func ParseMagicPacket(packet []byte) (net.HardwareAddr, error) {
	if len(packet) < MagicPacketLength {
		return nil, fmt.Errorf("packet too short: got %d bytes", len(packet))
	}

	start := -1
	for i := 0; i <= len(packet)-MagicPacketLength; i++ {
		if bytes.Equal(packet[i:i+6], []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}) {
			start = i
			break
		}
	}
	if start == -1 {
		return nil, fmt.Errorf("magic packet header not found")
	}

	mac := append(net.HardwareAddr(nil), packet[start+6:start+12]...)
	for i := 1; i < 16; i++ {
		offset := start + 6 + (i * 6)
		if !bytes.Equal(packet[offset:offset+6], mac) {
			return nil, fmt.Errorf("MAC repetition %d does not match", i+1)
		}
	}

	return mac, nil
}

func BuildMagicPacket(mac net.HardwareAddr) ([]byte, error) {
	if len(mac) != 6 {
		return nil, fmt.Errorf("MAC must contain 6 bytes")
	}

	packet := make([]byte, 0, MagicPacketLength)
	packet = append(packet, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff)
	for i := 0; i < 16; i++ {
		packet = append(packet, mac...)
	}
	return packet, nil
}

func SendUnicast(ip string, port int, packet []byte) error {
	addr := &net.UDPAddr{IP: net.ParseIP(ip), Port: port}
	if addr.IP == nil {
		return fmt.Errorf("invalid destination IP %q", ip)
	}

	conn, err := net.DialUDP("udp4", nil, addr)
	if err != nil {
		return fmt.Errorf("open UDP connection: %w", err)
	}
	defer conn.Close()

	if _, err := conn.Write(packet); err != nil {
		return fmt.Errorf("send UDP packet: %w", err)
	}
	return nil
}
