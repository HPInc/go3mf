package common

import (
	"math"
	"reflect"
	"testing"

	"github.com/qmuntal/go3mf/pkg/semaphore"
	"github.com/qmuntal/go3mf/pkg/stack"
)

var callbackTrue = func(progress int, id Stage, data interface{}) bool {
	return true
}

var callbackFalse = func(progress int, id Stage, data interface{}) bool {
	return false
}

func TestNewMonitor(t *testing.T) {
	tests := []struct {
		name string
		want *Monitor
	}{
		{
			name: "",
			want: &Monitor{
				lastCallbackResult: true,
				callbackMutex:      semaphore.NewSemaphore(),
				levels:             stack.NewItemStack(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMonitor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMonitor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMonitor_QueryCancelled_True(t *testing.T) {
	p := NewMonitor()
	p.progressCallback = callbackTrue
	ret := p.QueryCancelled()
	if ret == false {
		t.Error("QueryCancelled should return true if callback return true but returned false")
	}
}

func TestMonitor_QueryCancelled_False(t *testing.T) {
	p := NewMonitor()
	p.progressCallback = callbackFalse
	ret := p.QueryCancelled()
	if ret == true {
		t.Error("QueryCancelled should return false if callback return false but returned true")
	}
}

func TestMonitor_Progress_True(t *testing.T) {
	p := NewMonitor()
	p.lastCallbackResult = false
	p.progressCallback = callbackTrue
	ret := p.Progress(1.0, StageDone)
	if ret == false {
		t.Error("Progress should return true if callback return true but returned false")
	}
}

func TestMonitor_Progress_False(t *testing.T) {
	p := NewMonitor()
	p.lastCallbackResult = true
	p.progressCallback = callbackFalse
	ret := p.Progress(1.0, StageDone)
	if ret == true {
		t.Error("Progress should return false if callback return false but returned true")
	}
}

func TestMonitor_Progress_Nil(t *testing.T) {
	p := NewMonitor()
	p.lastCallbackResult = false
	p.progressCallback = nil
	ret := p.Progress(1.0, StageDone)
	if ret == false {
		t.Error("Progress should return true if callback is nil")
	}
}

func TestMonitor_Progress_CantRun(t *testing.T) {
	p := NewMonitor()
	p.lastCallbackResult = false
	p.progressCallback = callbackFalse
	p.callbackMutex.CanRun()
	ret := p.Progress(1.0, StageDone)
	if ret == false {
		t.Error("Progress should return true if can't run")
	}
}

func TestMonitor_Progress_Done(t *testing.T) {
	p := NewMonitor()
	p.progressCallback = callbackFalse
	p.callbackMutex.CanRun()
	p.callbackMutex.Done()
	ret := p.Progress(1.0, StageDone)
	if ret == true {
		t.Error("Progress should return false if callback return false but returned true")
	}
}

func TestMonitor_PushLevel(t *testing.T) {
	type args struct {
		relativeStart float64
		relativeEnd   float64
	}
	type want struct {
		A float64
		B float64
	}
	p := NewMonitor()
	tests := []struct {
		name string
		p    *Monitor
		args args
		want want
	}{
		{"0-1", p, args{0.0, 1.0}, want{0.0, 1.0}},
		{"0.2-1.0", p, args{0.2, 1.0}, want{0.2, 1.0}},
		{"0.2-1.0", p, args{0.2, 1.0}, want{0.36, 1.0}},
		{"0.4-0.8", p, args{0.4, 0.8}, want{0.616, 0.872}},
		{"1.0-1.0", p, args{1.0, 1.0}, want{0.872, 0.872}},
		{"0.0-0.0", p, args{0.0, 0.0}, want{0.872, 0.872}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.PushLevel(tt.args.relativeStart, tt.args.relativeEnd)
			if top := tt.p.level(); math.Abs(top.A-tt.want.A) > 0.001 || math.Abs(top.B-tt.want.B) > 0.001 {
				t.Errorf("wrong level values, expected %f - %f but got %f - %f", tt.want.A, tt.want.B, top.A, top.B)
			}
		})
	}
}

func TestMonitor_PopLevel(t *testing.T) {
	p := NewMonitor()
	p.PushLevel(0.0, 1.0)
	p.PushLevel(0.2, 1.0)
	tests := []struct {
		name string
		p    *Monitor
		want float64Pair
	}{
		{"2", p, float64Pair{0.2, 1.0}},
		{"1", p, float64Pair{0.0, 1.0}},
		{"0", p, float64Pair{0.0, 1.0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if a, b := tt.p.PopLevel(); a != tt.want.A || b != tt.want.B {
				t.Errorf("Monitor.PopLevel() = %v, %v, want %v", a, b, tt.want)
			}
		})
	}
}

func TestMonitor_ResetLevels_Empty(t *testing.T) {
	p := NewMonitor()
	tests := []struct {
		name string
		p    *Monitor
	}{
		{"2", p},
		{"1", p},
		{"0", p},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.ResetLevels()
			if got := tt.p.level(); !reflect.DeepEqual(got, float64Pair{0.0, 1.0}) {
				t.Errorf("expect initial values but got %f - %f", got, float64Pair{0.0, 1.0})
			}
		})
	}
}

func TestMonitor_ResetLevels(t *testing.T) {
	p := NewMonitor()
	p.PushLevel(0.0, 1.0)
	p.PushLevel(0.2, 1.0)
	tests := []struct {
		name string
		p    *Monitor
	}{
		{"2", p},
		{"1", p},
		{"0", p},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.ResetLevels()
			if got := tt.p.level(); !reflect.DeepEqual(got, float64Pair{0.0, 1.0}) {
				t.Errorf("expect initial values but got %f - %f", got, float64Pair{0.0, 1.0})
			}
		})
	}
}

