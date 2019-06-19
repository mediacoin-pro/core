package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

//var allTestMode = len(os.Args) == 2 && os.Args[1] == "-test.v=true"

func TestGeneratePrivateKey(t *testing.T) {

	prv1 := NewPrivateKey()
	prv2 := NewPrivateKey()
	pub1 := prv1.PublicKey()

	assert.Equal(t, 33, len(prv1.Bytes()))
	assert.NotEqual(t, prv1, prv2)
	assert.Equal(t, PublicKeySize, len(pub1.Bytes()))
}

func TestGeneratePrivateKeyByPassword(t *testing.T) {
	//if allTestMode { // skip test
	//	return
	//}
	password := "SuperPuperSecret-10003"

	prv1 := NewPrivateKeyBySecret(password)
	pub1 := prv1.PublicKey()

	prv2 := NewPrivateKeyBySecret(password)
	pub2 := prv2.PublicKey()

	prv3 := NewPrivateKeyBySecret(password + ".")
	pub3 := prv3.PublicKey()

	assert.Equal(t, "AZsjpgKFOsdnnBTPxdXtu-EXdcMyaIzXqwynAQ47SQr1", prv1.String())
	assert.Equal(t, "0x04c093844a25ee795885be70be20d445d19ada4c4f1369f9730ed4d827a1b5d35481276a7fb9fcae6da867c106ad81faeb67f7b4ab45954612f4715e98e6abea96", pub1.String())
	assert.Equal(t, 33, len(prv1.Bytes()))
	assert.Equal(t, prv1.Bytes(), prv2.Bytes())
	assert.Equal(t, pub1.Bytes(), pub2.Bytes())
	assert.NotEqual(t, prv1.Bytes(), prv3.Bytes())
	assert.NotEqual(t, pub1.Bytes(), pub3.Bytes())
}

func TestPrivateKey_D(t *testing.T) {
	org := NewPrivateKey()
	buf := org.Bytes()

	decPrv, _ := decodePrivateKey(buf)

	assert.Equal(t, org, decPrv)
}

func TestSign(t *testing.T) {
	prv := NewPrivateKey()
	data := []byte("Airbus откажется от самолетов A380")

	sign1 := prv.Sign(data)
	sign2 := prv.Sign(data)

	assert.Equal(t, PublicKeySize, len(sign1))
	assert.Equal(t, PublicKeySize, len(sign2))
	assert.False(t, bytes.Equal(sign1, sign2))
}

func TestVerify(t *testing.T) {
	prv := NewPrivateKey()
	pub := prv.PublicKey()
	data := []byte("Совет по туризму Норвегии определил места обитания редких покемонов")
	sign := prv.Sign(data)

	verify1 := pub.Verify(data, sign)
	verify2 := pub.Verify(data, sign)

	assert.True(t, verify1)
	assert.True(t, verify2)
}

func TestVerifyFail(t *testing.T) {
	prv := NewPrivateKey()
	pub := prv.PublicKey()
	data := []byte("Роскачество не нашло хороших пельменей на российском рынке")

	sign := prv.Sign(data)
	sign[0]++
	verify := pub.Verify(data, sign)

	assert.False(t, verify)
}
