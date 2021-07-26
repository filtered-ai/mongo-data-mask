package generators

type boolgen struct {
	Base
}

func NewBoolGenerator(base Base) *boolgen {
	return &boolgen{Base: base}
}

func (b *boolgen) Bool() bool {
	if b.randomByte() == 0 {
		return true
	}
	return false
}

func (b *boolgen) randomByte() byte {
	return byte(b.Pcg32.Random() & 0x01)
}
