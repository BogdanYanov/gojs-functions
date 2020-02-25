package jstime

import (
	"sync"
	"testing"
	"time"
)

func TestSetInterval(t *testing.T) {
	var (
		timer1   = time.NewTimer(2000 * time.Millisecond)
		timer2   = time.NewTimer(2500 * time.Millisecond)
		interval = 250 * time.Millisecond
		counter  = 0
		wg       = &sync.WaitGroup{}
	)

	wg.Add(1)
	stop := SetInterval(func() {
		counter++
	}, interval)

OuterLoop:
	for {
		select {
		case <-timer2.C:
			if counter != 8 {
				t.Errorf("Counter = %d, want - %d", counter, 8)
				stop()
				wg.Done()
			}
			break OuterLoop
		case <-timer1.C:
			stop()
			wg.Done()
		}
	}
	wg.Wait()
}

func TestSetTimeout(t *testing.T) {
	var (
		timer1   = time.NewTimer(time.Second)
		timer2   = time.NewTimer(2500 * time.Millisecond)
		timer3   = time.NewTimer(3 * time.Second)
		interval = 250 * time.Millisecond
		timeout1 = 2 * time.Second
		timeout2 = 1500 * time.Millisecond
		counter  = 0
		wg       = &sync.WaitGroup{}
	)

	wg.Add(1)
	stopI := SetInterval(func() {
		counter++
	}, interval)

	wg.Add(1)
	stopT := SetTimeout(func() {
		stopI()
	}, timeout1)

	_ = SetTimeout(func() {
		stopT()
	}, timeout2)

OuterLoop:
	for {
		select {
		case <-timer1.C:
			stopT()
			wg.Done()
		case <-timer2.C:
			stopI()
			wg.Done()
		case <-timer3.C:
			if counter != 10 {
				t.Errorf("Counter = %d, want - %d", counter, 10)
				stopI()
				wg.Done()
			}
			break OuterLoop
		default:
		}
	}
	wg.Wait()
}

func Test_timeJSHelper_closeChannel(t *testing.T) {
	type fields struct {
		millis   uint32
		do       func(*timeJSHelper)
		secondDo func(chan struct{})
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "timeJSHelper.closeChannel() case 1",
			fields: fields{
				millis:   1,
				do:       func(helper *timeJSHelper) {},
				secondDo: func(ch chan struct{}) {},
			},
			want: false,
		},
		{
			name: "timeJSHelper.closeChannel() case 2",
			fields: fields{
				millis: 1,
				do: func(helper *timeJSHelper) {
					helper.closeChannel()
				},
				secondDo: func(ch chan struct{}) {},
			},
			want: true,
		},
		{
			name: "timeJSHelper.closeChannel() case 3",
			fields: fields{
				millis: 1,
				do: func(helper *timeJSHelper) {
					helper.stopExecute()
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
			helper := newTimeJSHelper()
			go tt.fields.secondDo(helper.stopCh)
			tt.fields.do(helper)
			if helper.isChannelClosed() != tt.want {
				t.Errorf("closeChannel() : isClosed = %v, want - %v", helper.isClosed, tt.want)
			}
		})
	}
}

func Test_timeJSHelper_stopExecute(t *testing.T) {
	type fields struct {
		millis   uint32
		do       func(*timeJSHelper)
		secondDo func(chan struct{})
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "timeJSHelper.stopExecute() case 1",
			fields: fields{
				millis:   1,
				do:       func(helper *timeJSHelper) {},
				secondDo: func(ch chan struct{}) {},
			},
			want: false,
		},
		{
			name: "timeJSHelper.stopExecute() case 2",
			fields: fields{
				millis: 1,
				do: func(helper *timeJSHelper) {
					helper.closeChannel()
				},
				secondDo: func(ch chan struct{}) {},
			},
			want: true,
		},
		{
			name: "timeJSHelper.stopExecute() case 3",
			fields: fields{
				millis: 1,
				do: func(helper *timeJSHelper) {
					helper.stopExecute()
				},
				secondDo: func(ch chan struct{}) {
					<-ch
				},
			},
			want: true,
		},
		{
			name: "timeJSHelper.stopExecute() case 4",
			fields: fields{
				millis: 1,
				do: func(helper *timeJSHelper) {
					helper.closeChannel()
					helper.stopExecute()
				},
				secondDo: func(ch chan struct{}) {},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			helper := newTimeJSHelper()
			go tt.fields.secondDo(helper.stopCh)
			tt.fields.do(helper)
			if helper.isChannelClosed() != tt.want {
				t.Errorf("stopExecute() : isClosed = %v, want - %v", helper.isClosed, tt.want)
			}
		})
	}
}
