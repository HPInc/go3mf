package common

import (
	"math"
	"reflect"
	"testing"

	"github.com/qmuntal/go3mf/pkg/semaphore"
	"github.com/qmuntal/go3mf/pkg/stack"
)

var callbackTrue = func(progress int, id ProgressIdentifier, data interface{}) bool {
	return true
}

var callbackFalse = func(progress int, id ProgressIdentifier, data interface{}) bool {
	return false
}

func TestNewProgressMonitor(t *testing.T) {
	tests := []struct {
		name string
		want *ProgressMonitor
	}{
		{
			name: "",
			want: &ProgressMonitor{
				lastCallbackResult: true,
				callbackMutex:      semaphore.NewSemaphore(),
				levels:             stack.NewItemStack(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewProgressMonitor(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewProgressMonitor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgressMonitor_QueryCancelled_True(t *testing.T) {
	p := NewProgressMonitor()
	p.progressCallback = callbackTrue
	ret := p.QueryCancelled()
	if ret == false {
		t.Error("QueryCancelled should return true if callback return true but returned false")
	}
}

func TestProgressMonitor_QueryCancelled_False(t *testing.T) {
	p := NewProgressMonitor()
	p.progressCallback = callbackFalse
	ret := p.QueryCancelled()
	if ret == true {
		t.Error("QueryCancelled should return false if callback return false but returned true")
	}
}

func TestProgressMonitor_Progress_True(t *testing.T) {
	p := NewProgressMonitor()
	p.lastCallbackResult = false
	p.progressCallback = callbackTrue
	ret := p.Progress(1.0, ProgressDone)
	if ret == false {
		t.Error("Progress should return true if callback return true but returned false")
	}
}

func TestProgressMonitor_Progress_False(t *testing.T) {
	p := NewProgressMonitor()
	p.lastCallbackResult = true
	p.progressCallback = callbackFalse
	ret := p.Progress(1.0, ProgressDone)
	if ret == true {
		t.Error("Progress should return false if callback return false but returned true")
	}
}

func TestProgressMonitor_Progress_Nil(t *testing.T) {
	p := NewProgressMonitor()
	p.lastCallbackResult = false
	p.progressCallback = nil
	ret := p.Progress(1.0, ProgressDone)
	if ret == false {
		t.Error("Progress should return true if callback is nil")
	}
}

func TestProgressMonitor_Progress_CantRun(t *testing.T) {
	p := NewProgressMonitor()
	p.lastCallbackResult = false
	p.progressCallback = callbackFalse
	p.callbackMutex.CanRun()
	ret := p.Progress(1.0, ProgressDone)
	if ret == false {
		t.Error("Progress should return true if can't run")
	}
}

func TestProgressMonitor_Progress_Done(t *testing.T) {
	p := NewProgressMonitor()
	p.progressCallback = callbackFalse
	p.callbackMutex.CanRun()
	p.callbackMutex.Done()
	ret := p.Progress(1.0, ProgressDone)
	if ret == true {
		t.Error("Progress should return false if callback return false but returned true")
	}
}

func TestProgressMonitor_PushLevel(t *testing.T) {
	type args struct {
		relativeStart float64
		relativeEnd   float64
	}
	type want struct {
		A float64
		B float64
	}
	p := NewProgressMonitor()
	tests := []struct {
		name string
		p    *ProgressMonitor
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

func TestProgressMonitor_PopLevel(t *testing.T) {
	p := NewProgressMonitor()
	p.PushLevel(0.0, 1.0)
	p.PushLevel(0.2, 1.0)
	tests := []struct {
		name string
		p    *ProgressMonitor
		want Float64Pair
	}{
		{"2", p, Float64Pair{0.2, 1.0}},
		{"1", p, Float64Pair{0.0, 1.0}},
		{"0", p, Float64Pair{0.0, 1.0}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.PopLevel(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProgressMonitor.PopLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProgressMonitor_ResetLevels_Empty(t *testing.T) {
	p := NewProgressMonitor()
	tests := []struct {
		name string
		p    *ProgressMonitor
	}{
		{"2", p},
		{"1", p},
		{"0", p},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.ResetLevels()
			if got := tt.p.level(); !reflect.DeepEqual(got, Float64Pair{0.0, 1.0}) {
				t.Errorf("expect initial values but got %f - %f", got, Float64Pair{0.0, 1.0})
			}
		})
	}
}

func TestProgressMonitor_ResetLevels(t *testing.T) {
	p := NewProgressMonitor()
	p.PushLevel(0.0, 1.0)
	p.PushLevel(0.2, 1.0)
	tests := []struct {
		name string
		p    *ProgressMonitor
	}{
		{"2", p},
		{"1", p},
		{"0", p},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.ResetLevels()
			if got := tt.p.level(); !reflect.DeepEqual(got, Float64Pair{0.0, 1.0}) {
				t.Errorf("expect initial values but got %f - %f", got, Float64Pair{0.0, 1.0})
			}
		})
	}
}

func TestProgressMonitor_SetProgressCallback(t *testing.T) {
	pr := Float64Pair{0.0, 1.0}
	type args struct {
		callback ProgressCallback
		userData interface{}
	}
	tests := []struct {
		name string
		p    *ProgressMonitor
		args args
		want *ProgressMonitor
	}{
		{"true", NewProgressMonitor(), args{callbackTrue, 2}, &ProgressMonitor{callbackTrue, 2, true, stack.NewItemStack(), semaphore.NewSemaphore()}},
		{"false", NewProgressMonitor(), args{callbackFalse, "aaa"}, &ProgressMonitor{callbackFalse, "aaa", true, stack.NewItemStack(), semaphore.NewSemaphore()}},
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

func TestProgressMonitor_ClearProgressCallback(t *testing.T) {
	p := NewProgressMonitor()
	p.progressCallback = callbackTrue
	p.userData = 2
	tests := []struct {
		name string
		p    *ProgressMonitor
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

func TestProgressMonitor_WasAborted_True(t *testing.T) {
	p := NewProgressMonitor()
	p.lastCallbackResult = true
	ret := p.WasAborted()
	if ret == true {
		t.Error("WasAborted should return true if callback return true but returned false")
	}
}

func TestProgressMonitor_WasAborted_False(t *testing.T) {
	p := NewProgressMonitor()
	p.lastCallbackResult = false
	ret := p.WasAborted()
	if ret == false {
		t.Error("WasAborted should return false if callback return false but returned true")
	}
}

func TestProgressMonitor_GetProgressMessage(t *testing.T) {
	p := NewProgressMonitor()
	type args struct {
		progressIdentifier ProgressIdentifier
	}
	tests := []struct {
		name string
		p    *ProgressMonitor
		args args
		want string
	}{
		{"ProgressQueryCanceled", p, args{ProgressQueryCanceled}, ""},
		{"ProgressDone", p, args{ProgressDone}, "Done"},
		{"ProgressCleanup", p, args{ProgressCleanup}, "Cleaning up"},
		{"ProgressReadStream", p, args{ProgressReadStream}, "Reading stream"},
		{"ProgressExtractOPCPackage", p, args{ProgressExtractOPCPackage}, "Extracting OPC package"},
		{"ProgressReadNonRootModels", p, args{ProgressReadNonRootModels}, "Reading non-root models"},
		{"ProgressReadRootModel", p, args{ProgressReadRootModel}, "Reading root model"},
		{"ProgressReadResources", p, args{ProgressReadResources}, "Reading resources"},
		{"ProgressReadMesh", p, args{ProgressReadMesh}, "Reading mesh data"},
		{"ProgressReadSlices", p, args{ProgressReadSlices}, "Reading slice data"},
		{"ProgressReadBuild", p, args{ProgressReadBuild}, "Reading build definition"},
		{"ProgressCreateOPCPackage", p, args{ProgressCreateOPCPackage}, "Creating OPC package"},
		{"ProgressWriteModelsToStream", p, args{ProgressWriteModelsToStream}, "Writing models to stream"},
		{"ProgressWriteRootModel", p, args{ProgressWriteRootModel}, "Writing root model"},
		{"ProgressWriteNonRootModels", p, args{ProgressWriteNonRootModels}, "Writing non-root models"},
		{"ProgressWriteAttachements", p, args{ProgressWriteAttachements}, "Writing attachments"},
		{"ProgressWriteContentTypes", p, args{ProgressWriteContentTypes}, "Writing content types"},
		{"ProgressWriteObjects", p, args{ProgressWriteObjects}, "Writing objects"},
		{"ProgressWriteNodes", p, args{ProgressWriteNodes}, "Writing Nodes"},
		{"ProgressWriteTriangles", p, args{ProgressWriteTriangles}, "Writing triangles"},
		{"ProgressWriteSlices", p, args{ProgressWriteSlices}, "Writing slices"},
		{"Unknown", p, args{-1}, "Unknown Progress Identifier"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.p.GetProgressMessage(tt.args.progressIdentifier); got != tt.want {
				t.Errorf("ProgressMonitor.GetProgressMessage() = %v, want %v", got, tt.want)
			}
		})
	}
}
