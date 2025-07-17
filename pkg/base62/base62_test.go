package base62

import (
	"testing"
)

func TestMain(m *testing.M) {
	MustInit("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	m.Run()
}

func TestInt2String(t *testing.T) {
	type args struct {
		seq uint64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "case:0", args: args{seq: 0}, want: "0"},
		{name: "case:1", args: args{seq: 1}, want: "1"},
		{name: "case:62", args: args{seq: 62}, want: "10"},
		{name: "case:6347", args: args{seq: 6347}, want: "1En"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Int2String(tt.args.seq); got != tt.want {
				t.Errorf("Int2String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestString2Int(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		wantSeq uint64
	}{
		{name: "case 0:", args: args{s: "0"}, wantSeq: 0},
		{name: "case 10:", args: args{s: "10"}, wantSeq: 62},
		{name: "case 1En:", args: args{s: "1En"}, wantSeq: 6347},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotSeq := String2Int(tt.args.s); gotSeq != tt.wantSeq {
				t.Errorf("String2Int() = %v, want %v", gotSeq, tt.wantSeq)
			}
		})
	}
}
