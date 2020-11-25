# go-publisher-subscriber

1. clone repo, navigate to cmd/pubsub
2. go run pubsub.go
3. start a new terminal window, navigate to the same folder
4. go run test.go
5. test code completes on its own, pubsub.go must be stopped with a control-c interrupt 

Notes:
-websocket setup with channels is based on gorilla server/client example
-implementation of mutex lock on subscribers[] slice is not a performance problem for small numbers of publishers and subscribers, but would introduce latency issues at scale, and would need to be reworked 
-no implementation of topics, but would be relatively straight forward to add
