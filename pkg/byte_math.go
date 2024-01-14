package pkg

// Add8 adds two bytes and a carryIn byte and returns the sum and carryOut byte.
// The carryOut byte is 1 if the sum overflows, otherwise 0.
// The sum is truncated to 8 bits.
// 'a' and 'b' are the bytes to add.
// 'carryIn' is the carry in byte. 0 or 1.
func Add8(a, b, carryIn byte) (sum, carryOut byte) {
	sum = a + b + carryIn
	carryOut = ((a & b) | ((a | b) &^ sum)) >> 7
	return
}

// Add8BigEndian adds two byte slices and a carryIn byte and returns the sum byte slice and carryOut byte.
// The carryOut byte is 1 if the sum overflows, otherwise 0.
// 'a' and 'b' are the byte slices to add.
// 'carryIn' is the carry in byte. 0 or 1.
func Add8BigEndian(a, b []byte, carryIn byte) (sum []byte, carryOut byte) {
	idxA := len(a)
	idxB := len(b)
	if idxA < idxB {
		a, b = b, a
		idxA, idxB = idxB, idxA
	}
	sum = make([]byte, idxA)
	if idxA == 0 && idxB == 0 {
		if carryIn != 0 {
			sum = append(sum, carryIn)
		}
		return
	}
	for idxA > 0 {
		idxA--
		idxB--
		vA := a[idxA]
		vB := byte(0)
		if idxB >= 0 {
			vB = b[idxB]
		}
		sum[idxA], carryOut = Add8(vA, vB, carryIn)
		carryIn = carryOut
	}
	return
}
