package main
import (
    "net"
    "os"
)

var conn Conn
var tcpAddr TCPAddr 

func openConnection(servAddr string) {
    // servAddr := "localhost:6666"
    tcpAddr, err = net.ResolveTCPAddr("tcp", servAddr)
    if err != nil {
        println("ResolveTCPAddr failed:", err.Error())
        os.Exit(1)
    }
    conn, err = net.DialTCP("tcp", nil, tcpAddr)
    if err != nil {
        println("Dial failed:", err.Error())
        os.Exit(1)
    }
}

func closeConnection() {
    conn.Close()
}

func sendTcp(strEcho string) {
    _, err = conn.Write([]byte(strEcho))
    if err != nil {
        println("Write to server failed:", err.Error())
        os.Exit(1)
    }

    println("write to server = ", strEcho)

    reply := make([]byte, 1024)

    _, err = conn.Read(reply)
    if err != nil {
        println("Write to server failed:", err.Error())
        os.Exit(1)
    }

    println("reply from server=", string(reply))

}