# go-publisher-subscriber

Instructions:  <br/>

1. clone repo, navigate to cmd/pubsub
2. go run pubsub.go
3. start a new terminal window, navigate to the same folder
4. go run test.go
5. test code completes on its own, pubsub.go must be stopped with a control-c interrupt 

Notes: <br/> <br/>
 -websocket setup with channels is based on gorilla server/client example: https://github.com/gorilla/websocket/tree/master/examples/echo <br/> <br/>
  -implementation of mutex lock on subscribers[] slice is not a performance problem for small numbers of publishers and subscribers, but would introduce latency issues at scale, and would need to be reworked <br/> <br/>
 -no implementation of topics, but would be relatively straight forward to add
