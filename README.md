# pq is a cli to help manage podman quadlets

> [!NOTE]
> Uner Development

See how quadlets are stored in a git repository https://github.com/rgolangh/podman-quadlets.

This git repo is used by default. Override with `--repo https://my/git/repo`

## Usage

```console
$ pq list
Listing quadlets from repo https://github.com/rgolangh/podman-quadlets (default in ~/.config/pq/pq.yaml)

- nginx
- redpanda
- wordpress

$ pq install wordpress
Installing quadlet "wordpress"
[#############             ]
Reload systemd daemon?[y/N]y
Reloading systemd daemon for the current user
Starting service wordpress.service for current user
Starting service wordpress-db.service for current user

$ pq install wordpress --repo https://github.com/rgolangh/podman-quadlets
Installing quadlet "wordpress" from https://github.com/rgolangh/podman-quadlets
[#############             ]
Reload systemd daemon?[y/N]y
Reloading systemd daemon for the current user
Starting service wordpress.service for current user
Starting service wordpress-db.service for current user

$ pq list --installed
- wordpress (on 24/01/2024)

$ pq remove wordpress
Stopping service wordpress-db.service for current user
Stopping service wordpress.service for current user
Remove quadlet "wordpress" from path /var/home/rgolan/.config/containers/systemd/wordpress?[y/n]y
removed "wordpress" from path /var/home/rgolan/.config/containers/systemd/wordpress
Reload systemd daemon?[y/N]y
Reloading systemd daemon for the current user

$ pq list-services
nextcloud - nextcloud-aio-master.service active (running)
redpanda - console.service inactive (dead)
redpanda - redpanda.service inactive (dead)

$ pq inspect nextcloud
Inspect quadlet "nextcloud"
# Source: https://github.com/rgolangh/podman-quadlets nextcloud/nextcloud-aio-master.container
[Unit]
Description=Nextcloud AIO Master Container
Documentation=https://github.com/nextcloud/all-in-one/blob/main/docker-rootless.md
After=local-fs.target
Requires=podman.socket

[Service]
TimeoutStartSec=900

[Container]
ContainerName=nextcloud-aio-mastercontainer
Image=docker.io/nextcloud/all-in-one:latest
AutoUpdate=registry
PublishPort=127.0.0.1:11001:8080
Volume=nextcloud_aio_mastercontainer:/mnt/docker-aio-config
Volume=/run/user/%U/podman/podman.sock:/var/run/docker.sock:ro
Network=bridge
SecurityLabelDisable=true

Environment=APACHE_PORT=11000
Environment=APACHE_IP_BINDING=127.0.0.1
Environment=WATCHTOWER_DOCKER_SOCKET_PATH=/run/user/%U/podman/podman.sock

[Install]
WantedBy=multi-user.target default.target

# Source: https://github.com/rgolangh/podman-quadlets nextcloud/nextcloud-aio-master.volume
[Volume]
VolumeName=nextcloud_aio_mastercontainer

```


