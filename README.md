# knein

Terminal CLI to search for contexts and then open k9s with the selected context

> [!NOTE]  
> This was fully vibe coded to make my life easier working with many contexts. 

## Installation

### Download a binary

Every release attaches a binary per platform, so no Go toolchain is needed:

```console
curl -fsSL -o knein https://github.com/chris-cmsoft/knein/releases/latest/download/knein_darwin_arm64
chmod +x knein
mv knein /usr/local/bin/knein
```

Swap `darwin_arm64` for `darwin_amd64`, `linux_amd64` or `linux_arm64`. Each
release also ships `checksums.txt`:

```console
sha256sum -c checksums.txt --ignore-missing
```

If you download through a browser instead of `curl`, macOS quarantines the file
and refuses to run it. Clear the flag with
`xattr -d com.apple.quarantine knein`.

### With Go

```console
go install github.com/chris-cmsoft/knein@latest
```

Or build locally with `make build`, or cross compile every released target into
`dist/` with `make dist`.

## Usage

```
$ knein
```

### Flags

| Flag | Default | Purpose |
| --- | --- | --- |
| `--kubeconfig` | standard resolution (`$KUBECONFIG`, then `~/.kube/config`) | Path to the kubeconfig to read contexts from. |
| `--limit` | `9` | Maximum contexts the picker shows at once. |
| `--version` | | Print the running build and the newest release. |

## Version

`--version` prints the running build and asks GitHub for the newest release:

```console
$ knein --version
knein   v0.1.1
latest  v0.1.2 (https://github.com/chris-cmsoft/knein/releases/latest)
```

When the two match, the second line reads `up to date` instead. The lookup is
best effort: with no network it reports `latest  unknown (...)` and still prints
the running version.

Release builds are stamped with their tag. Builds from source report what
`git describe` says, and `go install` builds report the module version Go
records in the binary.

## Architecture

Context discovery, filtering and the picker live in
[kubecontext](https://github.com/chris-cmsoft/gotool-kubecontext-picker), shared with `kctx`.
It splits your input by spaces, and then finds contexts which match all those
texts.
