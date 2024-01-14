# ccid_go

CCID Go lang implementation

[Reference](https://github.com/Pencroff/ccid)


## How to install

    go get -u github.com/Pencroff/ccid_go

## How to use

```go
package main

import (
	"fmt"
	c "github.com/Pencroff/ccid_go"
	e "github.com/Pencroff/ccid_go/extras"
	p "github.com/Pencroff/ccid_go/pkg"
)

func main() {
	// Create new random bytes reader
	r, err := e.NewHybridRandReaderWithSize(e.SIZE_32k)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Create new strategy for monotonic mutations
	s := p.NewFiftyPercentMonotonicStrategy(r)
	// Create new CcId generator
	ccidGen, err := c.NewMonotonicCcIdGen(p.ByteSliceSize96, r, s)
	if err != nil {
		fmt.Println(err)
		return
	}
	// Generate new CcId
	ccid, err := ccidGen.Next()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", ccid)
	// CcId{size: 12, timestamp: 305250750 (2024-01-14T16:45:50Z), payload: 0x9461896a128ac957}
}
```


## How to increase version

* commit all required changes
* `git tag <version - v0.0.2>`
* `git push origin --tags`
* done - check docs on pkg.go.dev

