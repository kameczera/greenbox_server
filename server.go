// server.go
package main

import (
    "bufio"
    "fmt"
    "net"
    "os"
    "log"
    "strings"
    "time"
)

func main() {

    if len(os.Args) < 2 {
        fmt.Println("Por favor, informe o nome do arquivo de lista de servidores e o id do servidor.")
    }

    serverMap := readServerList(os.Args[1])

    var keys []string

    // Itera sobre o mapa e adiciona as chaves ao slice
    for key := range serverMap {
        keys = append(keys, key)
    }

    selfPort := os.Args[2]
    
    listener, err := net.Listen("tcp", fmt.Sprintf(":%s", selfPort))
    if err != nil {
        fmt.Println("Erro ao iniciar o servidor:", err)
        return
    }
    defer listener.Close()
    
    fmt.Println("Dicionário de servidores:", serverMap)
    fmt.Println("Servidor ouvindo na porta", selfPort, serverMap[selfPort])

    fmt.Println("Decidindo líder...")
    for {
        decideLeader(serverMap, keys, selfPort)
        time.Sleep(3 * time.Second)
    }
}

func decideLeader(serverMap map[string]string, keys []string, selfPort string) {
    connectionArray := make([]bool, len(serverMap))

    for i := range connectionArray {
        connectionArray[i] = true
    }
    
    index := 0
    for index < len(keys) {
        if keys[index] == selfPort {
            index++
            continue
        }
        fmt.Printf("Tentando conectar ao servidor %s\n", keys[index])
        conn, err := net.DialTimeout("tcp", fmt.Sprintf("localhost:%s", keys[index]), 5 * time.Second)
        if err != nil {
            fmt.Printf("%s tentou conectar ao servidor %s e falhou\n", selfPort, keys[index])
            connectionArray[index] = false
        } else {
            fmt.Printf("%s tentou conectar ao servidor %s e conseguiu\n", selfPort, keys[index])
            connectionArray[index] = true
            conn.Close()
            break
        }
        index++  // Move para o próximo servidor
    }
}


func readServerList(f string) map[string]string {
    serverMap := make(map[string]string)
    
    file, err := os.Open(f) // Usa uma variável diferente para o arquivo
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        serverInfo := strings.Split(line, ",")
        serverMap[serverInfo[0]] = serverInfo[1]
    }

    if err := scanner.Err(); err != nil {
        fmt.Println("Erro:", err)
    }

    return serverMap
}

func handleConnection(conn net.Conn) {
    defer conn.Close()
    buffer := make([]byte, 1024)
    n, err := conn.Read(buffer)
    if err != nil {
        fmt.Println("Erro ao ler do cliente:", err)
        return
    }
    fmt.Printf("Mensagem recebida: %s\n", string(buffer[:n]))

    // Envia uma resposta para o cliente
    response := "Olá, cliente! Mensagem recebida com sucesso."
    conn.Write([]byte(response))
}