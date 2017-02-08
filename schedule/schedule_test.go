package schedule

import (
	"reflect"
	"testing"
	"time"
)

const (
	EmptyScheduleText = `
{
	"Users": ["a", "b", "c"],
	"Start": "2017-02-01T10:00:00Z",
	"RotationLength": "168h",
	"ScheduleFor": "504h"
}`
	FilledScheduleText = `
{
	"Users": ["b", "c", "a"],
	"Start": "2017-03-01T10:00:00Z",
	"RotationLength": "168h",
	"ScheduleFor": "504h",
	"Rotations": [
		{
			"Start": "2017-02-01T10:00:00Z",
			"Primary": "a",
			"Secondary": "b"
		},
		{
			"Start": "2017-02-08T10:00:00Z",
			"Primary": "b",
			"Secondary": "c"
		},
		{
			"Start": "2017-02-15T10:00:00Z",
			"Primary": "c",
			"Secondary": "a"
		},
		{
			"Start": "2017-02-22T10:00:00Z",
			"Primary": "a",
			"Secondary": "b"
		}
	]
}`
)

var (
	Start = time.Date(2017, time.February, 1, 10, 0, 0, 0, time.UTC)
)

func EmptySchedule() *Schedule {
	return &Schedule{
		Users: []string{"a", "b", "c"},
		Start: Start,
		RotationLength: "168h",
		rotationLength: 7 * 24 * time.Hour,
		ScheduleFor: "504h",
		scheduleFor: 3 * 7 * 24 * time.Hour,
	}
}

func FilledSchedule() *Schedule {
	return &Schedule{
		Users: []string{"b", "c", "a"},
		Start: time.Date(2017, time.March, 1, 10, 0, 0, 0, time.UTC),
		RotationLength: "168h",
		rotationLength: 7 * 24 * time.Hour,
		ScheduleFor: "504h",
		scheduleFor: 3 * 7 * 24 * time.Hour,
		Rotations: []Rotation{
			{
				Start: Start,
				Primary: "a",
				Secondary: "b",
			},
			{
				Start: time.Date(2017, time.February, 8, 10, 0, 0, 0, time.UTC),
				Primary: "b",
				Secondary: "c",
			},
			{
				Start: time.Date(2017, time.February, 15, 10, 0, 0, 0, time.UTC),
				Primary: "c",
				Secondary: "a",
			},
			{
				Start: time.Date(2017, time.February, 22, 10, 0, 0, 0, time.UTC),
				Primary: "a",
				Secondary: "b",
			},
		},
	}
}

func TestParseEmptySchedule(t *testing.T) {
	s, err := NewSchedule([]byte(EmptyScheduleText))
	if err != nil {
		t.Error(err)
	}
	expected := EmptySchedule()
	if !reflect.DeepEqual(expected, s) {
		t.Errorf("parsed schedule does not match expected\nExpected:\n%+v\n---\nGot:\n%+v\n", expected, s)
	}
}

func TestParseFilledSchedule(t *testing.T) {
	s, err := NewSchedule([]byte(FilledScheduleText))
	if err != nil {
		t.Error(err)
	}
	expected := FilledSchedule()
	if !reflect.DeepEqual(expected, s) {
		t.Errorf("parsed schedule does not match expected\nExpected:\n%+v\n---\nGot:\n%+v\n", expected, s)
	}
}

func TestGenerateFromEmpty(t *testing.T) {
	empty := EmptySchedule()
	empty.now = Start
	filled := FilledSchedule()
	filled.now = empty.now
	s, err := empty.Generate()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(filled, s) {
		t.Errorf("generated schedule does not match expected\nExpected:\n%+v\n---\nGot:\n%+v\n", filled, s)
	}
}

func TestGenerateIsIdempotent(t *testing.T) {
	filled := FilledSchedule()
	filled.now = Start
	s, err := filled.Generate()
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(filled, s) {
		t.Errorf("generated schedule does not match expected\nExpected:\n%+v\n---\nGot:\n%+v\n", filled, s)
	}
}
