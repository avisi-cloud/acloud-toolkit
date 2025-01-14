package helpers

import (
	"context"
	"errors"
	"log"
	"sync"
	"testing"
	"time"
)

func TestRetry(t *testing.T) {
	globalAttempt := 0

	type args struct {
		attempts int
		sleep    time.Duration
		f        func() error
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		wantedAttempts int
	}{
		{
			name: "1 successful attempt",
			args: args{
				attempts: 1,
				sleep:    50 * time.Millisecond,
				f: func() error {
					globalAttempt = globalAttempt + 1
					return nil
				},
			},
			wantErr:        false,
			wantedAttempts: 1,
		},
		{
			name: "1 error attempt",
			args: args{
				attempts: 1,
				sleep:    50 * time.Millisecond,
				f: func() error {
					globalAttempt = globalAttempt + 1
					return errors.New("failed with an error")
				},
			},
			wantErr:        true,
			wantedAttempts: 1,
		},
		{
			name: "1 error attempt then 1 successful",
			args: args{
				attempts: 2,
				sleep:    50 * time.Millisecond,
				f: func() error {
					globalAttempt = globalAttempt + 1
					if globalAttempt == 1 {
						return errors.New("first attempt fails")
					}
					return nil
				},
			},
			wantErr:        false,
			wantedAttempts: 2,
		},
		{
			name: "3 error attempt",
			args: args{
				attempts: 3,
				sleep:    50 * time.Millisecond,
				f: func() error {
					globalAttempt = globalAttempt + 1
					return errors.New("failed with an error")
				},
			},
			wantErr:        true,
			wantedAttempts: 3,
		},
	}
	for _, tt := range tests {
		globalAttempt = 0
		t.Run(tt.name, func(t *testing.T) {
			if err := Retry(tt.args.attempts, tt.args.sleep, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("Retry() error = %v, wantErr %v", err, tt.wantErr)
			}
			if globalAttempt != tt.wantedAttempts {
				t.Errorf("Retry() globalAttampt = %d, wantedAttepts %d", globalAttempt, tt.wantedAttempts)
			}
		})
	}
}

func TestWithDefault(t *testing.T) {
	returnCh := make(chan error)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
	Loop:
		for {
			select {
			case err := <-returnCh:
				log.Println("returnCh: " + err.Error())
				break Loop
				// default:
				//	log.Println("default:")
			}
		}
		wg.Done()
	}()

	time.Sleep(5 * time.Second)
	returnCh <- errors.New("error")
	wg.Wait()
}

func TestRetryWithCancel(t *testing.T) {
	globalAttampt := 0
	globalContext, globalCancel := context.WithCancel(context.Background())

	type args struct {
		ctx      context.Context
		attempts int
		sleep    time.Duration
		f        func() error
	}
	tests := []struct {
		name           string
		args           args
		wantErr        bool
		wantedAttempts int
	}{
		{
			name: "1 successful attempt",
			args: args{
				ctx:      context.Background(),
				attempts: 1,
				sleep:    50 * time.Millisecond,
				f: func() error {
					globalAttampt = globalAttampt + 1
					return nil
				},
			},
			wantErr:        false,
			wantedAttempts: 1,
		},
		{
			name: "1 error attempt",
			args: args{
				ctx:      context.Background(),
				attempts: 1,
				sleep:    50 * time.Millisecond,
				f: func() error {
					globalAttampt = globalAttampt + 1
					return errors.New("failed with an error")
				},
			},
			wantErr:        true,
			wantedAttempts: 1,
		},
		{
			name: "1 error attempt then 1 successful",
			args: args{
				ctx:      context.Background(),
				attempts: 2,
				sleep:    50 * time.Millisecond,
				f: func() error {
					globalAttampt = globalAttampt + 1
					if globalAttampt == 1 {
						return errors.New("first attempt fails")
					}
					return nil
				},
			},
			wantErr:        false,
			wantedAttempts: 2,
		},
		{
			name: "3 error attempt",
			args: args{
				ctx:      context.Background(),
				attempts: 3,
				sleep:    50 * time.Millisecond,
				f: func() error {
					globalAttampt = globalAttampt + 1
					return errors.New("failed with an error")
				},
			},
			wantErr:        true,
			wantedAttempts: 3,
		},
		{
			name: "cancel after first try",
			args: args{
				ctx:      globalContext,
				attempts: 3,
				sleep:    50 * time.Millisecond,
				f: func() error {
					globalAttampt = globalAttampt + 1
					globalCancel()
					return errors.New("failed with an error")
				},
			},
			wantErr:        true,
			wantedAttempts: 1,
		},
	}
	for _, tt := range tests {
		globalAttampt = 0
		t.Run(tt.name, func(t *testing.T) {
			if err := RetryWithCancel(tt.args.ctx, tt.args.attempts, tt.args.sleep, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("RetryWithCancel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if globalAttampt != tt.wantedAttempts {
				t.Errorf("RetryWithCancel() globalAttampt = %d, wantedAttepts %d", globalAttampt, tt.wantedAttempts)
			}
		})
		globalCancel()
	}
}
