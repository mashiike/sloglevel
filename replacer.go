package sloglevel

import (
	"log/slog"
	"strings"
)

type options struct {
	next             func([]string, slog.Attr) slog.Attr
	additionalLevels map[slog.Level]string
	toLower          bool
	newKey           string
}

// AttrReplacerOption is a function that configures the options.
type AttrReplacerOption func(*options)

// NewAttrReplacer returns a function that replaces the attributes.
func NewAttrReplacer(opts ...AttrReplacerOption) func([]string, slog.Attr) slog.Attr {
	o := options{
		additionalLevels: map[slog.Level]string{},
	}
	for _, opt := range opts {
		opt(&o)
	}
	if o.next == nil {
		o.next = func(groups []string, a slog.Attr) slog.Attr {
			return a
		}
	}
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key != slog.LevelKey {
			return o.next(groups, a)
		}
		level, ok := a.Value.Any().(slog.Level)
		if !ok {
			return o.next(groups, a)
		}
		key := a.Key
		levelStr := level.String()
		if v, ok := o.additionalLevels[level]; ok {
			levelStr = v
		}
		if o.toLower {
			levelStr = strings.ToLower(levelStr)
		}
		if o.newKey != "" {
			key = o.newKey
		}
		current := slog.Attr{
			Key:   key,
			Value: slog.StringValue(levelStr),
		}
		return o.next(groups, current)
	}
}

// NextAttrReplacer returns an option that sets the next attribute replacer. for other key
func NextAttrReplacer(next func([]string, slog.Attr) slog.Attr) AttrReplacerOption {
	return func(o *options) {
		o.next = next
	}
}

// AddLevel returns an option that adds a level to the replacer.
func AddLevel(level slog.Level, value string) AttrReplacerOption {
	return func(o *options) {
		o.additionalLevels[level] = value
	}
}

// ToLower() returns an option that converts the level to lowercase.
func ToLower() AttrReplacerOption {
	return func(o *options) {
		o.toLower = true
	}
}

// ChangeKey returns an option that changes the key of the level.
func ChangeKey(key string) AttrReplacerOption {
	return func(o *options) {
		o.newKey = key
	}
}
