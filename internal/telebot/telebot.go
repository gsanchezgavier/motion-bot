package telebot

import (
	"fmt"
	"gsanchezgavier/telegram-bot/internal/metrics"
	"gsanchezgavier/telegram-bot/internal/motion"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	tb "gopkg.in/tucnak/telebot.v2"
)

// Telebot wraps the bot client with credential information
type Telebot struct {
	bot             *tb.Bot
	adminsID        []int
	subscribedChats []*tb.Chat
	motion          *motion.Motion
}

func New(token string, adminsID []int, subscribedChatsID []int64, motion *motion.Motion) (*Telebot, error) {
	b, err := tb.NewBot(tb.Settings{
		Token:  token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})
	if err != nil {
		return nil, err
	}

	if len(adminsID) < 1 {
		return nil, fmt.Errorf("adminID cannot be empty")
	}
	var ac []*tb.Chat
	for _, acID := range subscribedChatsID {
		ac = append(ac, &tb.Chat{ID: acID})
	}
	metrics.SubscribedChats.Set(float64(len(subscribedChatsID)))

	telebot := &Telebot{
		bot:             b,
		adminsID:        adminsID,
		subscribedChats: ac,
		motion:          motion,
	}

	b.Handle("/subscribe", telebot.subscribeChat)
	b.Handle("/status", telebot.healthCheck)
	b.Handle("/start", telebot.start)
	b.Handle("/pause", telebot.pause)
	b.Handle("/snapshot", telebot.snapshot)

	go b.Start()

	return telebot, nil
}

func (t *Telebot) send(r tb.Recipient, what interface{}) {
	log.Debug().Str("recipient", r.Recipient()).Interface("sendable", what).Msg("sending to telebot")
	if _, err := t.bot.Send(r, what); err != nil {
		log.Error().Err(err).Interface("recipient", r).Interface("sendable", what).Msg("fail to send to telegram bot")
	}
}

// SendToAuthChats sends the the sendable to all subscribed chats
func (t *Telebot) SendToAuthChats(what interface{}) {
	for _, s := range t.subscribedChats {
		t.send(s, what)
	}
}

// SendPicture sends the the picture to all subscribed chats
func (t *Telebot) SendPicture(picturePath string) {
	photo := &tb.Photo{File: tb.FromDisk(picturePath)}
	t.SendToAuthChats(photo)
}

// SendVideo sends the the video to all subscribed chats
func (t *Telebot) SendVideo(videoPath string) {
	video := &tb.Video{File: tb.FromDisk(videoPath)}
	t.SendToAuthChats(video)
}
func (t *Telebot) isAuthorizedSender(m *tb.Message) bool {
	id := m.Sender.ID
	for _, adminID := range t.adminsID {
		if id == adminID {
			return true
		}
	}
	metrics.UnauthorizedRequestsTotal.Inc()
	t.send(m.Sender, "You need admin privileges to execute this command")
	log.Debug().Msgf("unauthorized id request: %d", id)
	return false
}

func (t *Telebot) subscribeChat(m *tb.Message) {
	if !t.isAuthorizedSender(m) {
		return
	}
	t.subscribedChats = append(t.subscribedChats, m.Chat)
	metrics.SubscribedChats.Inc()
	t.send(m.Chat, "chat subscribed")
}
func (t *Telebot) healthCheck(m *tb.Message) {
	metrics.HealthCheckRequestedTotal.Inc()
	t.send(m.Chat, "🟢 Bot")
	resp, err := t.motion.SendCommand(motion.DetectionStatus)
	if err != nil {
		t.send(m.Chat, "🟥 Motion")
		return
	}
	t.send(m.Chat, "🟢 Motion")
	if strings.Contains(resp, "Active") {
		metrics.MovementDetectionActivated.Set(1)
		t.send(m.Chat, "🟢 Detection Active")
		return
	}
	metrics.MovementDetectionActivated.Set(0)
	t.send(m.Chat, "🟢 Detection Paused")

}

func (t *Telebot) start(m *tb.Message) {
	if !t.isAuthorizedSender(m) {
		return
	}
	if _, err := t.motion.SendCommand(motion.DetectionStart); err != nil {
		t.send(m.Chat, "failed to start motion detection")
		return
	}
	metrics.MovementDetectionActivated.Set(1)
	t.send(m.Chat, "motion detection started 👁️")
}
func (t *Telebot) pause(m *tb.Message) {
	if !t.isAuthorizedSender(m) {
		return
	}
	if _, err := t.motion.SendCommand(motion.DetectionPause); err != nil {
		t.send(m.Chat, "failed to pause motion detection")
		return
	}
	metrics.MovementDetectionActivated.Set(0)
	t.send(m.Chat, "motion detection paused 🛑")
}
func (t *Telebot) snapshot(m *tb.Message) {
	if !t.isAuthorizedSender(m) {
		return
	}
	if _, err := t.motion.SendCommand(motion.ActionSnapshot); err != nil {
		t.send(m.Chat, "snapshot commnad failed")
		return
	}
	t.send(m.Chat, "snapshot commnad sent 📸")
}
