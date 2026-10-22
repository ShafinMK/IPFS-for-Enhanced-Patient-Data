# Patient-Centric Healthcare Data Management on Hyperledger Fabric + IPFS

Reference implementation and deployment artifacts for the paper:

> **A Comprehensive Framework Using Hyperledger-Based Private Blockchain and IPFS for Enhanced Patient Data Security and Accessibility**

This repository contains the complete source code, deployment configuration, benchmark harness, and statistical analysis scripts needed to reproduce the results reported in the manuscript (Sections 3–4, Tables 2 and 6, and Algorithms 1–5).

## Contents

| Path | Description | Manuscript Section |
|---|---|---|
| `chaincode/accesscontrol` | CRUD smart contract for patient records | §3.3, Algorithm 1 |
| `chaincode/datasharing` | Consent-gated data-sharing contract | §3.3, Algorithm 4 |
| `chaincode/crypto` | AES-256-GCM, RSA-OAEP, HKDF, and key-epoch logic | §3.4, Algorithm 5 |
| `client/ipfs` | Chunking, per-chunk encryption, CID upload and retrieval | §3.2, Algorithm 2 |
| `network/` and `config/` | Fabric v2.2.14 network, Raft orderer, and Fabric CA | §3.1, §3.6, Table 2 |
| `deploy/` (Docker Compose) | Containerized deployment | §3.6 |
| `benchmark/` | Load driver, 30-run repetition, confidence intervals, and Mann–Whitney testing | §4.3, Table 6 |
| `data/synthetic/` | Synthetic, non-PHI record generator | §4.1 |

> **Privacy notice:** No real patient data or protected health information (PHI) are included. All records used in the evaluation are synthetic and can be generated using the tools in `data/synthetic/`.

## Prerequisites

Ensure the following tools are installed before deployment:

- Docker 24 or later
- Docker Compose v2
- Go 1.21 or later
- Hyperledger Fabric binaries v2.2.14:
  - `peer`
  - `orderer`
  - `configtxgen`
  - `cryptogen`
- IPFS (Kubo) v0.29 or later
- Python 3.10 or later
- Python packages:
  - `scipy`
  - `numpy`
  - `pandas`

## Quick Start

Run the following commands from the repository root.

### 1. Start the Network

Generate cryptographic materials and start the Hyperledger Fabric and IPFS services.

```bash
make network-up
```

### 2. Create the Application Channel

Create the application channel and join the participating peers.

```bash
make channel-create
```

### 3. Deploy the Chaincode

Package, install, approve, and commit the three smart contracts.

```bash
make chaincode-deploy
```

### 4. Generate and Seed Synthetic Data

Generate synthetic patient records, encrypt the data chunks, and upload them to IPFS.

```bash
make seed-data
```

### 5. Run the Benchmark Suite

Execute the complete benchmark suite using 30 repeated runs.

```bash
make benchmark
```

### 6. Generate Statistical Results

Calculate the mean, 95% confidence intervals (CI), coefficient of variation (CV), and Mann–Whitney test results for Table 6.

```bash
make stats
```

## Reproducibility

The repository is organized to support reproduction of the implementation and evaluation described in the manuscript.

- **Blockchain layer:** Hyperledger Fabric network and smart contracts.
- **Storage layer:** IPFS-based storage with encrypted data chunks.
- **Cryptographic layer:** AES-256-GCM, RSA-OAEP, HKDF, and key-epoch management.
- **Benchmarking layer:** Workload generation and repeated performance measurements.
- **Statistical analysis:** Confidence intervals, coefficient of variation, and Mann–Whitney testing.

For reproducible results, use the specified software versions and retain the benchmark configuration, generated datasets, and raw measurements alongside the statistical outputs.

## License

This project is distributed under the Apache License 2.0. See the [`LICENSE`](LICENSE) file for the full license text.

## Citation

If you use this implementation in your research, please cite the associated paper using the metadata provided in [`CITATION.cff`](CITATION.cff).

---

**Disclaimer:** This repository provides a research reference implementation. Deployment in a production healthcare environment requires additional security review, privacy controls, access-policy validation, operational monitoring, and compliance assessment.
