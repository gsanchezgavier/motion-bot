package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Config is the unique source of configuration for the app
type Config struct {
	MetricsPort          string  `yaml:"metrics_port" env:"MBOT_METRICS_PORT" env-default:"9802" env-description:"Port where prometheus metrics are exposed with /metrics"`
	WatcherPort          string  `yaml:"watcher_port" env:"MBOT_WATCHER_PORT" env-default:"8000" env-description:"Port where motion-bot listens for motion events"`
	MotionControlURL     string  `yaml:"motion_control_port" env:"MBOT_MOTION_CONTROL_URL" env-default:"http://motion:8080" env-description:"Motion Control URL"`
	TelToken             string  `yaml:"tel_token" env:"MBOT_TEL_TOKEN" env-default:"" env-description:"Telegram Bot Token"`
	TelAdminID           []int   `yaml:"tel_admin_id" env:"MBOT_TEL_ADMIN_ID" env-default:"" env-description:"Telegram user IDs of motion-bot Admins"`
	TelSubscribedChatsID []int64 `yaml:"tel_subscribed_chats_id" env:"MBOT_TEL_SUBSCRIBED_CHATS_ID" env-default:"" env-description:"Default chats subscribed to this bot"`
	Verbose              bool    `yaml:"verbose" env:"MBOT_VERBOSE" env-default:"false"`
}

// New creates a Config object parsing a yml file. Configs could also be added as env variables
func New() Config {
	var cfg Config

	fset := flag.NewFlagSet("motion-bot", flag.ExitOnError)
	configPath := fset.String("config", "/etc/telebot/config.yml", "path to config file")
	fset.Usage = cleanenv.FUsage(fset.Output(), &cfg, nil, fset.Usage)

	if err := fset.Parse(os.Args[1:]); err != nil {
		log.Fatal().Err(err).Msg("error reading parsing args")
	}
	if c := os.Getenv("MBOT_CONFIG_PATH"); c != "" {
		*configPath = c
	}

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	err := cleanenv.ReadConfig(*configPath, &cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("error reading config file")
	}

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if cfg.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	return cfg
}
