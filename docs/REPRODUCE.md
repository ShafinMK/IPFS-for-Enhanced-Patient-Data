# Reproduction guide

## 1. Hardware
Each IPFS node is intended to run on a separate physical host with the
configuration in Section 4.1 (i5-12600K, 64 GB RAM, NVMe SSD, 100 Mbps
wireless). The Fabric peers and orderer can share one host; IPFS nodes must
be distributed across two ISP-connected networks for the conservative
topology reported in the paper.

## 2. Two-ISP topology
- `ipfs0`–`ipfs4` behind ISP-A
- `ipfs5`–`ipfs9` behind ISP-B
Both networks must be reachable from the Fabric hosts; the Fabric hosts are
placed on a third segment and connect to the IPFS peers via their swarm
multiaddresses.

## 3. Steps
```bash
make network-up
make channel-create
make chaincode-deploy
make seed-data
make benchmark
make stats