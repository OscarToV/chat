package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // clientes conectados
var broadcast = make(chan Message)           // canal broadcast

//toma una conexion HTPP y la actualiza a un websocket
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Objeto mensaje
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func main() {

	fs := http.FileServer(http.Dir("../public"))
	http.Handle("/", fs)

	// Configurar la ruta websocket
	http.HandleFunc("/ws", handleConnections)

	// escuchamos mensajes recibidos
	go handleMessages()

	// Iniciamos el servidor en localhost en el puerto 8000
	log.Println("http server started on :5000")
	err := http.ListenAndServe(":5000", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Actualizamos el request Get inicial a un socket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	// Nos aseguramos de cerrar la conexion
	defer ws.Close()

	// Registramos a un nuevo cliente
	clients[ws] = true

	for {
		var msg Message
		// leemos el mensaje con formato json y los mapeamos
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
			break
		}
		// enviamos en mensaje recibido al canal broadcast
		broadcast <- msg
	}
}

func handleMessages() {
	for {
		// tomamos el siguiente mensaje desde el canal broadcast
		msg := <-broadcast
		// lo mandamos a todos los clientes que actualmente estan conectados
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %v", err)
				client.Close()
				delete(clients, client)
			}
		}
	}
}
