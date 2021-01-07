package time2

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

type mock struct {
	D Duration `json:"d"`
}

func TestDuration_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		d       Duration
		want    []byte
		wantErr bool
	}{
		{
			name:    "1",
			d:       NewDuration(90 * time.Second),
			want:    []byte(`"1m30s"`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.d.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("Duration.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Duration.MarshalJSON() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	type args struct {
		b []byte
	}
	tests := []struct {
		name    string
		d       Duration
		args    args
		wantErr bool
		want    time.Duration
	}{
		{
			name: "1",
			d:    Duration{},
			args: args{
				b: []byte(`"1m30s"`),
			},
			wantErr: false,
			want:    90 * time.Second,
		},

		{
			name: "1",
			d:    Duration{},
			args: args{
				b: []byte(`""`),
			},
			wantErr: false,
			want:    0 * time.Second,
		},
		{
			name: "1",
			d:    Duration{},
			args: args{
				b: []byte(`"`),
			},
			wantErr: true,
			want:    0 * time.Second,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.d.UnmarshalJSON(tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("Duration.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.d.String() != tt.want.String() {
				t.Errorf("Duration.UnmarshalJSON() = %v, want %v", tt.d.Duration, tt.want)
			}
		})
	}
}

func TestDuration_Marshal(t *testing.T) {

	m := &mock{
		D: NewDuration(time.Second),
	}
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	want := `{"d":"1s"}`
	if got != want {
		t.Fatalf("got: %v want: %v", got, want)
	}
}

func TestDuration_Unmarshal(t *testing.T) {
	b := `{"d":"1s"}`
	m := mock{}
	err := json.Unmarshal([]byte(b), &m)
	if err != nil {
		t.Fatal(err)
	}
	want := NewDuration(time.Second)
	if m.D.String() != want.String() {
		t.Fatalf("got: %v want: %v", m.D, want)
	}
}
