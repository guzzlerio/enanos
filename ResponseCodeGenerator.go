package main

type ResponseCodeGenerator interface {
	GenerateServerErrorCode() int
	GenerateRedirectionCode() int
	GenerateClientErrorCode() int
}

type RandomResponseCodeGenerator struct {
	responseCodes_3XX []int
	responseCodes_4XX []int
	responseCodes_5XX []int
	randomGen         Random
}

func (instance *RandomResponseCodeGenerator) GenerateServerErrorCode() int {
	from := 0
	to := len(instance.responseCodes_5XX)
	index := instance.randomGen.Int(from, to)
	return instance.responseCodes_5XX[index]
}

func (instance *RandomResponseCodeGenerator) GenerateRedirectionCode() int {
	from := 0
	to := len(instance.responseCodes_3XX)
	index := instance.randomGen.Int(from, to)
	return instance.responseCodes_3XX[index]
}

func (instance *RandomResponseCodeGenerator) GenerateClientErrorCode() int {
	from := 0
	to := len(instance.responseCodes_4XX)
	index := instance.randomGen.Int(from, to)
	return instance.responseCodes_4XX[index]
}

func NewRandomResponseCodeGenerator(responseCodes_3XX []int, responseCodes_4XX []int, responseCodes_5XX []int) *RandomResponseCodeGenerator {
	return &RandomResponseCodeGenerator{responseCodes_3XX, responseCodes_4XX, responseCodes_5XX, NewRealRandom()}
}

type FakeResponseCodeGenerator struct {
	use int
}

func (instance *FakeResponseCodeGenerator) Use(value int) {
	instance.use = value
}

func (instance *FakeResponseCodeGenerator) GenerateServerErrorCode() int {
	return instance.use
}

func (instance *FakeResponseCodeGenerator) GenerateRedirectionCode() int {
	return instance.use
}

func (instance *FakeResponseCodeGenerator) GenerateClientErrorCode() int {
	return instance.use
}

func NewFakeResponseCodeGenerator() *FakeResponseCodeGenerator {
	return &FakeResponseCodeGenerator{}
}
