# Word of Wisdom POW (Proof of work)

# Task

Test task for Server Engineer

Design and implement “Word of Wisdom” tcp server.  
• TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.  
• The choice of the POW algorithm should be explained.  
• After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
• Docker file should be provided both for the server and for the client that solves the POW challenge

# Demo

[![asciicast](https://asciinema.org/a/4ShE41NkxZizsUcFb8D6bbayK.svg)](https://asciinema.org/a/4ShE41NkxZizsUcFb8D6bbayK)

# Using

## Method #1 (Makefile)

```sh
make run_client
make run_server
```

## Method #2 (docker-compose)

```sh
docker-compose up
```

# Implementation

## Project structure

- docker - dockerfiles for build
- cmd
  - server - server side app
  - client - client side app
- internal
  - app - implement app singleton with factory 
  - book - service containing service (business logic) and repository (data extraction) 
  - handshake - sha256 based implementation of proof-of-work challenge response protocol
  - domain - folder containing domain object that can be possibly shared across multiple packages
  - client - client struct impl
  - server - server struct impl

## Algorithm choice explanation 
- It's well-spready and easy to implement in test task.
- SHA-256 can be part of a challenge-response PoW algorithm, but the overall effectiveness against DDoS attacks depends on the design, parameters, and other factors, not just the hash function choice.
- SHA-256 is a widely recognized and secure cryptographic hash function, providing a strong foundation for PoW algorithm security.
