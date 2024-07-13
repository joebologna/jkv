package jkv

import "context"

type baseCmd struct {
	err error
}

type StatusCmd struct {
	baseCmd
	val string
}

func NewStatusCmd(val string, err error) *StatusCmd {
	return &StatusCmd{baseCmd: baseCmd{err: err}, val: val}
}

type StringCmd struct {
	baseCmd
	val string
}

func NewStringCmd(val string, err error) *StringCmd {
	return &StringCmd{baseCmd: baseCmd{err: err}, val: val}
}

type IntCmd struct {
	baseCmd
	val int64
}

func NewIntCmd(val int64, err error) *IntCmd {
	return &IntCmd{baseCmd: baseCmd{err: err}, val: val}
}

type BoolCmd struct {
	baseCmd
	val bool
}

func NewBoolCmd(val bool, err error) *BoolCmd {
	return &BoolCmd{baseCmd: baseCmd{err: err}, val: val}
}

type StringSliceCmd struct {
	baseCmd

	val []string
}

func NewStringSliceCmd(val []string, err error) *StringSliceCmd {
	return &StringSliceCmd{baseCmd: baseCmd{err: err}, val: val}
}

func (s *StringCmd) Val() string        { return s.val }
func (s *StringCmd) Err() error         { return s.err }
func (s *IntCmd) Val() int64            { return s.val }
func (s *IntCmd) Err() error            { return s.err }
func (s *BoolCmd) Val() bool            { return s.val }
func (s *BoolCmd) Err() error           { return s.err }
func (s *StringSliceCmd) Val() []string { return s.val }
func (s *StringSliceCmd) Err() error    { return s.err }
func (s *StatusCmd) Val() string        { return s.val }
func (s *StatusCmd) Err() error         { return s.err }

type Client interface {
	Open() error
	Close()
	FlushDB(ctx context.Context) *StatusCmd
	Get(ctx context.Context, key string) *StringCmd
	Set(ctx context.Context, key, value string) *StatusCmd
	Del(ctx context.Context, keys ...string) *IntCmd
	Keys(ctx context.Context, pattern string) *StringSliceCmd
	Exists(ctx context.Context, keys ...string) *IntCmd
	HGet(ctx context.Context, hash, key string) *StringCmd
	HSet(ctx context.Context, hash, key string, values ...string) *IntCmd
	HDel(ctx context.Context, hash, key string) *IntCmd
	HKeys(ctx context.Context, hash string) *StringSliceCmd
	HExists(ctx context.Context, hash, key string) *BoolCmd
	Ping(ctx context.Context) *StatusCmd
}
