# imlog

Immutable log chain for Go.

A lightweight logging library that guarantees **log integrity using SHA256 hash chaining**.

Each log entry includes a hash calculated from:

```
SHA256(current_log_content + previous_log_hash)
```

This creates a **tamper-evident chain**, where modifying any past entry invalidates all subsequent logs.

---

# Table of Contents

- [Português](#português)
  - [Visão Geral](#visão-geral)
  - [Como Funciona](#como-funciona)
  - [Instalação](#instalação)
  - [Uso Básico](#uso-básico)
  - [Estrutura do Log](#estrutura-do-log)
  - [Verificação de Integridade](#verificação-de-integridade)
  - [Casos de Uso](#casos-de-uso)
- [English](#english)
  - [Overview](#overview)
  - [How It Works](#how-it-works)
  - [Installation](#installation)
  - [Basic Usage](#basic-usage)
  - [Log Structure](#log-structure)
  - [Integrity Verification](#integrity-verification)
  - [Use Cases](#use-cases)
- [License](#license)

---

# Português

## Visão Geral

O **imlog** é uma biblioteca de logging para Go que garante **imutabilidade e integridade criptográfica dos logs**.

Cada entrada de log contém um **hash SHA256** calculado a partir do conteúdo do log atual e do hash do log anterior.

Isso cria uma **cadeia criptográfica**, semelhante ao conceito usado em blockchains.

Qualquer modificação em um log antigo quebra toda a cadeia de hashes.

---

## Como Funciona

Cada novo log calcula seu hash usando o seguinte esquema:

```
hash = SHA256(
    log_conteudo +
    hash_do_log_anterior
)
```

Cadeia de logs:

```
Log0 -> Log1 -> Log2 -> Log3
  |      |      |      |
 H0     H1     H2     H3
```

Onde:

```
H1 = SHA256(Log1 + H0)
H2 = SHA256(Log2 + H1)
H3 = SHA256(Log3 + H2)
```

Se qualquer log for alterado:

- o hash daquele log muda
- todos os hashes seguintes deixam de ser válidos

Isso permite detectar **qualquer adulteração histórica**.

---

## Instalação

```
go get github.com/thomasbscj/imlog
```

---

## Uso Básico

Exemplo simples de criação de logs encadeados:

```
package main

import (
    "github.com/thomasbscj/imlog"
)

func main() {

    logger := imlog.New()

    logger.Log("user created")
    logger.Log("password changed")
    logger.Log("login success")

}
```

Cada chamada gera uma nova entrada encadeada ao log anterior.

---

## Estrutura do Log

Um log possui os seguintes campos:

```
type Entry struct {
    Index     uint64
    Timestamp int64
    Message   string
    PrevHash  [32]byte
    Hash      [32]byte
}
```

Campos:

| Campo | Descrição |
|------|-----------|
| Index | posição na cadeia |
| Timestamp | momento da criação |
| Message | conteúdo do log |
| PrevHash | hash do log anterior |
| Hash | hash do log atual |

---

## Verificação de Integridade

A cadeia de logs pode ser verificada recalculando todos os hashes.

Exemplo:

```
valid := logger.Verify()

if !valid {
    panic("log corruption detected")
}
```

O processo de verificação:

1. recalcula cada hash
2. compara com o hash armazenado
3. verifica se cada log referencia corretamente o anterior

---

## Casos de Uso

Essa biblioteca é útil para sistemas que exigem **auditabilidade forte**:

- trilhas de auditoria
- logs financeiros
- sistemas de compliance
- registros de segurança
- sistemas regulatórios
- logs críticos de infraestrutura

---

# English

## Overview

**imlog** is a Go logging library that guarantees **cryptographic integrity and immutability of logs**.

Each log entry contains a SHA256 hash calculated from:

```
current_log_content + previous_log_hash
```

This forms a **hash chain**, similar to concepts used in blockchains.

If any previous log entry is modified, the entire chain becomes invalid.

---

## How It Works

Each new log computes its hash using:

```
hash = SHA256(
    log_content +
    previous_log_hash
)
```

Log chain example:

```
Log0 -> Log1 -> Log2 -> Log3
  |      |      |      |
 H0     H1     H2     H3
```

Where:

```
H1 = SHA256(Log1 + H0)
H2 = SHA256(Log2 + H1)
H3 = SHA256(Log3 + H2)
```

If any entry is modified:

- its hash changes
- every subsequent hash becomes invalid

This makes tampering **detectable**.

---

## Installation

```
go get github.com/thomasbscj/imlog
```

---

## Basic Usage

Example:

```
package main

import (
    "github.com/thomasbscj/imlog"
)

func main() {

    logger := imlog.New()

    logger.Log("user created")
    logger.Log("password changed")
    logger.Log("login success")

}
```

Each log entry becomes part of the cryptographic chain.

---

## Log Structure

```
type Entry struct {
    Index     uint64
    Timestamp int64
    Message   string
    PrevHash  [32]byte
    Hash      [32]byte
}
```

Fields:

| Field | Description |
|------|-------------|
| Index | log position in chain |
| Timestamp | creation time |
| Message | log content |
| PrevHash | previous log hash |
| Hash | current log hash |

---

## Integrity Verification

The log chain can be validated by recalculating hashes.

Example:

```
valid := logger.Verify()

if !valid {
    panic("log corruption detected")
}
```

Verification steps:

1. recompute hashes
2. compare with stored hashes
3. verify correct chaining

---

## Use Cases

- audit logs
- financial logging
- compliance systems
- security event logging
- regulatory systems
- infrastructure event tracking

---

# License

MIT