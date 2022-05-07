package pairs

import "github.com/bluecolor/tractor/pkg/lib/types"

const (
	SendProgress = "send_progress"
)

func Get(key string, pairs ...types.Pair) (p *types.Pair) {
	for _, pair := range pairs {
		if pair.Key == key {
			p = &pair
			break
		}
	}
	return
}
func GetOr[T any](key string, defaultValue T, pairs ...types.Pair) T {
	var p = Get(key, pairs...)
	if p == nil {
		p = &types.Pair{Key: key, Value: defaultValue}
	}
	return p.Value.(T)
}

func WithSendProgress(value ...bool) (p types.Pair) {
	if len(value) > 0 {
		p = types.Pair{Key: SendProgress, Value: value[0]}
	} else {
		p = types.Pair{Key: SendProgress, Value: true}
	}
	return
}
