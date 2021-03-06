package pty_servers

import (
  "net"
  "strconv"
)

const READSIZE = 1024

type KeyServer struct {
  Port       int
  stay_alive chan int
  server     net.Listener
}

func NewKeyServer(client_limit int) (ks *KeyServer) {
  return &KeyServer{ stay_alive: make(chan int, client_limit) }
}

func (ks *KeyServer) Listen(port int, channel chan []byte, completed func()) {
  port_string := strconv.Itoa(port)
  server, err := net.Listen("tcp", ":"+port_string)
  if err != nil { panic(err) }

  ks.server = server
  ks.Port   = server.Addr().(*net.TCPAddr).Port

  for {
    conn, err := server.Accept()
    if err != nil { break }

    // put a placeholder in the stay-alive channel
    ks.stay_alive <- 1
    go ks.connection_to_channel(conn, channel)
  }
  completed()
}

func (ks *KeyServer) connection_to_channel(conn net.Conn, channel chan []byte) {
  for {
    bytes     := make([]byte, READSIZE)
    read, err := conn.Read(bytes)
    if err != nil { break }
    bytes = bytes[:read]
    channel <- bytes
  }

  // pop one off the stay live channel,
  //if empty kill the channel and the server
  _, _ = <-ks.stay_alive
  if (len(ks.stay_alive) == 0){
    close(channel)
    ks.server.Close()
  }
}


