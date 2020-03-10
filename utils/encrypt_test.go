package utils

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEncryptAES(t *testing.T) {
	txt := "This is nice"
	enc := AesEncrypt([]byte(txt), Key)
	bb := make([]byte, len(enc)+2)
	binary.BigEndian.PutUint16(bb, uint16(len(enc)+2))
	copy(bb[2:], enc)
	fmt.Printf("encrypted:%x,%x,%x\n", []byte(txt), enc, bb)
}
func TestWriteTo(t *testing.T) {
	str := `The MD5 message-digest algorithm is a widely used hash function producing a 128-bit hash value. Although MD5 was initially designed to be used as a cryptographic hash function, it has been found to suffer from extensive vulnerabilities. It can still be used as a checksum to verify data integrity, but only against unintentional corruption. It remains suitable for other non-cryptographic purposes, for example for determining the partition for a particular key in a partitioned database.[3]

MD5 was designed by Ronald Rivest in 1991 to replace an earlier hash function MD4,[4] and was specified in 1992 as RFC 1321.

One basic requirement of any cryptographic hash function is that it should be computationally infeasible to find two distinct messages that hash to the same value. MD5 fails this requirement catastrophically; such collisions can be found in seconds on an ordinary home computer.

The weaknesses of MD5 have been exploited in the field, most infamously by the Flame malware in 2012. The CMU Software Engineering Institute considers MD5 essentially "cryptographically broken and unsuitable for further use".[5]
`
	assert.True(t, len(str) <= 2048, "from str length should<=2048 byte,but it is "+fmt.Sprintf("%d", len(str)))
	from := bytes.NewBufferString(str)
	middle := bytes.NewBuffer(nil)
	if err := WriteTo(from, middle, false, true); err != nil {
		t.Log(err)
	}

	enc := AesEncrypt([]byte(str), Key)
	fmt.Printf("from:%x\n", enc)
	fmt.Printf("middle:%x\n", middle.Bytes())
	assert.True(t, bytes.Equal(middle.Bytes()[2:], enc), "encrypt error")
	target := bytes.NewBuffer(nil)
	if err := WriteTo(middle, target, true, false); err != nil {
		t.Log(err)
	}
	fmt.Printf("target:%s\n", target.Bytes())
	assert.True(t, string(target.Bytes()) == str, "decrypt error")
}

func BenchmarkEncryptAES(b *testing.B) {
	key := "S12@49I789AOqdef"
	txt := "This is nice"
	enc := AesEncrypt([]byte(txt), key)
	fmt.Printf("encrypted:%x\n", enc)
	dec := AesDecrypt(enc, key)
	fmt.Println("dec:", string(dec))
}
func TestBinary(t *testing.T) {
	by := make([]byte, 2)
	binary.BigEndian.PutUint16(by, uint16(18))
	fmt.Printf("%x\n", by)
}
