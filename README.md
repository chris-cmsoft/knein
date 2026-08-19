# knein

Terminal CLI to search for contexts and then open k9s with the selected context

> [!NOTE]  
> This was fully vibe coded to make my life easier working with many contexts. 

## Installation 

```
go install github.com/chris-cmsoft/knein@latest
```

## Usage

```
$ knein
```

### Flags

| Flag | Default | Purpose |
| --- | --- | --- |
| `--kubeconfig` | standard resolution (`$KUBECONFIG`, then `~/.kube/config`) | Path to the kubeconfig to read contexts from. |
| `--limit` | `9` | Maximum contexts the picker shows at once. |

## Architecture

Context discovery, filtering and the picker live in
[kubecontext](https://github.com/chris-cmsoft/gotool-kubecontext-picker), shared with `kctx`.
It splits your input by spaces, and then finds contexts which match all those
texts.
