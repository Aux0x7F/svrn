# RFC-0002: Sovereins Blob Service Protocol

**Status:** Draft (not yet ratified)
**Last updated:** 2025-11-06
**Authors:** 0x7F
**Applies to:** `svrn` Provider nodes, Blob client API, Envelope validation
**Related RFCs:** RFC-0001 (Core Architecture), RFC-0003 (CRDT Ops)

---

## 0. Abstract

This RFC defines the **Blob Service Protocol** — the layer of Sovereins responsible for storage and replication of immutable, encrypted data objects. Each object represents a *Last-Writer-Wins (LWW)* entry within a dataset.

This service provides the foundation for static file synchronization, versioned documents, and general encrypted file persistence in the Sovereins network.

---

## 1. Goals

* Enable encrypted data storage accessible over I2P with verifiable write authorization.
* Define a deterministic way to resolve conflicting writes (LWW semantics).
* Ensure all providers store data in a cryptographically verifiable form.
* Maintain end-to-end encryption—providers cannot decrypt content.

---

## 2. Design Summary

The Blob Service operates as an HTTP-based I2P endpoint that accepts and serves **Write Envelopes**. Each envelope contains:

1. **Capability Token (CapToken)** — signed authorization for write scope.
2. **WriteEnvelope** — signed metadata describing the blob and its lineage.
3. **Ciphertext** — AEAD-encrypted object body.

Providers validate the envelope, store it locally, and expose it via `/blob/get/<cid>` requests.

---

## 3. Endpoints

### 3.1 POST /blob/put

Stores a new encrypted blob.

#### Request Body
Binary concatenation of the following framed sections:
1. 4-byte length prefix (CapToken)
2. COSE_Sign1(CapToken)
3. 4-byte length prefix (WriteEnvelope)
4. COSE_Sign1(WriteEnvelope)
5. Ciphertext payload (remaining bytes)

#### Response
```json
{
  "ok": true,
  "cid": "<sha256>"
}
```

### 3.2 GET /blob/get/<cid>
Returns the ciphertext of the stored blob.

Headers:
```
Content-Type: application/octet-stream
X-SVRN-Envelope: base64(COSE_Sign1(WriteEnvelope))
```

### 3.3 GET /blob/index/<dataset_did>
Optional endpoint providing a signed index of known objects for the dataset.

```json
{
  "dataset": "did:key:...",
  "epoch": 3,
  "entries": [
    { "path": "/a/b.txt", "cid": "...", "seq": 42, "ts": 1730920000 }
  ],
  "sig": "base64(COSE_Sign1)"
}
```

---

## 4. Object Lifecycle

1. **Client encrypts content** with DEK (data encryption key).
2. **Client wraps DEK** with recipient KEKₑ (per dataset/keybag).
3. **Client constructs WriteEnvelope** with dataset DID, sequence number, and optional parent.
4. **Client signs WriteEnvelope** and attaches CapToken.
5. **Provider verifies both signatures** and capability validity.
6. **Provider stores** ciphertext, WriteEnvelope, and metadata.

---

## 5. Validation Rules

Providers MUST enforce the following rules before storing any blob:

| Check | Description | On failure |
|-------|--------------|-------------|
| Signature validity | CapToken and WriteEnvelope signatures valid | reject |
| Capability scope | Path and operation allowed in CapToken | reject |
| Expiry and caveats | `exp` and optional rate/size limits valid | reject |
| Dataset consistency | Envelope dataset == CapToken dataset | reject |
| Sequence check | `seq` > previous stored seq for path | reject (LWW) |
| CID match | `sha256(ciphertext)` == envelope CID | reject |
| Epoch validity | Envelope epoch <= dataset epoch | reject |

Rejected writes MUST NOT be logged with plaintext metadata.

---

## 6. Conflict Resolution (LWW)

Blobs are keyed by `(dataset_did, path)`.

When multiple valid envelopes exist for the same path, providers MUST select the one with the highest lexicographic `(seq, timestamp)` tuple.

Ties MAY be broken deterministically by lowest CID value.

---

## 7. Local Storage Format

### 7.1 Directory layout
```
<root>/blob/<dataset_did>/<epoch>/<sha256cid>.bin
<root>/blob/<dataset_did>/<epoch>/<sha256cid>.envelope
<root>/blob/<dataset_did>/index.json
```

### 7.2 Index schema
```json
{
  "path": "/a/b.txt",
  "cid": "...",
  "seq": 12,
  "epoch": 3,
  "ts": 1730920000
}
```

---

## 8. Cryptography

| Purpose | Algorithm | Notes |
|----------|------------|-------|
| Object encryption | XChaCha20-Poly1305 | AEAD, random nonce |
| Key wrapping | X25519 + HKDF | Dataset-specific KEKₑ |
| Record signing | Ed25519 | DIDs for nodes and datasets |
| Hashing | SHA-256 | CIDs and object IDs |

Ciphertext header MUST include the DEK-wrapped key and nonce in binary CBOR:
```cbor
{ dek_wrap: bstr, nonce: bstr }
```

---

## 9. Security Model

### Guarantees
* Providers cannot decrypt data.
* Consumers can verify write authenticity offline.
* Invalid or replayed writes are ignored.
* Datasets remain consistent under concurrent valid writes.

### Limitations
* Does not prevent Sybil attacks or storage omission.
* Does not guarantee global convergence (requires CRDT or consensus layer).

---

## 10. Error Codes

| Code | Meaning | Example |
|------|----------|----------|
| `400` | Malformed request | missing section or invalid frame |
| `401` | Unauthorized | invalid CapToken or signature |
| `409` | Conflict | seq/timestamp regression |
| `413` | Payload too large | exceeds CapToken limit |
| `429` | Rate limited | per CapToken caveat |
| `500` | Internal error | unexpected runtime issue |

---

## 11. Compliance Checklist

- [ ] Implements all endpoints `/blob/put`, `/blob/get/<cid>`
- [ ] Validates both CapToken and WriteEnvelope signatures
- [ ] Enforces LWW policy
- [ ] Rejects unsigned or invalid envelopes
- [ ] Stores only ciphertext (no plaintext path or content)
- [ ] Supports per-dataset index export (optional)

---

## 12. Change Log

| Version | Date | Changes |
|----------|------|----------|
| 0.1-draft | 2025-11-06 | Initial draft, LWW semantics, framing format |
| 0.2-draft | TBD | Define CBOR schema, compression, multi-chunk blobs |

---

_End of RFC-0002_

