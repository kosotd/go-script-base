package env

import (
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type Env interface {
	Bool(name string, value bool, usage string) *bool
	String(name string, value string, usage string) *string
	Duration(name string, value time.Duration, usage string) *time.Duration
	Float64(name string, value float64, usage string) *float64
	Int(name string, value int, usage string) *int
	Int64(name string, value int64, usage string) *int64
}

type env struct{}

func (e *env) Bool(name string, value bool, usage string) *bool {
	return pflag.Bool(name, value, usage)
}

func (e *env) String(name string, value string, usage string) *string {
	return pflag.String(name, value, usage)
}

func (e *env) Duration(name string, value time.Duration, usage string) *time.Duration {
	return pflag.Duration(name, value, usage)
}

func (e *env) Float64(name string, value float64, usage string) *float64 {
	return pflag.Float64(name, value, usage)
}

func (e *env) Int(name string, value int, usage string) *int {
	return pflag.Int(name, value, usage)
}

func (e *env) Int64(name string, value int64, usage string) *int64 {
	return pflag.Int64(name, value, usage)
}

func Init(callback func(e Env)) {
	callback(&env{})

	pflag.Parse()
	_ = viper.BindPFlags(pflag.CommandLine)
	viper.AutomaticEnv()
}
