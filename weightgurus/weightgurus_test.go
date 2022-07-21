package WeightGurus

import (
	"reflect"
	"testing"
)

func Test_convertWeightGuruNumToFloat(t *testing.T) {
	type args struct {
		weightGurusNum float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"convertWeightGuruNumToFloat", args{weightGurusNum: 2164}, 216.4},
		{"convertWeightGuruNumToFloat", args{weightGurusNum: 300}, 30.0},
		{"convertWeightGuruNumToFloat", args{weightGurusNum: 5500}, 550.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := convertWeightGuruNumToFloat(tt.args.weightGurusNum); got != tt.want {
				t.Errorf("convertWeightGuruNumToFloat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_removeDeletedOperation(t *testing.T) {
	type args struct {
		deletedOperation WeightGuruOperation
		weightHistory    []WeightGuruOperation
	}
	tests := []struct {
		name string
		args args
		want []WeightGuruOperation
	}{
		{"removeDeletedOperation", args{deletedOperation: WeightGuruOperation{entryTimestamp: "2019-01-01T00:00:00.000Z"},
			weightHistory: []WeightGuruOperation{{entryTimestamp: "2019-01-01T00:00:00.000Z"}, {entryTimestamp: "2020-01-01T00:00:00.000Z"}}}, []WeightGuruOperation{{entryTimestamp: "2020-01-01T00:00:00.000Z"}},
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeDeletedOperation(tt.args.deletedOperation, &tt.args.weightHistory); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("removeDeletedOperation() = %v, want %v", got, tt.want)
			}
		})
	}
}
