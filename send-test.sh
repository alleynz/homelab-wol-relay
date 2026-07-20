#!/bin/sh
set -eu

MAC="${1:-62:00:00:9c:72:6e}"
DEST="${2:-255.255.255.255}"
PORT="${3:-9}"

python3 - "$MAC" "$DEST" "$PORT" <<'PY'
import socket
import sys

mac = bytes.fromhex(sys.argv[1].replace(':', '').replace('-', ''))
if len(mac) != 6:
    raise SystemExit('MAC must contain 6 bytes')
packet = b'\xff' * 6 + mac * 16

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
sock.setsockopt(socket.SOL_SOCKET, socket.SO_BROADCAST, 1)
sock.sendto(packet, (sys.argv[2], int(sys.argv[3])))
print(f'Sent WoL for {sys.argv[1]} to {sys.argv[2]}:{sys.argv[3]}')
PY
