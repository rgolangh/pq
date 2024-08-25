# pq is a cli to help manage podman quadlets

NOTE: UNDER DEVELOPMENT

See how quadlets are stored in a git repository https://github.com/rgolangh/podman-quadlets.

This git repo is used by default. Override with `--repo https://my/git/repo`

## Usage

```bash
$ pq list
Listing quadlets from repo https://github.com/rgolangh/podman-quadlets (default in ~/.config/pq/pq.yaml)

- kind
- wordpress
- nginx

$ pq install kind
Installing quadlet 'kind'
[#############             ]

$ pq install kind --systemd-reload
Installing quadlet 'kind'
[#############             ]

$ pq install kind --repo https://github.com/rgolangh/podman-quadlets
Installing quadlet 'kind' from https://github.com/rgolangh/podman-quadlets
[#############             ]

$ pq list --installed
- kind (on 24/01/2024)

$ pq remove kind
```

