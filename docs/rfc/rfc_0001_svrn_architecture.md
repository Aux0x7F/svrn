# RFC‑0001: svrn — Sovereins Core Architecture

**Status:** Draft (not yet ratified)
**Last updated:** 2025‑11‑06
**Authors:** 0x7F
**Applies to:** `svrn` node, protocol, data model, and role system
**Supersedes:** None (first RFC)
**Related docs:** `ARCHITECTURE.md`, RFC‑0002 (Blob Service), RFC‑0003 (CRDT Ops)

---

## 0. Abstract

This document defines the core architecture of **svrn** — a decentralized, community‑scoped cloud platform built on I2P, providing encrypted object and collaborative data storage using a Kademlia‑based DHT, capability‑bound write authorization, and DID‑based cryptographic identity.

The goal of this RFC is to define the **minimum implementable specification** for:

* Node roles and behavior
* Transport and routing assumptions
* Identity, record, and dataset models
* Access control and revocation rules
* Expected invariants for providers, relays, and consumers

This RFC does **not** define:

* Full wire‑format of each service (covered in later RFCs)
* UI or application‑level conventions
* Federation outside I2P
* Non‑MVP optimizations such as PRE or DHT hardening

---

## 1. Motivation

Centralized cloud services place user data, authentication, and administrative power under the control of third parties. Existing self‑hosted alternatives require:

* Long‑term server stability and DNS exposure
* Complex TLS, NAT traversal, and firewall configuration
* Central administrators who control user access

**svrn** enables communities to operate a shared, encrypted data layer with:

* No clearnet exposure
* No single authoritative service
* No dependency on administrator identity management

Sovereins provides:

| Requirement | Solution |
|------------|----------|
| Privacy against network observers | I2P overlay, no direct IP routing |
| Decentralized service discovery | Kademlia DHT |
| Access control without accounts | Capability tokens (signed) |
| Revocable and rotating encryption | Dataset epochs and keybags |
| Data portability and exitability | Local consumer nodes can mirror datasets |

---

## 2. Definitions

| Term | Definition |
|------|------------|
| **Node** | An instance of the `svrn` binary running in any role |
| **Consumer** | Node that fetches/decrypts data but does not serve it |
| **Provider** | Node that stores and serves encrypted data objects |
| **Relay** | Node that participates in the DHT but does not store datasets |
| **Seed** | A relay that publishes a **signed seedlist** for bootstrap |
| **Dataset** | A logical container with its own key material and write rules |
| **Capability Token (CapToken)** | Signed permission object authorizing writes |
| **Keybag** | Encrypted key distribution bundle mapping dataset keys → recipients |
| **Epoch** | A dataset encryption generation; used for revocation and rotation |

---

## 3. Architectural Constraints

Each of the following MUST be true for an implementation to conform to RFC‑0001:

1. **All traffic MUST route through I2P**. No public IP listeners.
2. **All writes MUST be authenticated** using COSE‑signed envelopes.
3. **All stored content MUST be encrypted at rest** using per‑dataset symmetric keys.
4. **All metadata stored in the DHT MUST be signed** and MUST implement sequence‑based replay protection.
5. **Nodes MUST accept and ignore invalid writes** (fail‑closed behavior).
6. **Nodes MUST NOT require centralized identity or ACL servers**.
7. **Consumers MUST be able to sync and decrypt data without running a provider.**

---

## 4. Node Roles

### 4.1 Role Selection
A node declares roles via config or CLI flags. If no roles are defined, it operates as a **consumer**.

```
svrn --roles provider,relay --services blob,crdt --router auto
```

### 4.2 Role Matrix
| Node type | Stores data | Serves data | DHT participation | Publishes seedlist |
|-----------|-------------|-------------|-------------------|--------------------|
| Consumer  | optional local cache | ❌ | optional | ❌ |
| Provider  | ✅ | ✅ | ✅ | ❌ |
| Relay     | ❌ | ❌ | ✅ | ❌ |
| Seed      | ❌ | ❌ | ✅ | ✅ |

A single binary MAY assume multiple roles.

---

## 5. Transport

