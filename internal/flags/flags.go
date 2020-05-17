package flags

import (
	"log"

	"github.com/alecthomas/kong"
	"github.com/mitchellh/mapstructure"
)

// Events the flags passed to the events service
type Events struct {
	Version kong.VersionFlag `mapstructure:"version,omitempty"`

	Debug           bool   `help:"Enable debug logging." env:"DEBUG" mapstructure:"debug,omitempty"`
	RawEventLogging bool   `help:"Enable raw event logging." env:"RAW_EVENT_LOGGING" mapstructure:"raw_event_logging,omitempty"`
	Stage           string `help:"The development stage." env:"STAGE" mapstructure:"stage,omitempty"`
	Branch          string `help:"The git branch this code originated." env:"BRANCH" mapstructure:"branch,omitempty"`
}

// MustFields generate a fields map or panic
func (ev Events) MustFields() map[string]interface{} {

	fields := map[string]interface{}{}

	err := mapstructure.Decode(&ev, &fields)
	if err != nil {
		log.Fatalf("failed to build fields map: %+v", err) // game over 💥
	}

	return fields
}