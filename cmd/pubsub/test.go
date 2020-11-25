package main

import (
    "net/url"
    "os"
    "os/signal"
    "log"
    "flag"
    "time"
    "fmt"
    "strconv"
    "strings"

    "github.com/gorilla/websocket"
)

// based off of client/server example code from gorilla github 
// here: https://github.com/gorilla/websocket/tree/master/examples/echo

var addr = flag.String("addr", "localhost:8080", "http service address")

var msgMinSub1 = -1
var msgMinSub2 = -1
var msgMaxSub1 = 0
var msgMaxSub2 = 0
var endFlag = false

var upgrader = websocket.Upgrader{}

func createPublisher(pNum string) {

    log.Println("Creating publisher number", pNum)
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)

    u := url.URL{Scheme: "ws", Host: *addr, Path: "/publish"}
    log.Printf("        connecting to publisher-server at: %s", u.String())

    c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        log.Fatal("     dial:", err)
    }

    defer c.Close()

    done := make(chan struct{})

    count := 0

    // publisher sends a message to the pubsub every 5 seconds
    ticker := time.NewTicker(time.Second * 5)
    defer ticker.Stop()

    for {
        if endFlag{
            return
        }
        select {
            case <-done:
                return
            case t := <-ticker.C:
                count += 1
                strCount := strconv.Itoa(count)
                content := "publisher: " + pNum + " msg: " + strCount + " at " + t.String()
                msg := []byte(content)
                err := c.WriteMessage(websocket.TextMessage, []byte(msg))
                if err != nil {
                    log.Println("     write:", err)
                    return
                }
            case <-interrupt:
                err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
                if err != nil {
                    log.Println("     write close:", err)
                    return
                }
                select {
                    case <-done:
                    case <-time.After(time.Second):
                }
                return
        }
    }
}

func createSubscribedClient(sNum string) {


    log.Println("Creating subscriber number", sNum)
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt)

    u := url.URL{Scheme: "ws", Host: *addr, Path: "/subscribe"}
    log.Printf("        connecting to publisher-server at: %s", u.String())

    c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
    if err != nil {
        log.Fatal("dial:", err)
    }

    defer c.Close()

    done := make(chan struct{})

    go func() {
        defer close(done)
        for {
            _, message, err := c.ReadMessage()
            if err != nil {
                log.Println("     read:", err)
                return
            }
            msg := "subscriber " + sNum + " received following: " + string(message)
            log.Printf("     recv: %s", msg)

            s := strings.Split(msg, " ")
            i, _:= strconv.Atoi(s[7])

            // track first and last message number the subscribers get from the publisher
            if sNum == "1"{

                if msgMinSub1 == -1{
                    msgMinSub1 = i
                }
                msgMaxSub1 = i
            } else {

                if msgMinSub2 == -1{
                    msgMinSub2 = i
                }

                msgMaxSub2 = i
            }

            // stop the first subscriber at the publisher's 10th message
            if i == 10 && sNum == "1"{
                return
            }
            // stop the second subscriber at the publisher's 12th message
            if i == 12 && sNum == "2"{
                // end the test
                endFlag = true
                return
            }

        }
    }()

    for {
        select {
        case <-done:
            return
        case <-interrupt:
            err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
            if err != nil {
                log.Println("write close:", err)
                return
            }
            select {
            case <-done:
            case <-time.After(time.Second):
            }
            return
        }
    }
}

func main() {
    flag.Parse()
    log.SetFlags(0)

    // add a second subscriber 
    // new subscriber should start getting messages at msg #4
    timer1 := time.NewTimer(time.Second * 16)
        go func() {
            <-timer1.C
            fmt.Println("Timer1 fired")
            go createSubscribedClient("2")
        }()

    go createSubscribedClient("1")
    createPublisher("1")

    log.Println("\n\nTest done, analyzing results\n")
    
    log.Println("TEST: Subscriber 1 received messages 1 - 10")
    if msgMinSub1 == 1 && msgMaxSub1 == 10{
        log.Println("PASS\n")
    } else {
        log.Println("FAIL\n")
    }

    log.Println("TEST: Subscriber 2 received messages 4 - 12")
    if msgMinSub2 == 4 && msgMaxSub2 == 12{
        log.Println("PASS\n")
    }else {
        log.Println("FAIL\n")
    }

}