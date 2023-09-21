package result

import "testing"

func Test_checkName(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "郭嘉guojia",
			want: true,
		},
		{
			name: "郭嘉guojia123?",
			want: false,
		},
		{
			name: "郭嘉guojia123_",
			want: true,
		},
		{
			name: "郭嘉guojia123123123123",
			want: false,
		},
		{
			name: "郭嘉guojia123。",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkName(tt.name); got != tt.want {
				t.Errorf("checkName() = %v, want %v", got, tt.want)
			}
		})
	}
}
