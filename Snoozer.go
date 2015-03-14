package main

import (
	"time"
)

type Snoozer interface {
	RandomSnooze()
}

type RealSnoozer struct {
	Min    time.Duration
	Max    time.Duration
	random Random
}

func (instance *RealSnoozer) RandomSnooze() {
	randomSleep := instance.random.Duration(instance.Min, instance.Max)
	time.Sleep(randomSleep)
}

func NewRealSnoozer(min time.Duration, max time.Duration) *RealSnoozer {
	return &RealSnoozer{min, max, &RealRandom{}}
}

type FakeSnoozer struct {
	duration time.Duration
}

func (instance *FakeSnoozer) RandomSnooze() {
	time.Sleep(instance.duration)
}

func (instance *FakeSnoozer) SleepFor(duration time.Duration) {
	instance.duration = duration
}

func NewFakeSnoozer() *FakeSnoozer {
	return &FakeSnoozer{0}
}
