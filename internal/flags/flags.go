package flags

import (
	"github.com/alecthomas/kong"
)

// API the flags passed to the API service
type API struct {
	Version kong.VersionFlag `mapstructure:"version,omitempty"`

	Debug           bool   `help:"Enable debug logging." env:"DEBUG" mapstructure:"debug,omitempty"`
	RawEventLogging bool   `help:"Enable raw event logging." env:"RAW_EVENT_LOGGING" mapstructure:"raw_event_logging,omitempty"`
	Stage           string `help:"The development stage." env:"STAGE" mapstructure:"stage,omitempty"`
	Branch          string `help:"The git branch this code originated." env:"BRANCH" mapstructure:"branch,omitempty"`
	ReleaseTable    string `help:"The DynamoDB table used to store release data." env:"RELEASE_TABLE" mapstructure:"-"`
}