func TestMonitor_SetProgressCallback(t *testing.T) {
	pr := float64Pair{0.0, 1.0}
	type args struct {
		callback progressCallback
		userData interface{}
	}
	tests := []struct {
		name string
		p    *Monitor
		args args
		want *Monitor
	}{
		{"true", NewMonitor(), args{callbackTrue, 2}, &Monitor{callbackTrue, 2, true, stack.NewItemStack(), semaphore.NewSemaphore()}},
		{"false", NewMonitor(), args{callbackFalse, "aaa"}, &Monitor{callbackFalse, "aaa", true, stack.NewItemStack(), semaphore.NewSemaphore()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.SetProgressCallback(tt.args.callback, tt.args.userData)
			if tt.p.userData != tt.want.userData || tt.p.level() != pr || tt.p.lastCallbackResult != true {
				t.Error("expected restarted monitor")
			}
		})
	}
}

func TestMonitor_ClearProgressCallback(t *testing.T) {
	p := NewMonitor()
	p.progressCallback = callbackTrue
	p.userData = 2
	tests := []struct {
		name string
		p    *Monitor
	}{
		{"base", p},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.ClearProgressCallback()
			if p.progressCallback != nil {
				t.Error("callback expected to be nil")
			}

			if p.userData != nil {
				t.Error("userData expected to be nil")
			}
		})
	}
}

func TestMonitor_WasAborted_True(t *testing.T) {
	p := NewMonitor()
	p.lastCallbackResult = true
	ret := p.WasAborted()
	if ret == true {
		t.Error("WasAborted should return true if callback return true but returned false")
	}
}

func TestMonitor_WasAborted_False(t *testing.T) {
	p := NewMonitor()
	p.lastCallbackResult = false
	ret := p.WasAborted()
	if ret == false {
		t.Error("WasAborted should return false if callback return false but returned true")
	}
}

func TestMonitor_ProgressMessage(t *testing.T) {
	p := NewMonitor()
	type args struct {
		progressIdentifier Stage
	}
	tests := []struct {
		name string
		p    *Monitor
		args args
		want string
	}{
		{"StageQueryCanceled", p, args{StageQueryCanceled}, ""},
		{"StageDone", p, args{StageDone}, "Done"},
		{"StageCleanup", p, args{StageCleanup}, "Cleaning up"},
		{"StageReadStream", p, args{StageReadStream}, "Reading stream"},
		{"StageExtractOPCPackage", p, args{StageExtractOPCPackage}, "Extracting OPC package"},
		{"StageReadNonRootModels", p, args{StageReadNonRootModels}, "Reading non-root models"},
		{"StageReadRootModel", p, args{StageReadRootModel}, "Reading root model"},
		{"StageReadResources", p, args{StageReadResources}, "Reading resources"},
		{"StageReadMesh", p, args{StageReadMesh}, "Reading mesh data"},
		{"StageReadSlices", p, args{StageReadSlices}, "Reading slice data"},
		{"StageReadBuild", p, args{StageReadBuild}, "Reading build definition"},
		{"StageCreateOPCPackage", p, args{StageCreateOPCPackage}, "Creating OPC package"},
		{"StageWriteModelsToStream", p, args{StageWriteModelsToStream}, "Writing models to stream"},
		{"StageWriteRootModel", p, args{StageWriteRootModel}, "Writing root model"},
		{"StageWriteNonRootModels", p, args{StageWriteNonRootModels}, "Writing non-root models"},
		{"StageWriteAttachements", p, args{StageWriteAttachements}, "Writing attachments"},
		{"StageWriteContentTypes", p, args{StageWriteContentTypes}, "Writing content types"},
		{"StageWriteObjects", p, args{StageWriteObjects}, "Writing objects"},
		{"StageWriteNodes", p, args{StageWriteNodes}, "Writing Nodes"},
		{"StageWriteTriangles", p, args{StageWriteTriangles}, "Writing triangles"},
		{"StageWriteSlices", p, args{StageWriteSlices}, "Writing slices"},
		{"Unknown", p, args{-1}, "Unknown Progress Identifier"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.ProgressMessage(tt.args.progressIdentifier); got != tt.want {
				t.Errorf("Monitor.ProgressMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
