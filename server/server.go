// Package server provides an implementation of a server that handles the expected http POST requests
package server

import (
   "encoding/json"
   "io/ioutil"
   "log"
   "net/http"

   "backend/model"
)

type Server interface {
   Start()
}

type server struct {
   server      *http.Server
   sessionData map[string]*model.Data // key is session id, value is the data corresponding to that sesssion
}

func postRequestHandler(sessionData map[string]*model.Data) http.HandlerFunc {
   return func(w http.ResponseWriter, req *http.Request) {
      if req.Method != http.MethodPost || req.Header.Get("Content-Type") != "application/json" {
         http.Error(w, "must be a POST request with a JSON format body", http.StatusBadRequest)
         return
      }

      body, err := ioutil.ReadAll(req.Body)
      if err != nil {
         log.Printf("ERROR in postRequestHandler: read request body error: %s\n", err.Error())
         http.Error(w, "Internal server error!", http.StatusInternalServerError)
         return
      }

      var e model.Event
      err = json.Unmarshal(body, &e)
      if err != nil {
         log.Printf("ERROR in postRequestHandler: unmarshal request body error: %s\n", err.Error())
         http.Error(w, "Internal server error!", http.StatusInternalServerError)
         return
      }

      complete := false
      if _, ok := sessionData[e.SessionId]; ok {
         complete = sessionData[e.SessionId].Fill(e)
      } else {
         sessionData[e.SessionId] = model.NewData(e) // assuming that TimeTaken event won't be the first event received for a session
      }

      if complete {
         log.Printf("Completed: %+v\n", *sessionData[e.SessionId])
         sessionData[e.SessionId].LogUrlHash()
      } else {
         log.Printf("%+v\n", *sessionData[e.SessionId])
      }

   }
}

func (s *server) Start() {
   err := s.server.ListenAndServe()
   if err != nil {
      log.Fatal(err)
   }
}

// New returns a Server. To start the server, run Start() on the server
// The argument host is the address of the server. If nil is passed in as the argument, the default address ":8080" will be used
func New(host *string) Server {
   mux := http.NewServeMux()

   addr := ":8080"
   if host != nil {
      addr = *host
   }

   s := &server{
      server: &http.Server{
         Addr:    addr,
         Handler: mux,
      },
      sessionData: make(map[string]*model.Data),
   }

   mux.HandleFunc("/", postRequestHandler(s.sessionData))

   return s
}