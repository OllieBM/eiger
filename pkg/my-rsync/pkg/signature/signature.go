package signature

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"io"

	"github.com/rs/zerolog/log"
)

type signature struct {
	chunkSize uint64
	strong    [][]byte // strong checksums based on a hash
	weak      []uint64 // rolling checksum
	// TODO can be extended to include
	// version string
	// base  uint8
	// prime uint64
	// so different applications can communicate
}

type SignatureWriter struct {
	r bufio.Reader
}

func CreateSignature(in io.Reader, chunkLen uint) *signature {

	sig := &signature{

		chunkSize: uint64(chunkLen),
		strong:    make([][]byte, 0, 16),
		weak:      make([]uint64, 0, 16), // rolling checksum
	}

	strongHasher := md5.New()

	buf := make([]byte, chunkLen)
	for {
		n, err := in.Read(buf)

		if err != nil {
			if err != io.EOF {
				log.Error().Err(err).Msg("reading file")
			}
			break
		}
		fmt.Println(buf)
		if n != 0 {
			// strong hash

			buf = buf[:n] // trim off garbage
			// create strong hash
			strongHasher.Reset()
			strongHasher.Write(buf)
			sh := strongHasher.Sum(nil)
			// create rolling hash
			weak := rollingHash(buf)
			sig.strong = append(sig.strong, sh)
		}

	}
	return sig
}

// func CreateSignature(in io.Reader, ChunkSize uint) {

// 	buf := make([]byte, 0, ChunkSize)
// 	r := bufio.NewReader(in)
// 	for {
// 		n, err := r.Read(buf)
// 		buf = buf[:n]
// 		if n == 0 {
// 			if err == nil {
// 				continue
// 			}
// 			if err == io.EOF {
// 				break
// 			}
// 			log.Fatal(err)
// 		}
// 		if n != 0 {
// 			nChunks++
// 			nBytes += int64(len(buf))

// 		}
// 		// process buf
// 		if err != nil && err != io.EOF {
// 			log.Fatal(err)
// 		}
// 	}
// }

// ChunkSize is the number of characters in
func ReadAsChunks(in io.Reader, ChunkLen uint) [][]byte {

	chunks := make([][]byte, 0)
	buf := make([]byte, ChunkLen)
	for {
		n, err := in.Read(buf)

		if err != nil {
			if err != io.EOF {
				log.Error().Err(err).Msg("reading file")
			}
			break
		}
		fmt.Println(buf)
		if n != 0 {
			buf = buf[:n] // trim off garbage
			chunks = append(chunks, buf)
		}

	}
	return chunks
	// buf := make([]byte, 0, ChunkLen)
	// //r := bufio.NewReader(in)
	// count := 0
	// byteCount := 0
	// chunks := make([][]byte, 0, 8)
	// for {
	// 	fmt.Println("loop")
	// 	n, err := in.Read(buf)
	// 	buf = buf[:n] // trim off garbage
	// 	if n != 0 {
	// 		count++
	// 		byteCount += (n)
	// 		fmt.Println(buf)
	// 		chunks = append(chunks, buf)
	// 	}
	// 	if err != nil {
	// 		if err != io.EOF {
	// 			log.Fatal(err)
	// 		}
	// 		return chunks
	// 	}
	// }
	return chunks
}
