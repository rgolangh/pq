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
```


