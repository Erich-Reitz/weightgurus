# goWeightGurus

## Usage
```
package main

import (
	"weightgurus/weightgurus"
)

func main() {
	entries := WeightGurus.GetNonDeletedEntries("email", "password") 
	// or
	WeightGurus.WriteAllEntriesToFile("email", "password", "results.json")
}
```