package main

type ResponseBodyGenerator interface {
	Generate() string
}

type MaxResponseBodyGenerator struct {
	maxLength int
}

func (instance *MaxResponseBodyGenerator) Generate() string {
	var returnArray = make([]rune, instance.maxLength)
	for i := range returnArray {
		returnArray[i] = '-'
	}
	return string(returnArray)
}

func NewMaxResponseBodyGenerator(maxLength int) *MaxResponseBodyGenerator {
	return &MaxResponseBodyGenerator{maxLength}
}

type RandomResponseBodyGenerator struct {
	minLength int
	maxLength int
	random    Random
}

func (instance *RandomResponseBodyGenerator) Generate() string {
	randValue := instance.random.Int(instance.minLength, instance.maxLength)
	var returnArray = make([]rune, randValue)
	for i := range returnArray {
		returnArray[i] = '-'
	}
	return string(returnArray)
}

func NewRandomResponseBodyGenerator(minLength int, maxLength int) *RandomResponseBodyGenerator {
	random := NewRealRandom()
	return &RandomResponseBodyGenerator{minLength, maxLength, random}
}

type FakeResponseBodyGenerator struct {
	use string
}

func (instance *FakeResponseBodyGenerator) Generate() string {
	return instance.use
}

func (instance *FakeResponseBodyGenerator) UseString(value string) {
	instance.use = value
}

func NewFakeResponseBodyGenerator() *FakeResponseBodyGenerator {
	return &FakeResponseBodyGenerator{""}
}
