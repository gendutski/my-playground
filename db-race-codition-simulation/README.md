# Simulation of race condition and double query

This is to simulate race conditions and multiple queries, using mysql and redis

The case example used represents an e-wallet, where the topup or balance transfer feature works as expected

### How to run
1. Initialize the .env file and recreate the table using:
```bash
go run main.go
```
2. Test race conditions using:
```bash
go test -race main_test.go
```
