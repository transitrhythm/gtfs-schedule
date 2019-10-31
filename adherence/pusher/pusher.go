package main

import "github.com/pusher/pusher-http-go"

func main(){
  client := pusher.Client{
    AppID: "836050",
    Key: "5d4a393761eabb3b834f",
    Secret: "0e20d72b0dbcafa89b34",
    Cluster: "us3",
    Secure: true,
  }

  data := map[string]string{"message": "hello world"}
  client.Trigger("my-channel", "my-event", data)
}
