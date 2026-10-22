# Architecture

This document maps the repository to the manuscript.

## Fabric network (Section 3.1, Table 2)
- 2 organisations: `HospitalMSP`, `LabMSP`
- 1 Raft orderer, 1 Fabric CA per org
- Application channel `healthcare`, endorsement policy MAJORITY
- CouchDB state database
- TLS 1.3 on every peer/orderer/client channel

## Hybrid storage (Section 3.2, Algorithms 2 & 5)
- Textual PHI → `chaincode/accesscontrol`, encrypted with AES-256-GCM
- Multimedia → `client/ipfs` chunks each encrypted with a per-chunk HKDF key,
  uploaded to a 10-node IPFS swarm, CIDs written on-chain

## Encryption (Section 3.4, Algorithm 5)
- AES-256-GCM for records and chunks
- RSA-2048 OAEP key wrapping (per recipient)
- HKDF-SHA256 key derivation
- Key epochs for revocation (`chaincode/crypto/keys.go`)

## Benchmark (Section 4.3, Table 6)
- `benchmark/run_all.sh` executes 30 independent runs
- network restart + 120 s warm-up per run
- `benchmark/stats.py` computes mean ± 95% CI, CV, Mann–Whitney U