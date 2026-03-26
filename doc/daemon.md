# git-syn daemon

The `daemon` command runs the synchronization process periodically.

The daemon will push to all remotes in `.gitremotes` at the specified interval.

## Usage

```sh
git syn daemon [flags]
```

### Flags

- `--interval <duration>`: Sync interval (e.g., `5m`, `1h`). Defaults to the value in `~/.config/git-syn/config.yaml` or `5m`.
- `--path <dir>`: Path to the git repository (default: current directory).
- `--once`: Run exactly once and exit.

## How it works

The daemon uses a simple timer to trigger synchronization. It reads the `.gitremotes` file in the repository and attempts to push the current state of the repository to all remotes listed.

### Example

To run the daemon with a 10-minute interval:

```sh
git syn daemon --interval 10m
```

## Configuration

The default sync interval can be configured in the global configuration file:

```yaml
# ~/.config/git-syn/config.yaml
sync_interval: 5m
```
