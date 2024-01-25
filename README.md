# pq is a cli to help manage podman quadlets

```
$ pq list
Listing quadlets from repo https://rgolangh/podman-quadlets (default in ~/.config/pq/pq.yaml)

- kind
- wordpress
- nginx

$ pq install kind
Installing quadlet 'kind'
[#############             ]
Done

$ pq install kind --systemd-reload
Installing quadlet 'kind'
[#############             ]
Done

$ pq install kind https://github.com/rgolangh/podman-quadlets
Installing quadlet 'kind' from https://github.com/rgolangh/podman-quadlets
[#############             ]
Done

$ pq list installed
- kind (on 24/01/2024)

```