### 5.1 I2P Requirement
All node‑to‑node communication MUST occur over I2P. Implementations MAY:
* Start an embedded router (i2p‑zero model)
* Connect to external SAM bridge (`external:host:port`)

### 5.2 Local Bind Rules
* All services MUST bind only to 127.0.0.1 locally
* I2P server tunnels expose externally routable endpoints

---

## 6. Identity Model

* Every node MUST generate or import an Ed25519 DID (`did:key:` form)
* Identity MUST be used for:
  * Signing NodeRecords
  * Signing DatasetRecords
  * Signing Capability Tokens
  * Signing Write Envelopes

DID rotation is permitted; linkage is not required.

---

## 7. Data Model Summary

### 7.1 NodeRecord (for DHT discovery)
* MUST be COSE_Sign1
* MUST include at least one `dht` endpoint
* MUST include `seq` and `exp` fields

### 7.2 DatasetRecord
* MUST include dataset DID, epoch, and encryption algorithm
* MUST NOT leak KEK or DEK material

### 7.3 Capability Token
* MUST be signed by dataset owner or delegated issuer
* MUST express scope, caveats, and expiry

### 7.4 WriteEnvelope
* MUST include dataset ID, sequence, and operation type
* MUST be signed by writer DID
* MUST be validated by providers before storage

(More formal CBOR/COSE struct definitions appear in later RFCs.)

---

## 8. DHT Behavior Requirements

* Nodes MUST implement Kademlia routing semantics (XOR distance)
* Providers MUST service STORE and FIND requests for metadata
* Providers MUST NOT store encrypted dataset objects in DHT
* Seeds MUST serve a public `/bootstrap` endpoint over I2P containing a signed seedlist

---

## 9. Data Storage Requirements

| Object Type | Stored by | Encryption | Replication |
|-------------|-----------|------------|-------------|
| NodeRecord | Relay, Provider, Seed | Signed only | default Kademlia |
| DatasetRecord | Relay, Provider | Signed only | default Kademlia |
| Blob object | Provider only | AEAD encrypted | provider policy |
| CRDT ops | Provider only | AEAD encrypted | provider policy |

---

## 10. Revocation and Epoch Handling

1. Dataset owner increments `epoch` in a new DatasetRecord
2. Providers MUST treat old‑epoch writes as invalid
3. Consumers MAY continue decrypting old data while keys remain available
4. Rotation of recipient keybags is out of scope for this RFC but MUST NOT require provider trust

---

## 11. Security Properties

| Threat | Mitigation |
|--------|------------|
| Metadata spoofing | COSE signatures + seq + exp |
| Unauthorized writes | Capability token verification |
| Passive traffic observers | I2P transport + AEAD objects |
| Replay attacks | `seq`, `exp`, optional `nbf` |
| Malicious providers | Client‑side verification of every envelope |

---

## 12. Non‑Goals (for RFC‑0001)

* Global, cross‑community routing
* Non‑I2P transports (Tor, WAN, LAN broadcast)
* Reputation‑weighted DHT routing
* Multi‑writer key escrow and secure deletion
* Multi‑dataset atomic transactions

These MAY be defined in future RFCs.

---

## 13. Compliance Checklist

A node is compliant with RFC‑0001 if:

- [ ] It never opens a clearnet listener
- [ ] It signs all published metadata
- [ ] It rejects unsigned or expired writes
- [ ] It stores no plaintext dataset content
- [ ] It fails closed (invalid operations are ignored, not processed)
- [ ] It can operate as a consumer without administrator approval

---

## 14. Open Questions

* Should capability tokens support on‑chain delegation breadcrumbs?
* Should dataset epochs mandate automatic re‑encryption sweeps?
* Should consumers be permitted to run partial DHT caches?
* Should we adopt MLS or HPKE in a future revision?

---

## 15. Change Log

| Version | Date | Changes |
|----------|------|----------|
| 0.1‑draft | 2025‑11‑06 | Initial spec skeleton |
| 0.2‑draft | TBD | Add formal CBOR schema + service error codes |

---

## 16. Approval

**This RFC becomes active when:**
* At least 2 maintainers sign off
* Reference implementation passes the RFC‑0001 test suite

Once accepted, changes require a new RFC or a version bump.

---

_End of RFC‑0001_

