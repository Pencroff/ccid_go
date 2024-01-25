package ccid_go

import (
	"bytes"
	"fmt"
	e "github.com/Pencroff/ccid_go/extras"
	p "github.com/Pencroff/ccid_go/pkg"
	"sync"
	"testing"
)

func TestCcIdGen_Race(t *testing.T) {

	const (
		numRoutines = 10
		numCycles   = 100
	)
	sizes := []byte{
		p.ByteSliceSize64,
		p.ByteSliceSize96,
		p.ByteSliceSize128,
		p.ByteSliceSize160,
	}
	for _, size := range sizes {
		t.Run(fmt.Sprintf("size_%d", size), func(t *testing.T) {
			r, _ := e.NewHybridRandReader()
			s := p.NewFiftyPercentMonotonicStrategy(r)
			notLockedGen, _ := NewMonotonicCcIdGen(size, r, s)
			gen := NewCcIdGenLocked(notLockedGen)
			prev, _ := gen.Next()
			queue := make(chan p.CcId)
			var wg sync.WaitGroup
			defer wg.Wait()
			wg.Add(numRoutines)
			for i := 0; i < numRoutines; i++ {
				go func(i int, g CcIdGen, q chan p.CcId) {
					defer wg.Done()
					for j := 0; j < numCycles; j++ {
						id, _ := g.Next()
						q <- id
					}
				}(i, gen, queue)
			}

			for i := 0; i < numRoutines*numCycles; i++ {
				next := <-queue
				if bytes.Compare(prev.Bytes(), next.Bytes()) >= 0 {
					fmt.Printf("%d - Not sequential in channel.\nPrev: %x\nmore then\nNext: %x\nDetails:\nPrev: %p - %#v\nNext: %p - %#v\n",
						i, prev.Bytes(), next.Bytes(), &prev, prev, &next, next)
				}
				prev = next
			}
		})
	}

}
