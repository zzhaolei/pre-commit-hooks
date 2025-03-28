# pre-commit-hooks

Some out-of-the-box hooks for pre-commit.

See also: https://github.com/pre-commit/pre-commit

### Using pre-commit-hooks with pre-commit

Add this to your `.pre-commit-config.yaml`

```yaml
- repo: https://github.com/zzhaolei/pre-commit-hooks
  rev: v1.0.0
  hooks:
    - id: check-added-large-files
  # -   id: ...
```

### Hooks available

#### `check-added-large-files`

Prevent giant files from being committed. Like `jujutsu`.
