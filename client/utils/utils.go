package utils

import (
	"bufio"
	"bytes"
	"client/globals"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Mensaje struct {
	Mensaje string `json:"mensaje"`
}

type Paquete struct {
	Valores []string `json:"valores"`
}

func IniciarConfiguracion(filePath string) *globals.Config {
	var config *globals.Config
	configFile, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)

	return config
}

func LeerConsola() {
	// Leer de la consola
	reader := bufio.NewReader(os.Stdin)
	log.Println("Ingrese los mensajes (o vacio para terminar)")
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error al leer entrada:", err)
			continue
		}
		text = strings.TrimSpace(text)

		if text == "" {
			log.Println("Fin del programa.")
			break
		}
		log.Println("Ingresado:", text)
	}
}

func GenerarYEnviarPaquete() {
	paquete := Paquete{}

	// Leemos y cargamos el paquete
	reader := bufio.NewReader(os.Stdin)
	log.Println("Ingrese las líneas del paquete (una vacía para finalizar):")

	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error al leer entrada:", err)
			continue
		}
		text = strings.TrimSpace(text)

		if text == "" {
			log.Println("Entrada vacía detectada. Finalizando carga del paquete.")
			break
		}
		log.Println("Se anadio con exito la linea:", text, "al paquete")
		paquete.Valores = append(paquete.Valores, text)
	}
	log.Printf("paqute a enviar: %+v", paquete)
	EnviarPaquete(globals.ClientConfig.Ip, globals.ClientConfig.Puerto, paquete)
	// Enviamos el paqute
}

func EnviarMensaje(ip string, puerto int, mensajeTxt string) {
	mensaje := Mensaje{Mensaje: mensajeTxt}
	body, err := json.Marshal(mensaje)
	if err != nil {
		log.Printf("error codificando mensaje: %s", err.Error())
	}

	url := fmt.Sprintf("http://%s:%d/mensaje", ip, puerto)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("error enviando mensaje a ip:%s puerto:%d", ip, puerto)
	}

	log.Printf("respuesta del servidor: %s", resp.Status)
}

func EnviarPaquete(ip string, puerto int, paquete Paquete) {
	body, err := json.Marshal(paquete)
	if err != nil {
		log.Printf("error codificando mensajes: %s", err.Error())
	}

	url := fmt.Sprintf("http://%s:%d/paquetes", ip, puerto)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("error enviando mensajes a ip:%s puerto:%d", ip, puerto)
	}

	log.Printf("respuesta del servidor: %s", resp.Status)
}

func ConfigurarLogger() {
	logFile, err := os.OpenFile("tp0.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)
}
