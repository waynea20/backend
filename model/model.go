// Package model defines the Data and Dimension structs as specified in the readme.
// It also defines an Event struct that maps to the expected http request body, and
// methods that transform an Event to a Data
package model

import (
   "log"
)

type Event struct {
   Type string `json:"eventType"`

   // assumed to have the same value for data of same session Id
   WebsiteUrl string `json:"websiteUrl"`
   SessionId  string `json:"sessionId"`

   // screen resize event values
   ResizeFrom Dimension `json:"resizeFrom"`
   ResizeTo   Dimension `json:"resizeTo"`

   // copy and paste event value
   FieldId string `json:"formId"`

   // time taken event value
   Time int `json:"time"` // Seconds
}

type Data struct {
   WebsiteUrl         string
   SessionId          string
   ResizeFrom         Dimension
   ResizeTo           Dimension
   CopyAndPaste       map[string]bool // map[fieldId]true
   FormCompletionTime int             // Seconds
}

type Dimension struct {
   Width  string `json:"width"`
   Height string `json:"height"`
}

const (
   // Event types
   CopyAndPasteEvent = "copyAndPaste"
   TimeTakenEvent    = "timeTaken"
   ScreenResizeEvent = "screenResize"
)

// NewData creates and fills a data with event data
//
// Assumes that TimeTaken event won't be the first event received for
// a session (so that an incomplete struct is not treated as a complete struct)
func NewData(e Event) *Data {
   d := &Data{
      WebsiteUrl: e.WebsiteUrl,
      SessionId:  e.SessionId,
   }

   d.Fill(e)

   return d
}

// Fill fills Data d with the new values from Event e based on the event type
// assumptions: screen resizing event has "screenResize" as event type
//              websiteUrl value will be the same for data with same session Id
//              receiving a TimeTaken event means that the submit button has been clicked
func (d *Data) Fill(e Event) bool {
   if d.SessionId != e.SessionId {
      return false
   }

   switch e.Type {
   case CopyAndPasteEvent:
      if d.CopyAndPaste == nil {
         d.CopyAndPaste = make(map[string]bool)
      }
      d.CopyAndPaste[e.FieldId] = true
   case TimeTakenEvent:
      d.FormCompletionTime = e.Time
      return true
   case ScreenResizeEvent:
      d.ResizeFrom = e.ResizeFrom
      d.ResizeTo = e.ResizeTo
   }

   return false
}

func (d *Data) LogUrlHash() {
   log.Printf("Hash value of url %s: %d", d.WebsiteUrl, hash(d.WebsiteUrl))
}

func hash(s string) int {
   hash := 0

   for _, r := range s {
      hash = 31*hash + int(r)
   }

   return hash
}
