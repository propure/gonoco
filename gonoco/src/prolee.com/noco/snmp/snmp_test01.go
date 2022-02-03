package snmp

import (
    "fmt"
    "log"

    g "github.com/gosnmp/gosnmp"
)

func main() {
    // Default is a pointer to a GoSNMP struct that contains
    // sensible defaults eg port 161, community public, etc
    g.Default.Target = "172.18.0.2"
    g.Default.Community = "community"
    err := g.Default.Connect()
    if err != nil {
        log.Fatalf("Connect() err: %v", err)
    }
    defer g.Default.Conn.Close()

    oids := []string{"1.3.6.1.2.1.1.5.0", "1.3.6.1.2.1.1.7.0"}
    result, err := g.Default.Get(oids) // Get() accepts up to g.MAX_OIDS
    if err != nil {
        log.Fatalf("Get() err: %v", err)
    }

    for i, v := range result.Variables {
        fmt.Printf("%d. oid: %s ", i, v.Name)

        // the Value of each variable returned by Get() implements
        // interface{}. You could do a type switch...
        switch v.Type {
        case g.OctetString:
            fmt.Printf("string: %s\n", string(v.Value.([]byte)))
        default:
            // ... or often you're just interested in numeric values.
            // ToBigInt() will return the Value as a BigInt, for plugging
            // into your calculations.
            fmt.Printf("number: %d\n", g.ToBigInt(v.Value))
        }
    }
}
