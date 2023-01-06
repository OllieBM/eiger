package rolling_hash

import (
	"fmt"
	"math"
)

// https://www.infoarena.ro/blog/rolling-hash

//const base = 16777619
const prime = uint64(1190494759) // selected from http://www.primos.mat.br/2T_en.html
const base = uint64(128)

const old = "abc"
const new = "bcd"

func hash(s string, prime uint64) uint64 {
	h := uint64(0)
	pos := uint64(1)
	//for i := 0; i < len(s); i++ {
	for _, c := range s {
		//h = (h*prime + uint64(s[i]))
		h = (h*base + uint64(c)) % prime
		pos = (pos * base) % prime
		//     rolling_hash = (rolling_hash * a + S[i]) % MOD
		//     an = (an * a) % MOD
	}
	return h
}

// func hash2(s string, prime uint32) uint64 {

// 	var h uint32
// 	var p_pow = uint64(1)
// 	for i := 0; i < len(s); i++ {
// 		h = (h*prime + uint32(s[i])) % prime
// 		p_pow = (p_pow * base) % prime
// 	}
// 	return h
// }

func rolling_hash(s string, prime uint64) []uint64 {

	n := uint64(3)
	result := make([]uint64, 0, 8)
	//hash = hash(s[:n])
	base := uint64(128) // every ascii character

	fmt.Println(s[:n])
	pos := uint64(1)
	rolling_hash := uint64(0)
	// initial hash window
	for _, c := range s[:n] {
		rolling_hash = (rolling_hash*base + uint64(c)) % prime
		pos = (pos * base) % prime
	}
	fmt.Println(rolling_hash)
	result = append(result, rolling_hash)
	// rolling hash
	lhs := 0
	for _, c := range s[n:] {
		// now we can compute hash based on a rolling window
		//prev := rolling_hash
		prev_c := uint64(s[lhs])

		rolling_hash = (rolling_hash*base + uint64(c) - pos*prev_c) % prime
		//rolling_hash = ((rolling_hash-pow(prev_c, n-1))*base + uint64(c)) % prime
		lhs++
		println(rolling_hash)
		result = append(result, rolling_hash)
	}
	return result
}

func pow(a, b uint64) uint64 {
	return uint64(math.Pow(float64(a), float64(b)))
}

// an = 1
//   rolling_hash = 0
//   for i in range(0, n):
//     rolling_hash = (rolling_hash * a + S[i]) % MOD
//     an = (an * a) % MOD
//   if rolling_hash == hash_p:
//     print "match"
//   for i in range(1, m - n + 1):
//     rolling_hash = (rolling_hash * a + S[i + n - 1] - an * S[i - 1]) % MOD
//     if rolling_hash == hash_p:
//         print "match"
