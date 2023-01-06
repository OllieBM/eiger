package rolling_hash

import "math"

// A Prime random large prime number
const prime = 1190494759 // selected from http://www.primos.mat.br/2T_en.html
const base = 101

// computeHash returns the hash value of string s
func compute_hash(s string) uint64 {

	hashValue := uint64(0)
	posPow := uint64(1)
	var a rune
	a = rune('a')
	for _, c := range s {
		// this character as an uint between 1...27;
		v := uint64(c - a + 1)
		hashValue = (hashValue + (v)*posPow) % prime
		posPow = (posPow * base) % prime
	}
	return hashValue
}

// RollingHash checks if previous string and current string are different
func rolling_hash_value(prev string, c rune) uint64 {
	hash := compute_hash(prev)
	nextHash := ((hash-math.Pow(prev[0], len(prev()-1))*base + int(c)) % prime
	return nextHash	
}

/*
// computes the hash value of the input string s
long long compute_hash(string s) {
    const int base = 31;   // base
    const int prime = 1e9 + 9; // large prime number
    long long hashValue = 0;
    long long posPow = 1;
    for (char c : s) {
        hashValue = (hashValue + (c - 'a' + 1) * posPow) % prime;
        posPow = (posPow * base) % prime;
    }
    return hashValue;
}
// finds the hash value of next substring given nxt as the ending character
// and the previous substring prev
long long rolling_hash(string prev,char nxt)
{
   const int p = 31;
   const int m = 1e9 + 9;
   long long H=compute_hash(prev);
   long long Hnxt=( ( H - pow(prev[0],prev.length()-1) ) * p + (int)nxt ) % m;
   return Hnxt;
}

*/
