package pkg

import (
	"testing"
)

func TestAdd8(t *testing.T) {
	keys := SortKeys(testCaseByteMathMap)
	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			tc := testCaseByteMathMap[key]
			sum, carryOut := Add8(tc.a, tc.b, tc.carryIn)
			if sum != tc.sum || carryOut != tc.carryOut {
				t.Errorf("Add8(%d, %d, %d) =\n%d, %d; want\n%d, %d", tc.a, tc.b, tc.carryIn, sum, carryOut, tc.sum, tc.carryOut)
			}
		})
	}
}

func TestAdd8Slice(t *testing.T) {
	keys := SortKeys(testCaseAdd8SliceMap)
	for _, key := range keys {
		t.Run(key, func(t *testing.T) {
			tc := testCaseAdd8SliceMap[key]
			sum, carryOut := Add8BigEndian(tc.a, tc.b, tc.carryIn)
			if !SliceEqual(sum, tc.sum) || carryOut != tc.carryOut {
				t.Errorf("Add8BigEndian(%v, %v, %d) =\n%v, %d; want\n%v, %d", tc.a, tc.b, tc.carryIn, sum, carryOut, tc.sum, tc.carryOut)
			}
		})
	}
}
