# RFC-0003: Sovereins CRDT Ops Protocol

**Status:** Draft (not yet ratified)
**Last updated:** 2025-11-06
**Authors:** Sovereins Project Maintainers
**Applies to:** `svrn` Provider nodes, CRDT client sync API, Envelope validation
**Related RFCs:** RFC-0001 (Core Architecture), RFC-0002 (Blob Service)

---

## 0. Abstract

This RFC defines the **CRDT Ops Service Protocol**, the Sovereins subsystem for synchronized, append-only collaborative state. Unlike the Blob Service (RFC‑0002), CRDT Ops enables *causally ordered*, multi‑writer updates using encrypted operation batches.

The CRDT layer is used for boards, lists, shared JSON documents, wiki‑like structures, and any data that requires **convergent, conflict‑free, shared mutation**.

---

## 1. Goals

* Allow multiple writers to append ordered encrypted CRDT ops to a shared dataset.
* Guarantee deterministic client‑side convergence without provider trust.
* Preserve privacy: providers MUST NOT interpret, decrypt, or merge CRDT content.
* Support eventual synchronization via vector clock batching (`/since`, `/heads`).

---

## 2. Design Summary

A CRDT dataset consists of **encrypted op batches**, each wrapped in a signed **CrdtOpEnvelope**. Provider nodes:

* accept appends (`POST /crdt/append`),
* store them in local log order, and
* serve them to consumers via `/crdt/since` and `/crdt/heads`.

All merging happens **client‑side**, with the CRDT library chosen by the consumer (e.g. Yjs, Automerge).

---

## 3. Endpoints

### 3.1 POST /crdt/append

Append one or more encrypted CRDT ops.

#### Request Body
Binary framing (same pattern as RFC‑0002):
1. 4‑byte length prefix (CapToken)
2. COSE_Sign1(CapToken)
3. 4‑byte length prefix (CrdtOpEnvelope)
4. COSE_Sign1(CrdtOpEnvelope)
5. Ciphertext (opaque `ops[]` batch)

#### Response
```json
{
  "ok": true,
  "seq": 57
}
```

### 3.2 GET /crdt/since?dataset=<did>&clock=<base64>
Returns all op batches newer than the provided vector clock.

Response body (streamed NDJSON or CBOR array):
```
[{ envelope_b64, ciphertext_b64 }, ...]
```

### 3.3 GET /crdt/heads?dataset=<did>
Returns the provider’s current vector clock summary for the dataset.

```json
{
  "dataset": "did:key:…",
  "clock": { "node1": 12, "node2": 20, ... },
  "seq": 104
}
```

### 3.4 GET /crdt/snapshot?dataset=<did>
*Optional.* Providers MAY offer encrypted state snapshots for faster bootstrap.

---

## 4. CrdtOpEnvelope

COSE_Sign1(CBOR):
```cbor
{
  v: 1,
  dataset: "did:key:...",
  path: "/board/alpha",        ; logical doc pathway
  seq: 57,                      ; monotonic provider sequence
  clock: { ... },               ; vector clock (per writer DID)
  nbf: 1730920000,              ; optional not‑before
  exp: 1767225600               ; optional expiry
}
```

Properties:
* `seq` is assigned by writer; provider enforces monotonicity.
* `clock` allows causality tracking without decrypting ops.
* `path` MAY be used for sharding multiple docs inside dataset.

---

## 5. Validation Rules

| Check | Description | On failure |
|--------|-------------|-------------|
| CapToken signature | CapToken MUST validate | reject |
| Envelope signature | CrdtOpEnvelope MUST validate | reject |
| Capability scope | CapToken MUST permit append on `path` | reject |
| Sequence check | `seq` MUST be >= previous stored | reject |
| Epoch check | Envelope epoch MUST <= dataset epoch | reject |
| Clock format | MUST contain valid map of DID→int | reject |

Providers MUST NOT inspect decrypted ops.

---

## 6. Storage Model

```
<root>/crdt/<dataset_did>/<epoch>/<seq>.op
<root>/crdt/<dataset_did>/<epoch>/<seq>.envelope
<root>/crdt/<dataset_did>/heads.json
```

Op file format:
```
[ ciphertext bytes ]
```

Envelope stored separately for verification and sync.

---

## 7. Client Merge Model

Clients MUST:
* fetch `/crdt/heads` to determine missing ops
* call `/crdt/since` to receive encrypted batches
* decrypt via KEKₑ → DEK → plaintext ops
* merge ops via CRDT engine

Providers are **not** part of convergence logic.

---

## 8. Cryptography

Same crypto suite as Blob Service (RFC‑0002):

| Purpose | Algorithm |
|----------|-----------|
| Op encryption | XChaCha20‑Poly1305 |
| Key wrapping | X25519 |
| Envelope & CapToken signatures | Ed25519 |
| Vector Clock integrity | inside COSE body |

---

## 9. Error Codes

| Code | Meaning |
|-------|---------|
| `400` | Malformed request |
| `401` | Unauthorized / bad capability |
| `409` | Sequence regression / causality violation |
| `413` | Payload too large |
| `422` | Invalid vector clock |

---

## 10. Compliance Checklist

- [ ] Implements `/crdt/append`, `/crdt/since`, `/crdt/heads`
- [ ] Rejects unsigned or unauthorized ops
- [ ] Does not decrypt or interpret CRDT payload
- [ ] Maintains monotonic `seq` per dataset
- [ ] Maintains vector clock index
- [ ] Can stream historical ops to consumer

---

## 11. Open Questions

* Should providers prune ops older than N snapshots?
* Should CRDT batching be chunked by size or time?
* Should `/since` support delta‑compression?
* Should we define a standard CRDT engine for interop (Yjs vs Automerge)?

---

## 12. Change Log

| Version | Date | Changes |
|----------|------|----------|
| 0.1‑draft | 2025‑11‑06 | Initial draft, basic endpoints + envelope definition |
| 0.2‑draft | TBD | Define CBOR schema + multi‑chunk batching |

---

_End of RFC‑0003_

