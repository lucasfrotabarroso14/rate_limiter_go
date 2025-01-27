# Rate Limiter com Redis

Este projeto implementa um **Rate Limiter** (limitador de requisições) utilizando [Go](https://golang.org/), [Redis](https://redis.io/) e o pacote [go-redis](https://github.com/redis/go-redis). Ele foi criado para demonstrar como limitar acessos simultâneos tanto por **endereço IP** quanto por **chave de API** (API_KEY) de forma simples e escalável.

---

## Sumário

- [Como Funciona](#como-funciona)
    - [Visão Geral](#visão-geral)
    - [Limitação por IP](#limitação-por-ip)
    - [Limitação por Token/API_KEY](#limitação-por-tokenapikey)
    - [Bloqueio de IP ou Token](#bloqueio-de-ip-ou-token)
- [Estrutura do Projeto](#estrutura-do-projeto)
    - [main.go](#main-go)
    - [config.go](#config-go)
    - [limiter/rate_limiter.go](#limiterrate_limitergo)
    - [limiter/redis_store.go](#limiterredis_storego)
    - [Testes (rate_limiter_test.go)](#testes-rate_limiter_testgo)
- [Configuração e Instalação](#configuração-e-instalação)
    - [Variáveis de Ambiente (.env)](#variáveis-de-ambiente-env)
    - [Instalando Dependências](#instalando-dependências)
- [Uso](#uso)
    - [Executando o Servidor](#executando-o-servidor)
    - [Testando Manualmente](#testando-manualmente)
- [Testes Automatizados](#testes-automatizados)
- [Como Contribuir](#como-contribuir)
- [Licença](#licença)

---

## Como Funciona

### Visão Geral

1. O **Rate Limiter** é configurado com limites de requisições (por IP e por Token) e tempos de bloqueio para cada caso.
2. As requisições HTTP entram no servidor e passam pelo `MiddlewareHTTP` do **Rate Limiter**, que verifica:
    - Se existe um **API_KEY** no header.
    - Caso exista, aplica o limite de requisições por **Token**.
    - Caso não exista, aplica o limite de requisições por **IP**.
3. Se o limite de requisições for ultrapassado, o IP/Token é **bloqueado** por um tempo determinado.
4. Todas as informações de contagem e bloqueio são armazenadas no **Redis**.

### Limitação por IP

- Cada requisição **sem** cabeçalho `API_KEY` é identificada pelo endereço IP (`clientIP`).
- O Rate Limiter mantém um **contador** de requisições para aquele IP no Redis.
- Ao atingir o limite de requisições configurado (`IPRateLimit`), o IP é **bloqueado** por um período (`IPBlockTime`).

### Limitação por Token/API_KEY

- Quando a requisição contém o cabeçalho `API_KEY`, o Rate Limiter considera esse token no lugar do IP.
- É mantido um contador de requisições para esse **Token** no Redis.
- Ao atingir o limite de requisições configurado para Tokens (`TokenRateLimit`), a chave de API é **bloqueada** por um período (`TokenBlockTime`).

### Bloqueio de IP ou Token

- Se o IP/Token exceder o limite, ele é marcado como **bloqueado** por meio de uma chave `block:ip:xxx` ou `block:tokenXXX` no Redis.
- Novas requisições daquele IP/Token são respondidas automaticamente com **HTTP 429 Too Many Requests** enquanto durar o bloqueio.

---