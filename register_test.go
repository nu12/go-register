package register

import (
	"errors"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string
		want *Register
	}{
		{name: "Get new Register", want: &Register{Step: 0, Err: nil}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegister_Error(t *testing.T) {
	type fields struct {
		Err error
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{name: "Register without error", fields: fields{Err: nil}, wantErr: false},
		{name: "Register without error", fields: fields{Err: errors.New("undefined error")}, wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Register{
				Err: tt.fields.Err,
			}
			if err := r.Error(); (err != nil) != tt.wantErr {
				t.Errorf("Error() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestRegister_Run(t *testing.T) {
	type fields struct {
		Step int
		Err  error
	}
	type args struct {
		fo []func()
	}
	tests := []struct {
		name     string
		register *Register
		args     args
		want     int
	}{
		{name: "Run 1 function", register: &Register{Step: 0, Err: nil}, args: args{fo: []func(){func() {}}}, want: 1},
		{name: "Run 2 functions", register: &Register{Step: 0, Err: nil}, args: args{fo: []func(){func() {}, func() {}}}, want: 2},
		{name: "Skip function if error exists", register: &Register{Step: 3, Err: errors.New("undefined error")}, args: args{fo: []func(){func() {}}}, want: 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := tt.register
			if got := r.Run(tt.args.fo...).Step; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegister_IfError(t *testing.T) {
	target := 0

	var changeTarget = func() {
		target += 1
	}

	type fields struct {
		Step int
		Err  error
	}
	type args struct {
		f func()
	}
	tests := []struct {
		name     string
		register *Register
		args     args
		want     int
	}{
		{name: "Run if error exists", register: &Register{Step: 0, Err: errors.New("undefined error")}, args: args{f: changeTarget}, want: 1},
		{name: "Skip if no error", register: &Register{Step: 0, Err: nil}, args: args{f: changeTarget}, want: 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.register.IfError(tt.args.f)
			if got := target; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRegister_If(t *testing.T) {
	type fields struct {
		Step int
		Err  error
	}
	type args struct {
		b  bool
		fo []func()
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   int
	}{
		{name: "Run functions if b is true", fields: fields{Step: 0, Err: nil}, args: args{b: true, fo: []func(){func() {}, func() {}}}, want: 2},
		{name: "Skip functions if b is false", fields: fields{Step: 0, Err: nil}, args: args{b: false, fo: []func(){func() {}, func() {}}}, want: 0},
		{name: "Skip functions if error exist", fields: fields{Step: 0, Err: errors.New("undefined error")}, args: args{b: true, fo: []func(){func() {}, func() {}}}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Register{
				Step: tt.fields.Step,
				Err:  tt.fields.Err,
			}
			if got := r.If(tt.args.b, tt.args.fo...).Step; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("If() = %v, want %v", got, tt.want)
			}
		})
	}
}
