# Online-Offline Indicator
Design implementation of online-offline indicator system. A client would send a heartbeat periodically and server will store it on redis(in a key-value pair, key will be user ID and value will be current time) with expiration time. Once the time is expired, that entry will be deleted from redis, hence, we can say user is offline and if entry is present means that user is still online.

## How to run?
- Install redis and run it on port `6379`
- Run command `go run cmd/main.go`
