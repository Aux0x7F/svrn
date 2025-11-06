# svrn — Sovereins Architecture

**Repository:** https://github.com/sovereins/svrn  
**Binary:** `svrn`  
**Subtext:** *Sovereins — a decentralized community cloud*  
**Last updated:** 2025-11-06  
**Status:** Draft

---

## 1. Overview

Sovereins is a decentralized, community-scoped “cloud” built over I2P.  
It provides two data surfaces — *LWW Blobs* and *CRDT Ops* — for versioned and collaborative data sharing, respectively.  
Each node runs the same universal binary (`svrn`) and assumes one or more roles (`consumer`, `provider`, `relay`, `seed`) configured declaratively.

Key principles:

- **No star topology** — all peers route via a distributed DHT.
- **Privacy first** — all transport via I2P, no clearnet ingress.
- **Capabilities, not accounts** — possession of keys and signed tokens define access.
- **Composability** — datasets and apps live inside signed *communities*.
- **Single binary** — same build runs everything; roles chosen at runtime.

---

## 2. Core Architecture Diagram

```
         ┌──────────────────────────────┐
         │          I2P Overlay         │
         │   (Tunnels, eepsites, SAM)   │
         └──────────────┬───────────────┘
                        │
              ┌─────────┴─────────┐
              │     DHT Layer     │
              │  (Kademlia RPCs)  │
              └─────────┬─────────┘
                        │
      ┌─────────────────┼──────────────────┐
      │                 │                  │
  ┌───▼───┐         ┌───▼───┐         ┌───▼───┐
  │ Blobs │         │ CRDTs │         │  Seed │
  │ LWW   │         │  Ops  │         │Lists  │
  └───┬───┘         └───┬───┘         └───┬───┘
      │                 │                 │
      │                 │                 │
      ▼                 ▼                 ▼
  Encrypted Data    Encrypted Ops     Bootstrap Peers
```

---

## 3. Roles

| Role | Purpose | Default | Notes |
|------|----------|----------|-------|
| **consumer** | Lookup, decrypt, sync | ✅ | No inbound ports |
| **provider** | Serve `/blob/*`, `/crdt/*` | opt-in | Can also relay |
| **relay** | DHT store/lookup RPCs | opt-in | Metadata only |
| **seed** | Relay + publishes signed seedlist | rare | No data services |

All roles share the same codepath, toggled by configuration or CLI flags.

---

## 4. Key Subsystems

| Subsystem | Function |
|------------|-----------|
| **Agent** | Loads config, spawns subsystems by role |
| **I2P Adapter** | Starts or attaches to router, manages tunnels |
| **DHT** | XOR-metric Kademlia network, peer routing, STORE/FIND RPCs |
| **Records** | COSE/CBOR structures for NodeRecord, DatasetRecord, capabilities |
| **Blob Service** | Versioned file storage, LWW conflict resolution |
| **CRDT Service** | Append-only operation logs for structured data |
| **Crypto** | Ed25519, X25519, XChaCha20-Poly1305, HKDF |
| **Storage** | On-disk persistence of peers, blobs, ops |
| **Logging** | Structured, redaction-safe logs |

---

## 5. Data Objects

### NodeRecord
COSE_Sign1(CBOR)
```cbor
{ v:1, node:"did:key:...", tunnels:[
  {n:"dht",  b32:"...", p:8083},
  {n:"blob", b32:"...", p:8081},
  {n:"crdt", b32:"...", p:8082}
], seq:int, exp:int }
```

### DatasetRecord
COSE_Sign1(CBOR)
```cbor
{ v:1, dataset:"did:key:...", epoch:int,
  alg:"xchacha20poly1305", keybag_index:"<pointer>",
  seq:int, exp:int }
```

### Capability Token
COSE_Sign1(CBOR)
```cbor
{ v:1, iss:"did:key:<owner>", aud:"did:key:<writer>",
  cap:"dataset:<DID>:write:/docs/**",
  caveats:{ exp:int, max_bytes:int?, rate:str? } }
```

---

## 6. Access Control

| Action | Auth Mechanism | Enforcement |
|--------|----------------|-------------|
| **Read** | Possession of keybag entry → unwrap KEKₑ → decrypt | client |
| **Write** | Signed capability + envelope | provider verifies |
| **Discover** | Public DHT | signed metadata only |

Revocation: increment dataset *epoch* and distribute new keybags.  
Providers ignore writes signed with expired or invalid capabilities.

---

## 7. Encryption Summary

| Layer | Algorithm | Purpose |
|--------|------------|----------|
| DID / sigs | Ed25519 | identity & signing |
| Data encryption | XChaCha20-Poly1305 | AEAD payload |
| Key wrapping | X25519 | per-recipient KEKₑ wrap |
| Hash | SHA-256 | object and record IDs |
| Format | CBOR + COSE | portable structured encoding |

---

## 8. Networking and Transport

- **All** traffic routed through I2P (server tunnels for services, client tunnels for lookups).  
- Local HTTP listeners bind **127.0.0.1** only.  
- DHT messages are simple HTTP POSTs with COSE payloads.  
- Optional embedded router (`--router auto`) via i2p-zero; or attach external SAM host.

---

## 9. Directory Layout

```
svrn/
├─ cmd/svrn/                # main entrypoint
├─ internal/
│  ├─ config/               # load YAML/env/flags
│  ├─ i2p/                  # router/tunnel management
│  ├─ http/                 # local mux & framing
│  └─ logging/              # log wrapper
├─ pkg/
│  ├─ agent/                # lifecycle & roles
│  ├─ dht/                  # Kademlia implementation
│  ├─ record/               # COSE/CBOR types
│  ├─ crypto/               # crypto helpers
│  └─ services/
│     ├─ blob/              # versioned blobs
│     └─ crdt/              # collaborative ops
├─ docs/                    # RFCs, architecture
└─ deploy/                  # systemd, docker, installer assets
```

---

## 10. Config Example

`%PROGRAMDATA%\svrn\node.yaml`
```yaml
node:
  roles: [provider, relay]
  services:
    blob: { enabled: true }
    crdt: { enabled: true }
 i2p:
  router: auto
communities:
  - http://<seed>.b32.i2p/bootstrap/seedlist.json
storage:
  data_dir: C:\svrn\data
  max_store_mib: 20480
policy:
  rate_limit: 2M
logging:
  level: info
```

Command-line overrides:
```powershell
svrn.exe --roles provider,relay --services blob,crdt --router auto --community http://<b32>.b32.i2p/bootstrap/seedlist.json
```

---

## 11. Development Workflow

1. **Clone & bootstrap**
   ```powershell
   git clone https://github.com/sovereins/svrn.git
   cd svrn
   go mod tidy
   ```
2. **Build**
   ```powershell
   go build -o build\svrn.exe .\cmd\svrn
   ```
3. **Run (consumer mode)**
   ```powershell
   .\build\svrn.exe
   ```
4. **Run (provider mode)**
   ```powershell
   .\build\svrn.exe --roles provider,relay --services blob,crdt --router auto
   ```

---

## 12. Testing & CI

- **Tests:** `go test ./...`
- **Coverage:** `go test -cover`
- **CI:** GitHub Actions runs on Windows & Linux
