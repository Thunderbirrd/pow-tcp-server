# pow-tcp-server
TCP server &amp; client with Proof-of-work

## Task
Design and implement "Word of Wisdom" tcp server:

- TCP server should be protected from DDOS attacks with the Prof of Work (https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge.


## Run

To run both server && client use:
```
docker-compose up -d
```
