# ddns
Cloudflare DDNS client

## Usage:

In ```ddns.toml```:
```toml
apiToken = "abcd"
zoneID = "efgh"
recordNames = ["example.test", "sub.example.test"]
```

```console
> sudo snap set ddns config="$(cat ddns.toml)"
```
