# ddns
Cloudflare DDNS client

## Usage:

In ```ddns.toml```:
```toml
email = "user@example.test"
apiKey = "abcd"
zoneID = "efgh"
records = ["example.test", "sub.example.test"]
```

```bash
sudo snap set ddns config="$(cat ddns.toml)"
```