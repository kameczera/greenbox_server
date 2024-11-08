// client.go
package main

import (
    "fmt"
    "net"
)

func main() {
    // Conecta ao servidor
    conn, err := net.Dial("tcp", "localhost:8080")
    if err != nil {
        fmt.Println("Erro ao conectar ao servidor:", err)
        return
    }
    defer conn.Close()

    // Envia uma mensagem ao servidor
    message := "Chupa meu pau!"
    fmt.Printf("Enviando mensagem: %s\n", message)
    conn.Write([]byte(message))

    // Lê a resposta do servidor
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Erro ao ler resposta do servidor:", err)
        return
    }
    fmt.Printf("Resposta do servidor: %s\n", string(buffer[:n]))
}