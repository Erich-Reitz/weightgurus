# weightgurus

## Usage
```
package main

import (
	"weightgurus/weightgurus"
)

func main() {
	entries := WeightGurus.GetNonDeletedEntries("email", "password") 
	// or
	WeightGurus.WriteNonDeletedEntriesToFile("email", "password", "results.json")
}
```
