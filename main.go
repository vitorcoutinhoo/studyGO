package main

import (
	"encoding/json"
	"log"
	"net/http"
	_ "main.go/docs"  //não funfa

	httpSwagger "github.com/swaggo/http-swagger"
)

// Response representa a resposta da API
type Response struct {
	Message string `json:"message" example:"Olá! API funcionando perfeitamente"`
}

// @title API Simples
// @version 1.0
// @description API de exemplo com um endpoint
// @host localhost:8080
// @BasePath /

// @Summary Endpoint de saudação
// @Description Retorna uma mensagem de boas-vindas
// @Tags hello
// @Produce json
// @Success 200 {object} Response
// @Router /hello [get]
func hello(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "Olá! API funcionando perfeitamente",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	log.Println("Servidor rodando em http://localhost:8080")
	log.Println("Acesse: http://localhost:8080/hello")
	log.Println("Swagger UI: http://localhost:8080/swagger/index.html")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
