package main

import (
	"math/rand"
	"time"
)

type Random interface {
	Int(from int, to int) (randomInt int)
	Duration(from time.Duration, to time.Duration) time.Duration
}

type RealRandom struct{}

func (instance *RealRandom) Int(min int, max int) (randomInt int) {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func (instance *RealRandom) Duration(min time.Duration, max time.Duration) time.Duration {
	rand.Seed(time.Now().Unix())
	return time.Duration(rand.Int63n(int64(max-min) + int64(min)))
}

func NewRealRandom() *RealRandom {
	return &RealRandom{}
}

type FakeRandom struct {
	number   int
	duration time.Duration
}

func (instance *FakeRandom) Int(from int, to int) (randomInt int) {
	return instance.number
}

func (instance *FakeRandom) Duration(from time.Duration, to time.Duration) time.Duration {
	return instance.duration
}

func (instance *FakeRandom) ForIntUse(value int) {
	instance.number = value
}

func (instance *FakeRandom) ForDurationUse(value time.Duration) {
	instance.duration = value
}

func NewFakeRandom() (random *FakeRandom) {
	return &FakeRandom{0, 0}
}
