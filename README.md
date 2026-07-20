# WoL Relay

A small Wake-on-LAN relay for Guacamole and segmented networks.

It listens for broadcast WoL packets, validates the magic packet, extracts the target MAC address, looks up a static MAC-to-IP mapping, and sends a new unicast WoL packet to the configured device.

## GitHub Container Registry

The included GitHub Actions workflow runs tests and publishes multi-architecture images for `linux/amd64` and `linux/arm64` to:

```text
ghcr.io/<github-user-or-org>/<repository>:latest
```

The image is rebuilt whenever code is pushed to `main`. Version tags such as `v0.1.0` also publish a matching container tag.

After the first successful workflow, open the package in GitHub and make it public if anonymous pulls are desired.

## Persistent configuration

Create a persistent directory on the Docker host:

```sh
sudo mkdir -p /opt/wol-relay/config
sudo nano /opt/wol-relay/config/config.json
```

Example:

```json
{
  "listen": "0.0.0.0",
  "port": 9,
  "devices": [
    {
      "name": "Windows PC",
      "mac": "62:00:00:9c:72:6e",
      "ip": "192.168.1.20"
    }
  ]
}
```

The configuration is stored outside the container and survives updates, restarts, and recreation.

## Deployment

Use `docker/docker-cosmos.yml`

```text
ghcr.io/YOUR_GITHUB_USERNAME/wol-relay:latest
```

with the image shown on the GitHub package page.

The relay must share the same docker network with Guacamole so it can receive Guacamole's UDP broadcast.

After editing the configuration, restart the container in Cosmos.

## Logs

```sh
docker logs -f wol-relay
```

Expected startup:

```text
WoL relay listening on 0.0.0.0:9 with 1 configured device(s)
```

Expected relay event:

```text
relayed WoL for Windows PC to 192.168.1.20:9
```

## Packet verification

```sh
sudo tcpdump -ni any 'udp port 9' -nn -vv -X
```

You should see Guacamole's broadcast followed by the relay's unicast packet.

## Local development

```sh
go test ./...
docker compose up -d --build
```

## Future web GUI

Configuration, packet parsing, and packet delivery are separate packages so an HTTP API and web GUI can be added without replacing the relay engine.
