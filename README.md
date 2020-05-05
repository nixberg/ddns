# ddns
Cloudflare DDNS client

## Usage

Initially, the config is empty and the daemon is stopped.

### Example config

```toml
apiToken = "<TOKEN>"
zoneID = "<ID>"
recordNames = ["example.test", "sub.example.test"]
```

### Starting the daemon

#### Bash

```bash
> sudo snap set ddns config="$(cat ddns.toml)"
```

#### Fish

```fish
> read -z config < ddns.toml
> sudo snap set ddns config=$config
```

### Stopping the daemon

```bash
> sudo snap stop ddns
```

Or:

```bash
> sudo snap set ddns config=""
```
