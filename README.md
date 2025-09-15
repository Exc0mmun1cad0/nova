[![CI status](https://github.com/Exc0mmun1cad0/nova/actions/workflows/ci.yml/badge.svg)](https://github.com/Exc0mmun1cad0/nova/actions/workflows/ci.yml)
![Repository Top Language](https://img.shields.io/github/languages/top/Exc0mmun1cad0/nova)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Exc0mmun1cad0/nova)
![Github Repository Size](https://img.shields.io/github/repo-size/Exc0mmun1cad0/nova)
![License](https://img.shields.io/badge/license-MIT-green)
![GitHub last commit](https://img.shields.io/github/last-commit/Exc0mmun1cad0/nova)

# Nova

> My implementation of a Redis-like key-value in-memory database

Nova is a **key-value in-memory Redis-compatible database** written in Go.

## Features
- Communication via **RESP (Redis Serialization Protocol)**
- High-performance in-memory storage

## Supported data types
- **String**
- **Int**
- **List** (of strings)

## TODO
- `INFO` command
- **Write-Ahead Log (WAL)** for persistence