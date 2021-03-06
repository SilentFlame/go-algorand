// Copyright (C) 2019-2020 Algorand, Inc.
// This file is part of go-algorand
//
// go-algorand is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as
// published by the Free Software Foundation, either version 3 of the
// License, or (at your option) any later version.
//
// go-algorand is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with go-algorand.  If not, see <https://www.gnu.org/licenses/>.

package compactcert

import (
	"testing"

	"github.com/algorand/go-algorand/crypto"
)

func TestHashCoin(t *testing.T) {
	var slots [32]uint64
	var sigcom [32]byte

	crypto.RandBytes(sigcom[:])

	for j := uint64(0); j < 1000; j++ {
		coin := hashCoin(j, sigcom, uint64(len(slots)))
		if coin >= uint64(len(slots)) {
			t.Errorf("hashCoin out of bounds")
		}

		slots[coin]++
	}

	for i, count := range slots {
		if count < 3 {
			t.Errorf("slot %d too low: %d", i, count)
		}
		if count > 100 {
			t.Errorf("slot %d too high: %d", i, count)
		}
	}
}

func BenchmarkHashCoin(b *testing.B) {
	var sigcom [32]byte
	crypto.RandBytes(sigcom[:])

	for i := 0; i < b.N; i++ {
		hashCoin(uint64(i), sigcom, 1024)
	}
}

func TestNumReveals(t *testing.T) {
	billion := uint64(1000 * 1000 * 1000)
	microalgo := uint64(1000 * 1000)
	provenWeight := 2 * billion * microalgo
	secKQ := uint64(128)
	bound := uint64(1000)

	for i := uint64(3); i < 10; i++ {
		signedWeight := i * billion * microalgo
		n, err := numReveals(signedWeight, provenWeight, secKQ, bound)
		if err != nil {
			t.Error(err)
		}

		if n < 50 || n > 300 {
			t.Errorf("numReveals(%d, %d, %d) = %d looks suspect",
				signedWeight, provenWeight, secKQ, n)
		}
	}
}

func BenchmarkNumReveals(b *testing.B) {
	billion := uint64(1000 * 1000 * 1000)
	microalgo := uint64(1000 * 1000)
	provenWeight := 100 * billion * microalgo
	signedWeight := 110 * billion * microalgo
	secKQ := uint64(128)
	bound := uint64(1000)

	nr, err := numReveals(signedWeight, provenWeight, secKQ, bound)
	if nr < 900 {
		b.Errorf("numReveals(%d, %d, %d) = %d < 900", signedWeight, provenWeight, secKQ, nr)
	}

	for i := 0; i < b.N; i++ {
		_, err = numReveals(signedWeight, provenWeight, secKQ, bound)
		if err != nil {
			b.Error(err)
		}
	}
}
