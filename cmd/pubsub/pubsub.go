package main

import (
    "flag"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "fmt"
    "sync"
    "strconv"

    "github.com/gorilla/websocket"
)

// based off of client/server example code from gorilla github 
// here: https://github.com/gorilla/websocket/tree/master/examples/echo

var addr = flag.String("addr", "localhost:8080", "http service address")
var upgrader = websocket.Upgrader{}
var subscribers = []*websocket.Conn{}
var publisherCount = 0
var muSubscribers = sync.Mutex{}

func subscribe(w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }

    muSubscribers.Lock()
    subscribers = append(subscribers, c)
    log.Println("New subscriber connection stored! Pub/Sub subscriber count ", strconv.Itoa(len(subscribers)))
    muSubscribers.Unlock()

}

// establish ws with publisher, then when it sends messages to the ws, forward them to the clients
func publish(w http.ResponseWriter, r *http.Request) {
    c, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Print("upgrade:", err)
        return
    }
    publisherCount += 1
    log.Println("Publisher count", strconv.Itoa(publisherCount))

    defer c.Close()
    for {
        log.Println("\nForwarding new message to subscribers")
        mt, message, err := c.ReadMessage()
        if err != nil {
            log.Println("read:", err)
            break
        }

        log.Printf("recv: %s", message)


        muSubscribers.Lock()

        var subscribersRefresh = []*websocket.Conn{}

        // forwarding new publisher content to all subscribers
        for i := 0; i < len(subscribers); i++ {
            subscriber := subscribers[i]
            err = subscriber.WriteMessage(mt, message)
            if err != nil {
                log.Println("error writing to subscriber, removing:", err)
                continue
            }

            subscribersRefresh = append(subscribersRefresh, subscriber)
        }

        // remove disappeared subscribers 
        subscribers = subscribersRefresh
        log.Println("Current subscriber count: ", strconv.Itoa(len(subscribers)))
        muSubscribers.Unlock()
    }

    publisherCount -= 1
}

func home(w http.ResponseWriter, r *http.Request) {
    log.Println("Root established!")
    fmt.Fprintf(w, "subscriber count: %d\n", len(subscribers))
    fmt.Fprintf(w, "publisher count: %d\n", publisherCount)
}

func cleanup() {
    for i := 0; i < len(subscribers); i++ {
     log.Println("Closing a connection")
         subscriberConnection := subscribers[i]
         subscriberConnection.Close()
    }
}

func main() {

    c := make(chan os.Signal)
    signal.Notify(c, os.Interrupt, syscall.SIGTERM)
    go func() {
        <-c
        cleanup()
        os.Exit(1)
    }()

    flag.Parse()
    log.SetFlags(0)

    log.Println("Starting Publisher-Server")
    http.HandleFunc("/subscribe", subscribe)
    http.HandleFunc("/publish", publish)
    http.HandleFunc("/", home)
    log.Fatal(http.ListenAndServe(*addr, nil))
}