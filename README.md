## Backend (Go)

### Part 1

The Backend should:

1. Create a Server
2. Accept POST requests in JSON format similar to those specified above
3. Map the JSON requests to relevant sections of the data struct (specified below)
4. Print the struct for each stage of its construction
5. Also print the struct when it is complete (i.e. when the form submit button has been clicked)

### Part 2

6. Write a simple hashing function (your implementation - either of
   your own design or a known algorithm), that given a string will
   calculate a hash of that string.  We are not looking for you to
   wrap a standard function, but to provide the implementation itself.
7. Use that function to calculate the hash of the `WebSiteurl` field
   and print the hash, and print out the hash once calculated.

### Go Struct
```go
type Data struct {
	WebsiteUrl         string
	SessionId          string
	ResizeFrom         Dimension
	ResizeTo           Dimension
	CopyAndPaste       map[string]bool // map[fieldId]true
	FormCompletionTime int // Seconds
}

type Dimension struct {
	Width  string
	Height string
}
```