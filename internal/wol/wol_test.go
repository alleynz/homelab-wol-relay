package wol

import (
	"net"
	"testing"
)

func TestBuildAndParseMagicPacket(t *testing.T) {
	mac, err := net.ParseMAC("62:00:00:9c:72:6e")
	if err != nil {
		t.Fatal(err)
	}

	packet, err := BuildMagicPacket(mac)
	if err != nil {
		t.Fatal(err)
	}
	if len(packet) != MagicPacketLength {
		t.Fatalf("expected %d bytes, got %d", MagicPacketLength, len(packet))
	}

	parsed, err := ParseMagicPacket(packet)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.String() != mac.String() {
		t.Fatalf("expected %s, got %s", mac, parsed)
	}
}

func TestParseRejectsInvalidPacket(t *testing.T) {
	if _, err := ParseMagicPacket(make([]byte, 102)); err == nil {
		t.Fatal("expected invalid packet to fail")
	}
}
