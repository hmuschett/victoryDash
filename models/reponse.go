package models

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//Response represent a standar response
type Response struct {
	Status   int         `json:"status"`
	Menssage string      `json:"menssage"`
	Data     interface{} `json:"data"`
	writer   http.ResponseWriter
}

func CreateDefaultResponse(w http.ResponseWriter) Response {
	return Response{Status: http.StatusOK, Menssage: "Susseful", writer: w}
}
func SendNotFound(w http.ResponseWriter) {
	response := CreateDefaultResponse(w)
	response.NotFound()
	response.Send()
}
func (this *Response) NotFound() {
	this.Status = http.StatusNotFound
	this.Menssage = "Resource not found!!"
}
func SendNotContent(w http.ResponseWriter) {
	response := CreateDefaultResponse(w)
	response.NotContetnt()
	response.Send()
}
func (this *Response) NotContetnt() {
	this.Status = http.StatusNoContent
	this.Menssage = "StatusNoContent!!"
}
func SendUnprocessableEntity(w http.ResponseWriter) {
	response := CreateDefaultResponse(w)
	response.UnprocessableEntity()
	response.Send()
}
func (this *Response) UnprocessableEntity() {
	this.Status = http.StatusUnprocessableEntity
	this.Menssage = "Unprocessable Entity!!"
}
func SendData(w http.ResponseWriter, data interface{}) {
	response := CreateDefaultResponse(w)
	response.Data = data
	response.Send()
}
func (response *Response) Send() {
	response.writer.Header().Set("content-type", "application/json")
	response.writer.WriteHeader(response.Status)

	output, _ := json.Marshal(&response)
	fmt.Fprintf(response.writer, string(output))
}
