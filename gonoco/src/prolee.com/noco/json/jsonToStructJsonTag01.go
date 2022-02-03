package main

import (
    "encoding/json"
    "fmt"
)

type Email struct {
    Email string `json:"email_address"`
}

type EmailsList struct {
    IsSchemaConforming bool    `json:"isSchemaConforming"`
    SchemaVersion      int     `json:"schemaVersion"`
    Emails             []Email `json:"unknown.0"`
}

func main() {
    jsonStr := `{
                    "isSchemaConforming":true,
                    "schemaVersion":0,
                    "unknown.0":[
                        {"email_address":"test1@uber.com"},
                        {"email_address":"test2@uber.com"}
                    ]
                }`

    emails := EmailsList{}
    json.Unmarshal([]byte(jsonStr), &emails)

    fmt.Printf("%+v\n", emails)
}
