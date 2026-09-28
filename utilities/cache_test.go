package utilities

import (
	"strconv"
	"testing"
	"time"
)

func TestEntriesGetAdded(t *testing.T) {
	duration, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("Couldn't init duration")
	}

	cache := NewCache(duration)

	for i := range 100 {
		cache.Add(strconv.Itoa(i), nil)
	}

	if len(cache.entries) != 100 {
		t.Errorf("Cache len should be 100.")
	}
}

func TestEntryValuesArePreserved(t *testing.T) {
	duration, err := time.ParseDuration("1h")
	if err != nil {
		t.Errorf("Couldn't init duration")
	}

	cache := NewCache(duration)

	cache.Add("1", "111")
	cache.Add("2", "2")
	cache.Add("3", "3")

	expected := "111"
	value, ok := cache.Get("1")
	if !ok {
		t.Errorf("Couldn't retrieve value for key %v", "1")
	}

	if value.(string) != expected {
		t.Errorf("Retrieved value for key %v should be %v", "1", expected)
	}
}

func TestEntryValuesGetCleanedUp(t *testing.T) {
	// Setup
	cacheCleanupDuration, err := time.ParseDuration("1ns")
	if err != nil {
		t.Errorf("Couldn't init duration")
	}

	waitDuration, err := time.ParseDuration("1s")
	if err != nil {
		t.Errorf("Couldn't init duration")
	}

	cache := NewCache(cacheCleanupDuration)

	// Action
	cache.Add("1", "111")
	time.Sleep(waitDuration)

	// Validation
	if len(cache.entries) != 0 {
		t.Errorf("Cache len should be 0.")
	}
}

func TestEntryValuesWontGetCleanedUp(t *testing.T) {
	// Setup
	cacheCleanupDuration, err := time.ParseDuration("5s")
	if err != nil {
		t.Errorf("Couldn't init duration")
	}

	waitDuration, err := time.ParseDuration("1ns")
	if err != nil {
		t.Errorf("Couldn't init duration")
	}

	cache := NewCache(cacheCleanupDuration)

	// Action
	cache.Add("1", "111")
	time.Sleep(waitDuration)

	// Validation
	if len(cache.entries) != 1 {
		t.Errorf("Cache len should be 1. We shouldn't do cleanup.")
	}
}
