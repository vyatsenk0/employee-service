package models

type Employee struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Phone string `json:"phone"`
    City  string `json:"city"`
}
