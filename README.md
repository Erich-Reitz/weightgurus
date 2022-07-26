# weightgurus

## Usage
```
package main

import (
	"github.com/Erich-Reitz/weightgurus"
)

func main() {
	entries := weightgurus.GetNonDeletedEntries("email", "password") 
	// or
	weightgurus.WriteNonDeletedEntriesToFile("email", "password", "results.json"); 
}
```
