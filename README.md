# nox

## Getting Started

### Prerequisites

- Go 1.18+

### Installation

```bash
go install github.com/aottr/nox/cmd/nox@latest
```

## Usage

### Generate age key

```bash
mkdir -p keys secrets
```

```bash
age-keygen -o keys/key.txt
```

### Encrypt secrets

```bash
age -r <recipient> -o secrets/prod.env.age secrets/prod.env
```

### Configure

Create a `config.yaml` file with the following contents:

```yaml
interval: "10m"
age: 
  identity: "keys/key.txt"
  recipients: # optional, when used to encrypt secrets
    - "age1xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
statePath: ".nox-state.json"
defaultRepo: git@github.com:ShorkBytes/nox-secrets.git

apps:
  debug:
    branch: main
    files:
      - path: debug/debug.age
        output: ./secrets/.env
```

### Run

```bash
nox --help
```

### How to

#### Decrypt secret into custom file

```bash
nox decrypt --app debug --dry-run > secrets.env
```

### Contributing

Contributions are welcome!

```bash
go fmt ./...
```

## License

[MIT](LICENSE)