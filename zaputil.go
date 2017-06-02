package zaputil

import (
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ConfigOption func(*zap.Config)

func Encoding(enc string) ConfigOption {
	return func(c *zap.Config) {
		c.Encoding = enc
	}
}

func DisableStackTrace() ConfigOption {
	return func(c *zap.Config) {
		c.DisableStacktrace = true
	}
}

func OutputPaths(p ...string) ConfigOption {
	return func(c *zap.Config) {
		c.OutputPaths = p
	}
}

func Level(l zapcore.Level) ConfigOption {
	return func(c *zap.Config) {
		c.Level.SetLevel(l)
	}
}

func DisableCaller() ConfigOption {
	return func(c *zap.Config) {
		c.DisableCaller = true
	}
}

func DisableTimestamp() ConfigOption {
	return func(c *zap.Config) {
		c.EncoderConfig.TimeKey = ""
	}
}

func LevelEncoder(e zapcore.LevelEncoder) ConfigOption {
	return func(c *zap.Config) {
		c.EncoderConfig.EncodeLevel = e
	}
}

func AttachLevelHandler(m *http.ServeMux, p string) ConfigOption {
	return func(c *zap.Config) {
		m.Handle(p, c.Level)
	}
}

func Template(env string) zap.Config {
	var c zap.Config
	switch env {
	case "production", "prod":
		c = zap.NewProductionConfig()
		DisableTimestamp()(&c)
	case "staging", "stg":
		c = zap.NewDevelopmentConfig()
		DisableTimestamp()(&c)
	default:
		c = zap.NewDevelopmentConfig()
	}
	Encoding("console")(&c)
	LevelEncoder(zapcore.CapitalColorLevelEncoder)(&c)
	OutputPaths("stdout")(&c)
	return c
}

func Config(t zap.Config, opts ...ConfigOption) zap.Config {
	for _, o := range opts {
		o(&t)
	}
	return t
}

func Must(l *zap.Logger, e error) *zap.Logger {
	if e != nil {
		panic(e)
	}
	return l
}

var newFancyID = func() func() string {
	adjectives := [...]string{"autumn", "hidden", "bitter", "misty", "silent", "empty",
		"dry", "dark", "summer", "icy", "delicate", "quiet", "white", "cool", "spring",
		"winter", "patient", "twilight", "dawn", "crimson", "wispy", "weathered", "blue",
		"billowing", "broken", "cold", "damp", "falling", "frosty", "green", "long", "late",
		"lingering", "bold", "little", "morning", "muddy", "old", "red", "rough", "still",
		"small", "sparkling", "throbbing", "shy", "wandering", "withered", "wild", "black",
		"young", "holy", "solitary", "fragrant", "aged", "snowy", "proud", "floral",
		"restless", "divine", "polished", "ancient", "purple", "lively", "nameless"}
	nouns := [...]string{"waterfall", "river", "breeze", "moon", "rain", "wind", "sea",
		"morning", "snow", "lake", "sunset", "pine", "shadow", "leaf", "dawn", "glitter",
		"forest", "hill", "cloud", "meadow", "sun", "glade", "bird", "brook", "butterfly",
		"bush", "dew", "dust", "field", "fire", "flower", "firefly", "feather", "grass",
		"haze", "mountain", "night", "pond", "darkness", "snowflake", "silence", "sound",
		"sky", "shape", "surf", "thunder", "violet", "water", "wildflower", "wave", "water",
		"resonance", "sun", "wood", "dream", "cherry", "tree", "fog", "frost", "voice",
		"paper", "frog", "smoke", "star"}
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	return func() string {
		return fmt.Sprintf("%s-%s-%d", adjectives[rnd.Intn(len(adjectives))], nouns[rnd.Intn(len(nouns))], rnd.Intn(9999))
	}
}()

func InstanceID() zapcore.Field {
	id, err := os.Hostname()
	if err != nil {
		id = newFancyID()
	}
	return zap.String("instance", id)
}
