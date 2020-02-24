package jstime

import (
	"sync"
	"testing"
	"time"
)

func TestSetInterval(t *testing.T) {
	var (
		timer1  *time.Timer
		timer2  *time.Timer
		counter int
		wg      *sync.WaitGroup
	)

	counter = 0
	wg = &sync.WaitGroup{}
	timer1 = time.NewTimer(2000 * time.Millisecond)
	timer2 = time.NewTimer(2500 * time.Millisecond)

	stop := SetInterval(func() {
		counter++
	}, 250, wg)

OuterLoop:
	for {
		select {
		case <-timer2.C:
			if counter != 8 {
				t.Errorf("Counter = %d, want - %d", counter, 8)
				stop()
			}
			break OuterLoop
		case <-timer1.C:
			stop()
		default:
		}
	}
	wg.Wait()
}

func TestSetTimeout(t *testing.T) {
	var (
		timer1  *time.Timer
		timer2  *time.Timer
		timer3  *time.Timer
		counter int
		wg      *sync.WaitGroup
	)

	counter = 0
	wg = &sync.WaitGroup{}
	timer1 = time.NewTimer(1000 * time.Millisecond)
	timer2 = time.NewTimer(2500 * time.Millisecond)
	timer3 = time.NewTimer(3000 * time.Millisecond)

	stopI := SetInterval(func() {
		counter++
	}, 250, wg)

	stopT := SetTimeout(func() {
		stopI()
	}, 2000, wg)

	_ = SetTimeout(func() {
		stopT()
	}, 1500, wg)

OuterLoop:
	for {
		select {
		case <-timer1.C:
			stopT()
		case <-timer2.C:
			stopI()
		case <-timer3.C:
			if counter != 10 {
				t.Errorf("Counter = %d, want - %d", counter, 8)
				stopI()
			}
			break OuterLoop
		default:
		}
	}
	wg.Wait()
}

func Test_jsFields_closeChannel(t *testing.T) {
	type fields struct {
		millis   uint32
		do       func(*jsFields)
		secondDo func(chan struct{})
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "jsFields.closeChannel() case 1",
			fields: fields{
				millis:   1,
				do:       func(jsF *jsFields) {},
				secondDo: func(ch chan struct{}) {},
			},
			want: false,
		},
		{
			name: "jsFields.closeChannel() case 2",
			fields: fields{
				millis: 1,
				do: func(jsF *jsFields) {
					jsF.closeChannel()
				},
				secondDo: func(ch chan struct{}) {},
			},
			want: true,
		},
		{
			name: "jsFields.closeChannel() case 3",
			fields: fields{
				millis: 1,
				do: func(jsF *jsFields) {
					jsF.stopExecute()
				},
				secondDo: func(ch chan struct{}) {
					<-ch
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsF := newJSFields(tt.fields.millis)
			go tt.fields.secondDo(jsF.stopCh)
			tt.fields.do(jsF)
			if jsF.isChannelClosed() != tt.want {
				t.Errorf("closeChannel() : isClosed = %v, want - %v", jsF.isClosed, tt.want)
			}
		})
	}
}

func Test_jsFields_stopExecute(t *testing.T) {
	type fields struct {
		millis   uint32
		do       func(*jsFields)
		secondDo func(chan struct{})
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "jsFields.stopExecute() case 1",
			fields: fields{
				millis:   1,
				do:       func(jsF *jsFields) {},
				secondDo: func(ch chan struct{}) {},
			},
			want: false,
		},
		{
			name: "jsFields.stopExecute() case 2",
			fields: fields{
				millis: 1,
				do: func(jsF *jsFields) {
					jsF.closeChannel()
				},
				secondDo: func(ch chan struct{}) {},
			},
			want: true,
		},
		{
			name: "jsFields.stopExecute() case 3",
			fields: fields{
				millis: 1,
				do: func(jsF *jsFields) {
					jsF.stopExecute()
				},
				secondDo: func(ch chan struct{}) {
					<-ch
				},
			},
			want: true,
		},
		{
			name: "jsFields.stopExecute() case 4",
			fields: fields{
				millis: 1,
				do: func(jsF *jsFields) {
					jsF.closeChannel()
					jsF.stopExecute()
				},
				secondDo: func(ch chan struct{}) {},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsF := newJSFields(tt.fields.millis)
			go tt.fields.secondDo(jsF.stopCh)
			tt.fields.do(jsF)
			if jsF.isChannelClosed() != tt.want {
				t.Errorf("stopExecute() : isClosed = %v, want - %v", jsF.isClosed, tt.want)
			}
		})
	}
}
