# ddns
Cloudflare DDNS client

## Usage:

In ```config.toml```:
```toml
email = "invalid@example.test"
apiKey = "abcd"
zoneID = "efgh"
records = ["example.test", "sub.example.test"]
```

```bash
sudo snap set ddns config="$(cat config.toml)"
```