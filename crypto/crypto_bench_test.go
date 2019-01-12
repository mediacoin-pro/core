package crypto

import "testing"

/**
	PASS
	BenchmarkGenerateKeyBySecret-4	       1	1928276011 ns/op	55095848 B/op	 1238393 allocs/op
	BenchmarkGenerateKey-4        	   50000	     25783 ns/op	     976 B/op	      18 allocs/op
	BenchmarkSign-4               	   20000	     77236 ns/op	    3328 B/op	      52 allocs/op
	BenchmarkVerify-4             	    5000	    291804 ns/op	   50546 B/op	    1107 allocs/op
	ok  	xnet/crypto	7.305s
**/

func BenchmarkGenerateKeyBySecret(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPrivateKeyBySecret("secret-string")
	}
}

func BenchmarkGenerateKey(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewPrivateKey()
	}
}

func BenchmarkSign(b *testing.B) {
	prv := NewPrivateKey()
	data := []byte("Abc Ёпрст")

	for i := 0; i < b.N; i++ {
		prv.Sign(data)
	}
}

func BenchmarkVerify(b *testing.B) {
	prv := NewPrivateKey()
	pub := prv.PublicKey()
	data := []byte("Abc Ёпрст")
	sign := prv.Sign(data)

	for i := 0; i < b.N; i++ {
		if !pub.Verify(data, sign) {
			b.Fatal("Verify fail")
		}
	}
}
