# svrn — Sovereins

Reins without rulers.

A decentralized, community-scoped cloud platform built on I2P.  
Privacy-first, capability-based access, no central servers, no clearnet.

```
┌───────────────┐     ┌────────────────┐
│  I2P Overlay  │◀───▶│  svrn nodes    │
└───────────────┘     │ (consumer /    │
                       │  provider /    │
                       │  relay / seed) │
                       └────────────────┘
```

---

## ✨ Project Status
Status: **Pre‑alpha / RFC drafting phase**  
Implemented code: **not yet started**  
Specs drafted so far:

| RFC | Scope |
|------|--------|
| RFC‑0001 | Core architecture, identities, DHT, security |
| RFC‑0002 | Blob (LWW) encrypted object storage |
| RFC‑0003 | CRDT encrypted append‑only ops + sync API |
| RFC‑0004 | Keybags, epoch rotation, revocation model |

More RFCs coming: seedlists, DHT wire protocol, capability delegation, installer/packaging.

---

## 🔧 What is `svrn`?
`svrn` is a single self‑contained binary that can run in one or more roles:

| Role | Description |
|-------|-------------|
| **consumer** | Fetch + decrypt data, no inbound ports |
| **provider** | Stores + serves encrypted blobs/ops |
| **relay** | Participates in DHT routing, stores metadata only |
| **seed** | Bootstrap relay, publishes signed seedlist |

All traffic goes through I2P tunnels. No node ever needs a public IP.

---

## 🧠 Core Principles
✅ No star topology — all discovery via DHT  
✅ Providers cannot decrypt data (E2E encryption)  
✅ Write access = signed capability, not user account  
✅ Revocation = key rotation, not ACL changes  
✅ Same binary for all nodes, roles chosen at runtime

---

## 📦 Repo Layout
```
svrn/
├─ cmd/svrn/                # main entrypoint
├─ internal/                # not importable outside module
│  ├─ config/               # YAML/env/flag loader
│  ├─ i2p/                  # router + tunnel manager
│  ├─ http/                 # local service mux
│  └─ logging/              # zap/slog wrapper
├─ pkg/                     # public internal packages
│  ├─ agent/                # lifecycle + roles
│  ├─ dht/                  # Kademlia implementation
│  ├─ record/               # COSE/CBOR record types
│  ├─ crypto/               # ed25519/x25519/xchacha20
│  └─ services/             # blob + crdt impls
├─ docs/
│  ├─ ARCHITECTURE.md
│  └─ rfc/
│     ├─ RFC-0001-architecture.md
│     ├─ RFC-0002-blob-service.md
│     ├─ RFC-0003-crdt-ops.md
│     └─ RFC-0004-keybags-epoch-rotation.md
└─ deploy/                  # installers, systemd, docker, iso, etc
```

---

## 🚀 Roadmap (high‑level)

| Phase | Milestone |
|--------|-----------|
| ✅ 0 | RFCs drafted (ongoing) |
| ⏳ 1 | Bootstrap code scaffolding (config + agent + logging) |
| ⏳ 2 | Implement I2P router adapter + local tunnel model |
| ⏳ 3 | Implement DHT store/find RPC and node record signing |
| ⏳ 4 | Implement Blob service (RFC‑0002) |
| ⏳ 5 | Implement CRDT service (RFC‑0003) |
| ⏳ 6 | Implement keybag + rotation logic (RFC‑0004) |
| ⏳ 7 | MVP cluster demo + local seed bootstrap |
| ⏳ 8 | Packaging: Windows, Linux, Docker, ISO installer |

---

## 🔨 Build (future)
```sh
go build -o build/svrn ./cmd/svrn
```

---

## 🗣 Contributing
RFC review and discussion welcome.  
Code contributions will open after scaffold phase begins.

To participate in RFC review, open an issue titled:
```
RFC‑000X: <topic>
```

---

## 📜 License
MIT
---

## 🌐 Project Links
* Website: TBD
* Docs: `docs/`
* RFCs: `docs/rfc/`
* Discussion: GH Issues until bootstrap phase

---

EOF

