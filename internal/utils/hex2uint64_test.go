package utils

import "testing"

func TestHexToUint64(t *testing.T) {
	type args struct {
		hexString string
	}

	tests := []struct {
		name    string
		args    args
		want    uint64
		wantErr bool
	}{
		{
			name: "Valid Hex String (No Prefix)",
			args: args{
				hexString: "123456789abcdef0",
			},
			want:    0x123456789abcdef0,
			wantErr: false,
		},
		{
			name: "Valid Hex String (Short)",
			args: args{
				hexString: "0x10",
			},
			want:    0x10,
			wantErr: false,
		},
		{
			name: "Valid Hex String (Zero)",
			args: args{
				hexString: "0x0",
			},
			want:    0x0,
			wantErr: false,
		},
		{
			name: "Invalid Hex String (Non-Hex Characters)",
			args: args{
				hexString: "0x1234567890g",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Invalid Hex String (Empty)",
			args: args{
				hexString: "",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Invalid Hex String (Whitespace)",
			args: args{
				hexString: "   ",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Invalid Hex String (Only 0x)",
			args: args{
				hexString: "0x",
			},
			want:    0,
			wantErr: true,
		},
		{
			name: "Valid Hex String (Leading/Trailing Whitespace)",
			args: args{
				hexString: "  0x123  ",
			},
			want:    0x123,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := HexToUint64(tt.args.hexString)
			if (err != nil) != tt.wantErr {
				t.Errorf("HexToUint64() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("HexToUint64() = %v, want %v", got, tt.want)
			}
		})
	}
}
