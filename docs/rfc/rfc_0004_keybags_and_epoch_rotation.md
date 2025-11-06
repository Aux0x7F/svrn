# RFC-0004: Sovereins Keybags, Epoch Rotation, and Revocation Model

**Status:** Draft (not yet ratified)
**Last updated:** 2025-11-06
**Authors:** Sovereins Project Maintainers
**Applies to:** `svrn` client + provider nodes, dataset key lifecycle, capability validity
**Related RFCs:** RFC-0001 (Core Architecture), RFC-0002 (Blob Service), RFC-0003 (CRDT Ops)

---

## 0. Abstract

This RFC defines **dataset key management** in Sovereins: the encryption lifecycle for datasets, including keybag format, epoch rotation, member onboarding, revocation, and client responsibilities.

In Sovereins, **datasets are encrypted**, not nodes. Therefore, membership is enforced by **possession of decryption capability**, not network identity. A dataset may be accessed by any node able to unwrap the appropriate key material, regardless of hosting location or provider trust.

Keybags provide the distribution layer for these keys.

---

## 1. Goals

* Provide a cryptographically secure mechanism for distributing per-dataset decryption keys.
* Allow dataset owners to revoke access and rotate keys **without provider cooperation**.
* Support multi-member datasets without central authority.
* Prevent providers from accessing or deriving dataset keys.
* Avoid global re-encryption during membership changes (only per-epoch re-encryption).

---

## 2. Dataset Encryption Model

Each dataset has:

| Field | Meaning |
|--------|---------|
| **DEK** | Data Encryption Key (symmetric, rotates each epoch) |
| **KEKᵣ** | Key Encryption Key for each recipient (derived from member DID) |
| **Epoch** | Monotonic version integer for revocation/rotation |

Encryption flow:
```
plaintext → AEAD(DEK) → ciphertext
DEK → X25519 wrap → per-recipient encrypted DEK entry (keybag)
```

Providers store ciphertext only. DEK and KEK never appear in provider-visible storage.

---

## 3. Keybag Definition

A keybag is an encrypted map of: `recipient_did → wrapped_DEK`.

### 3.1 Binary Format (CBOR)
```cbor
{
  v: 1,
  dataset: "did:key:...",
  epoch: 3,
  alg: "xchacha20poly1305",
  wraps: [
    { did: "did:key:z6Mk...", wrap: bstr },
    { did: "did:key:z6Ln...", wrap: bstr }
  ],
  sig: COSE_Sign1        ; owner-signed keybag integrity
}
```

### 3.2 Wrap Format
```cbor
{
  kek_salt: bstr,      ; HKDF salt
  enc_dek: bstr        ; AEAD(DEK, kek_derived)
}
```

`kek_derived = HKDF( X25519(owner_priv, recipient_pub), salt )`

---

## 4. Epoch Rotation

Epoch rotation is the revocation model for Sovereins.

```
old epoch → revoked
new epoch → new DEK, new keybag
```

### Procedure
1. Owner creates new DEK
2. Owner wraps DEK for each active member → new keybag
3. Owner publishes updated `DatasetRecord` with `epoch += 1`
4. Providers MUST reject writes with older epoch values
5. Consumers MAY continue reading old encrypted objects if they have old epoch keys

### Notes
* Rotation does **not** require re-encrypting all stored ciphertext immediately
* Clients MAY support lazy re-encryption or dual-key reads
* Providers are not involved in rotation logic beyond epoch enforcement

---

## 5. Member Revocation

To revoke a member:

1. Owner generates new epoch (as above)
2. New keybag omits revoked member DID
3. No provider action required
4. Revoked member loses ability to decrypt future content

### Important
* Revocation does **not** remove access to previously downloaded ciphertext
* Secure deletion is out of scope (requires per-object re-encryption policies)

---

## 6. New Member Onboarding

1. Owner adds member DID to next or current epoch keybag
2. Owner republishes keybag
3. Member receives wrapped DEK for current epoch via:
   * direct transfer, or
   * dataset metadata pointer

Providers are not part of onboarding.

---

## 7. Client Responsibilities

Consumers MUST:
* detect dataset epoch change via updated `DatasetRecord`
* fetch matching keybag (pointer supplied in record)
* unwrap correct DEK and cache locally
* discard expired epochs if `exp` is set

Providers MUST:
* enforce epoch sequencing for writes
* ignore writes with `epoch < current epoch` for dataset

---

## 8. Storage Expectations

Keybags are **not stored in the DHT**.
Instead they are referenced via pointer in the DatasetRecord:

```cbor
{
  dataset: "did:key:...",
  epoch: 3,
  keybag_index: "bafkreia..."       ; CID or https:// or file:// or i2p://
}
```

Per RFC-0001, providers are forbidden from caching keybags.

---

## 9. Security Properties

| Property | Mechanism |
|----------|-----------|
| Provider cannot decrypt data | DEK never exposed; AEAD enforced |
| Member-based access control | Per-DID wrapped keys in keybag |
| Revocation without trust | Epoch bump + keybag omission |
| Compromise-localized | DEK leak affects only one epoch, not full dataset |

---

## 10. Limitations

* Does not prevent ex-members from retaining previously synced ciphertext
* Does not guarantee key erasure from compromised consumers
* Does not address multi-owner datasets (handled in future RFC)

---

## 11. Compliance Checklist

- [ ] Client unwraps DEK using X25519 + HKDF
- [ ] Providers reject lower-epoch writes
- [ ] Keybags are never cached or served by providers
- [ ] DatasetRecord MUST include `epoch` and `keybag_index`
- [ ] DEK MUST NOT be reused across epochs

---

## 12. Open Questions

* Should keybags support on-chain revocation attestations?
* Should we support HPKE as alternate wrapping mechanism?
* Should DEK derive from epoch via KDF instead of independent randomness?

---

## 13. Change Log

| Version | Date | Changes |
|----------|------|---------|
| 0.1-draft | 2025-11-06 | Initial definition, lifecycle + keybag format |
| 0.2-draft | TBD | Define multi-owner keybag signing, compressed format |

---

_End of RFC-0004_