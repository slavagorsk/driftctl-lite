# driftctl-lite

Lightweight CLI tool to detect infrastructure drift between Terraform state and actual cloud resources.

---

## Installation

```bash
go install github.com/yourusername/driftctl-lite@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftctl-lite.git
cd driftctl-lite
go build -o driftctl-lite .
```

---

## Usage

Point `driftctl-lite` at your Terraform state file and let it compare against your live cloud environment:

```bash
# Scan using a local state file
driftctl-lite scan --state terraform.tfstate --provider aws

# Scan using a remote S3 backend
driftctl-lite scan --state s3://my-bucket/terraform.tfstate --provider aws --region us-east-1
```

### Example Output

```
[✔] aws_s3_bucket.my-bucket        — in sync
[✗] aws_security_group.web         — DRIFTED (ingress rules modified)
[!] aws_iam_role.lambda_exec       — missing in state

Drift detected: 2 resource(s) out of sync
```

### Flags

| Flag | Description |
|------|-------------|
| `--state` | Path or URI to Terraform state file |
| `--provider` | Cloud provider (`aws`, `gcp`, `azure`) |
| `--region` | Cloud region (default: `us-east-1`) |
| `--output` | Output format: `text`, `json` (default: `text`) |

---

## Requirements

- Go 1.21+
- Valid cloud credentials configured (e.g., AWS credentials via `~/.aws/credentials` or environment variables)

---

## Contributing

Pull requests are welcome. Please open an issue first to discuss any significant changes.

---

## License

[MIT](LICENSE)