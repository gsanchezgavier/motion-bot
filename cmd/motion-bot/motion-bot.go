package main

import (
	"gsanchezgavier/telegram-bot/internal/config"
	"gsanchezgavier/telegram-bot/internal/metrics"
	"gsanchezgavier/telegram-bot/internal/motion"
	"gsanchezgavier/telegram-bot/internal/telebot"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg := config.New()

	metrics.Start(cfg.MetricsPort)

	motion := motion.New(cfg.MotionControlURL, cfg.WatcherPort)

	tb, err := telebot.New(
		cfg.TelToken,
		cfg.TelAdminID,
		cfg.TelSubscribedChatsID,
		motion,
	)
	if err != nil {
		log.Fatal().Err(err).Msg("fail to create telebot")
	}

	tb.SendToAuthChats("telebot started")

	log.Info().Str("metrics_port", cfg.MetricsPort).Str("watcher_port", cfg.WatcherPort).
		Str("motion_control_url", cfg.MotionControlURL).Msg("telebot started")
	//TODO Add gracefully shutdown

	for {
		select {
		case v := <-motion.Video:
			log.Debug().Msgf("video received:%v", v)
			tb.SendVideo(v)
			metrics.MovementDetectionEventsTotal.WithLabelValues("video").Inc()
		case p := <-motion.Picture:
			log.Debug().Msgf("picture received:%v", p)
			tb.SendPicture(p)
			metrics.MovementDetectionEventsTotal.WithLabelValues("picture").Inc()
		}
	}

}
