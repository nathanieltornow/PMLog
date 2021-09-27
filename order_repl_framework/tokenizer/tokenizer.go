package tokenizer

type Tokenizer struct {
	colorCtr map[uint32]uint32
}

func NewTokenizer() *Tokenizer {
	t := new(Tokenizer)
	t.colorCtr = make(map[uint32]uint32)
	return t
}

func (t *Tokenizer) NewToken(color)
