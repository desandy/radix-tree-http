package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/gbrlsnchs/radix"
)

var tr = radix.New(0)

func AddIP(ip string, handler http.HandlerFunc) {
    key := "/ip?=" + ip
    tr.Add(key, handler)
}

func main() {
    AddIP("1.1.1.1", ipHandler)
    AddIP("8.8.8.8", ipHandler)
    AddIP("192.168.0.1", ipHandler)

    tr.Sort(radix.PrioritySort)

    http.HandleFunc("/", routeHandler)
    log.Println("Server running on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func routeHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
        return
    }
    full := r.URL.Path
    if raw := r.URL.RawQuery; raw != "" {
        full += "?" + raw
    }
    node, _ := tr.Get(full)
    if node == nil {
        http.NotFound(w, r)
        return
    }
    node.Value.(http.HandlerFunc)(w, r)
}

func ipHandler(w http.ResponseWriter, r *http.Request) {
    ip := r.URL.Query().Get("ip")
    fmt.Fprintf(w, "Matched registered IP: %s\n", ip)
}
