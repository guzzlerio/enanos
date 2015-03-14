package main

type ResponseCodeGenerator interface {
	Generate() int
}

type RandomResponseCodeGenerator struct {
	responseCodes []int
	randomGen     Random
}

func (instance *RandomResponseCodeGenerator) Generate() int {
	from := 0
	to := len(instance.responseCodes)
	index := instance.randomGen.Int(from, to)
	return instance.responseCodes[index]
}

func NewRandomResponseCodeGenerator(responseCodes []int) *RandomResponseCodeGenerator {
	return &RandomResponseCodeGenerator{responseCodes, NewRealRandom()}
}

type FakeResponseCodeGenerator struct {
	use int
}

func (instance *FakeResponseCodeGenerator) Use(value int) {
	instance.use = value
}

func (instance *FakeResponseCodeGenerator) Generate() int {
	return instance.use
}

func NewFakeResponseCodeGenerator() *FakeResponseCodeGenerator {
	return &FakeResponseCodeGenerator{}
}
