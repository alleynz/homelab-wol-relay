# Homelab WoL Relay

A lightweight Wake-on-LAN relay written in Go.

The relay listens for broadcast Wake-on-LAN (WoL) magic packets, extracts the target MAC address, and retransmits the packet as a **unicast** WoL packet to a configured IP address.

This solves a common problem where applications (such as Apache Guacamole) only send broadcast WoL packets, which cannot cross routed networks or VLANs.

---

## Features

- Lightweight Go application
- Listens for broadcast WoL packets
- Retransmits WoL packets as unicast
- JSON configuration
- Docker support
- Multi-platform container image (planned)

---

## Project Status

This project is currently under active development.

Current goals:

- Broadcast WoL listener
- MAC address extraction
- Static MAC → IP mapping
- Unicast WoL transmission
- Docker image publishing

---

## Repository

```
https://github.com/alleynz/homelab-wol-relay
```

---

## Configuration

Create a directory for your configuration.

For example:

```bash
mkdir -p ./config
```

or

```bash
mkdir -p ~/docker/wol-relay/config
```

or

```bash
sudo mkdir -p /srv/wol-relay/config
```

Copy the example configuration:

```bash
cp config/config.example.json ./config/config.json
```

The application looks for:

```
/config/config.json
```

inside the container.

---

## Docker

An example Docker Compose file is provided in:

```
docker/docker-compose.yml
```

Adjust the volume path on the left-hand side to wherever you store your configuration.

Example:

```yaml
services:
  wol-relay:
    image: ghcr.io/alleynz/homelab-wol-relay:latest
    container_name: wol-relay
    restart: unless-stopped

    volumes:
      - ./config:/config:ro
```

The relay does not require any particular Docker network.

If you need to receive WoL packets from another container (for example Apache Guacamole), simply deploy both containers onto the same Docker network.

---

## Building

Clone the repository:

```bash
git clone https://github.com/alleynz/homelab-wol-relay.git
cd homelab-wol-relay
```

Build:

```bash
go build ./cmd/wol-relay
```

Run:

```bash
./wol-relay --config ./config/config.json
```

---

## Docker Build

```bash
docker build -t homelab-wol-relay .
```

Run:

```bash
docker run \
  -v ./config:/config:ro \
  ghcr.io/alleynz/homelab-wol-relay:latest
```

---

## Roadmap

- [ ] Broadcast WoL listener
- [ ] MAC address extraction
- [ ] Static MAC → IP mapping
- [ ] Unicast WoL transmission
- [ ] Logging
- [ ] Configuration validation
- [ ] Unit tests
- [ ] GitHub Actions
- [ ] Multi-architecture Docker images

---

## Contributing

Issues and pull requests are welcome.

---

## License

MIT License
