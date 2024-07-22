# Simulation of race condition and double query

This is to simulate race conditions and double queries, using mysql and redis

### How to run
1. Initialize the .env file and recreate the table using:
```bash
go run main.go
```
2. Test race conditions using:
```bash
go test -race main_test.go
```
