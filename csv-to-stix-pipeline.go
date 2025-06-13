package main

import (
    "bufio"
    "encoding/csv"
    "encoding/json"
    "fmt"
    "os"
    "time"

    "github.com/panther-labs/stix2"
)

func main() {
    f, err := os.Open("feed.csv")
    if err != nil {
        fmt.Fprintln(os.Stderr, "Failed to open CSV:", err)
        os.Exit(1)
    }
    defer f.Close()

    reader := csv.NewReader(bufio.NewReader(f))
    rows, err := reader.ReadAll()
    if err != nil {
        panic(err)
    }

    // Create STIX collection
    coll := stix2.New()

    for _, row := range rows {
        ip, source, score, ts := row[0], row[1], row[2], row[3]

        // Create IPv4Address SCO
        ipv4, _ := stix2.NewIPv4Address(ip)
        coll.Add(ipv4)

        // Create Indicator object
        indicator := stix2.NewIndicator(
            stix2.WithName(fmt.Sprintf("Malicious IP %s", ip)),
            stix2.WithPattern(fmt.Sprintf(`[ipv4-addr:value = '%s']`, ip)),
            stix2.WithCreated(time.Now().UTC()),
            stix2.WithModified(time.Now().UTC()),
        )
        // Add meta-data using extensions or description
        indicator.Description = fmt.Sprintf("source=%s; score=%s", source, score)
        coll.Add(indicator)
    }

    bundle := coll.ToBundle()
    data, _ := json.MarshalIndent(bundle, "", "  ")
    fmt.Println(string(data))
}
