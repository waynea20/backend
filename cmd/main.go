package main

import (
   "backend/server"
)

func main() {
   server := server.New(nil)
   server.Start()
}
