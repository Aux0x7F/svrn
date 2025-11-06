# RFC-0001: svrn  Sovereins, a Decentralized Community Cloud

**Status:** Draft  
**Created:** 2025-11-06  
**Repository:** github.com/aux0x7F/svrn  
**Binary:** svrn  
**Subtext:** *Sovereins  a decentralized community cloud*

## 0. Summary
svrn is a decentralized, community-scoped cloud over I2P with a Kademlia DHT for discovery.
One universal binary runs as consumer / provider / relay / seed via config.
Access is capability- and key-based. MVP ships LWW blobs + CRDT ops.
(No star topology; no clearnet ingress.)

## 1. Goals
- Community-scoped overlay; private by default.
- DID (did:key) identities; COSE-signed records.
- Two data modes: LWW blobs & CRDT ops.
- One binary; roles via config/flags.

## 2. Non-Goals
Centralized auth, clearnet exposure, global identity authority, massive (>500) member datasets in MVP.

## 3. Roles
consumer (default), provider, relay, seed  same binary, selected by config.

## 416
(Expanded in the full architecture doc; to be filled alongside implementation.)