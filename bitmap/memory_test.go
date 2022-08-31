package bitmap

import (
	"github.com/bits-and-blooms/bitset"
	"testing"
)

func TestLocal_CheckBits(t *testing.T) {
	type fields struct {
		bs *bitset.BitSet
		m  uint64
	}
	type args struct {
		locs []uint64
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		want      bool
		wantErr   bool
		doSetBits bool
	}{
		{
			name: "not exist",
			fields: fields{
				bs: bitset.New(500),
				m:  500,
			},
			args: args{
				locs: []uint64{
					12345,
					67890,
					13579,
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "exist",
			fields: fields{
				bs: bitset.New(500),
				m:  500,
			},
			args: args{
				locs: []uint64{
					12345,
					67890,
					13579,
				},
			},
			want:      true,
			wantErr:   false,
			doSetBits: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &InMemory{
				bs: tt.fields.bs,
				m:  tt.fields.m,
			}
			if tt.doSetBits {
				err := l.SetBits(tt.args.locs)
				if err != nil {
					t.Errorf("doSetBits failed: %v", err)
					return
				}
			}
			got, err := l.CheckBits(tt.args.locs)
			if (err != nil) != tt.wantErr {
				t.Errorf("CheckBits() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CheckBits() got = %v, want %v", got, tt.want)
			}
		})
	}
}
