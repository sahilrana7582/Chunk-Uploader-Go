# 📂 Content Addressable Distributed File System (CADFS)

A high-performance, peer-to-peer content-addressable file system written in Go. CADFS enables distributed file storage, efficient retrieval, and encrypted data transmission — all built on top of a custom TCP-based networking library.

---

## 🚀 Overview

Unlike traditional file systems that locate data by filenames or directory paths, **CADFS identifies files using the hash of their content**, ensuring **deduplication**, **integrity**, and **versioning**.

This project includes:

- A custom TCP-based peer-to-peer communication layer
- Chunked file transfer for large files
- Secure file storage with encryption and decryption
- Redundant data replication for fault tolerance

---

## ✨ Features

- 🔒 **End-to-End Encryption**  
  All file data is encrypted before storage and decrypted on retrieval.

- 🧩 **Content Addressable Storage**  
  Files are split into chunks, hashed, and addressed by their content hash.

- 🌍 **Peer-to-Peer Networking**  
  Fully decentralized — no central server required. Built using raw TCP.

- 📦 **Chunked Data Streaming**  
  Files are transferred in chunks to support large file streaming and recovery.

- 🛡 **Redundancy & Fault Tolerance**  
  Each chunk can be replicated to multiple peers to prevent data loss.

---


## 📂 Project Structure

```bash
cadfs/
├── server/
│   ├── main.go          # TCP server to handle client connections
│   ├── handler.go       # Logic for UPLOAD, DOWNLOAD, CHUNK, etc.
│   ├── chunks/          # Stored file chunks (by hash)
│   ├── manifest/        # JSON files with file-to-chunk mappings
│   └── files/           # Original input files (if needed)
├── client/
│   └── main.go          # CLI client to interact with server
├── utils/
│   └── hash.go          # SHA-256 hash utilities
├── go.mod
└── README.md
