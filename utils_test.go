package utils

import (
	"testing"
	"time"
)

func TestReverseList(t *testing.T) {
	intArray := []interface{}{1, 2, 3, 4, 5}
	stringArray := []interface{}{"hello", "how", "are", "you"}

	intList := ReverseList(intArray)
	if intList[0] == 5 {
		t.Logf("It works, here's the list: %v", intList)
	} else {
		t.Errorf("It didn't work, here's the list: %v", intList)
	}

	stringList := ReverseList(stringArray)
	if stringList[0] == "you" {
		t.Logf("It works, here's the list: %v", stringList)
	} else {
		t.Errorf("It didn't work, here's the list: %v", stringList)
	}
}

func TestGetMonthFromName(t *testing.T) {
	month := time.January
	secondMonth := time.February

	january, _ := GetMonthFromName("January")
	february, _ := GetMonthFromName("February")

	if january == month {
		t.Logf("Passed, returned %v", january)
	} else {
		t.Errorf("Didn't work, returned: %v", january)
	}
	if february == secondMonth {
		t.Logf("Passed, returned %v", february)
	} else {
		t.Errorf("Didn't work, returned: %v", february)
	}
}

// Can conduct more tests..
