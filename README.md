# pq is a cli to help manage podman quadlets

NOTE: UNDER DEVELOPMENT

See how quadlets are stored in a git repository https://github.com/rgolangh/podman-quadlets.

This git repo is used by default. Override with `--repo https://my/git/repo`

## Usage

```bash
$ pq list
Listing quadlets from repo https://github.com/rgolangh/podman-quadlets (default in ~/.config/pq/pq.yaml)

- nginx
- redpanda
- wordpress

$ pq install wordpress
Installing quadlet 'wordpress'
[#############             ]

$ pq install wordpress --no-systemd-daemon-reload
Installing quadlet 'wordpress'
[#############             ]

$ pq install wordpress --repo https://github.com/rgolangh/podman-quadlets
Installing quadlet 'wordpress' from https://github.com/rgolangh/podman-quadlets
[#############             ]

$ pq list --installed
- wordpress (on 24/01/2024)

$ pq remove wordpress
```


