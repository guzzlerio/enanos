package main

import (
	"time"
)

type Snoozer interface {
	Snooze()
}

type MaxSnoozer struct {
	Max time.Duration
}

func (instance *MaxSnoozer) Snooze() {
	time.Sleep(instance.Max)
}

func NewMaxSnoozer(max time.Duration) *MaxSnoozer {
	return &MaxSnoozer{max}
}

type RandomSnoozer struct {
	Min    time.Duration
	Max    time.Duration
	random Random
}

func (instance *RandomSnoozer) Snooze() {
	randomSleep := instance.random.Duration(instance.Min, instance.Max)
	time.Sleep(randomSleep)
}

func NewRandomSnoozer(min time.Duration, max time.Duration) *RandomSnoozer {
	return &RandomSnoozer{min, max, &RealRandom{}}
}

type FakeSnoozer struct {
	duration time.Duration
}

func (instance *FakeSnoozer) Snooze() {
	time.Sleep(instance.duration)
}

func (instance *FakeSnoozer) SleepFor(duration time.Duration) {
	instance.duration = duration
}

func NewFakeSnoozer() *FakeSnoozer {
	return &FakeSnoozer{0}
}
