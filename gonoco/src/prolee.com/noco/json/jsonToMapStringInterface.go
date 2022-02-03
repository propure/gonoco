package main

import (
    "encoding/json"
    "fmt"
)

func main() {
    jsonStr := `{
                    "isSchemaConforming":true,
                    "schemaVersion":0,
                    "unknown.0":[
                        {"email_address":"test1@uber.com"},
                        {"email_address":"test2@uber.com"}
                    ]
                }`

    dynamic := make(map[string]interface{})
    json.Unmarshal([]byte(jsonStr), &dynamic)

    firstEmail := dynamic["unknown.0"].([]interface{})[0].(map[string]interface{})["email_address"]

    fmt.Println(firstEmail)
}
