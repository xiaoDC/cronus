package timewheel

import (
	"reflect"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *TimeWheel
	}{
		{
			name: "first",
			want: &TimeWheel{
				interval:   time.Millisecond * 50,
				ticker:     time.NewTicker(time.Millisecond * 50),
				wheelSize:  60,
				currentPos: 0,
				startTime:  time.Now(),

				buckets:   []*Bucket{},
				TaskQueue: make(chan string, 1000),
				exitC:     make(chan struct{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

// func TestNewTimeWheel(t *testing.T) {
// 	type args struct {
// 		interval  time.Duration
// 		wheelSize int64
// 	}
// 	tests := []struct {
// 		name string
// 		args args
// 		want *TimeWheel
// 	}{
// 		// TODO: Add test cases.
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			if got := NewTimeWheel(tt.args.interval, tt.args.wheelSize); !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("NewTimeWheel() = %v, want %v", got, tt.want)
// 			}
// 		})
// 	}
// }
