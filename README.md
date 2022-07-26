# goWeightGurus

## Usage
```
package main

import (
	"github.com/Erich-Reitz/goWeightGurus"
)

func main() {
	entries := WeightGurus.GetNonDeletedEntries("email", "password") 
	// or
	WeightGurus.WriteNonDeletedEntriesToFile("email", "password", "results.json")
}
```