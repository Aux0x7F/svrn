package main

import (
  "flag"
  "fmt"
  "os"
  "strings"
)

var (
  flagRoles     = flag.String("roles", "", "comma-separated: consumer,provider,relay,seed (default consumer-only)")
  flagServices  = flag.String("services", "", "comma-separated: blob,crdt")
  flagConfig    = flag.String("config", "", "path to YAML config (optional)")
  flagCommunity = flag.String("community", "", "bootstrap community URI (file/http i2p/git)")
  flagRouter    = flag.String("router", "auto", "i2p router: auto | external:host:port")
  flagVersion   = flag.Bool("version", false, "print version")
)

const version = "0.0.0-dev"

func main() {
  flag.Parse()
  if *flagVersion {
    fmt.Println("svrn", version)
    return
  }

  roles := parseList(*flagRoles)
  services := parseList(*flagServices)

  if len(roles) == 0 {
    fmt.Println("svrn: consumer mode (default). Try --roles provider,relay --services blob,crdt")
    os.Exit(0)
  }

  fmt.Printf("svrn starting roles=%v services=%v router=%s community=%s config=%s\n",
    roles, services, *flagRouter, *flagCommunity, *flagConfig)

  // TODO: load config, init router, DHT, services per RFC.
}

func parseList(s string) []string {
  s = strings.TrimSpace(s)
  if s == "" { return nil }
  parts := strings.Split(s, ",")
  out := make([]string, 0, len(parts))
  for _, p := range parts {
    v := strings.TrimSpace(p)
    if v != "" { out = append(out, v) }
  }
  return out
}