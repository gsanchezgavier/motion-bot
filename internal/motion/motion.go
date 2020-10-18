package motion

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
)

type MotionCommand string

const (
	DetectionStatus     MotionCommand = "detection/status"     // Return the current status of the camera.
	DetectionConnection MotionCommand = "detection/connection" // Return the connection status of the camera.
	DetectionStart      MotionCommand = "detection/start"      // Start or resume motion detection.
	DetectionPause      MotionCommand = "detection/pause"      // Pause the motion detection.
	ActionEventStart    MotionCommand = "action/eventstart"    // Trigger a new event.
	ActionEventEnd      MotionCommand = "action/eventend"      // Trigger the end of a event.
	ActionSnapshot      MotionCommand = "action/snapshot"      // Create a snapshot
	ActionRestart       MotionCommand = "action/restart"       // Shutdown and restart Motion
	ActionQuit          MotionCommand = "action/quit"          // Close all connections to the camera
	ActionEnd           MotionCommand = "action/end"           // Entirely shutdown the Motion application
)

// Motion wraps functions related to Motion as well as some configuration
type Motion struct {
	controlHost    string
	httpClient     *http.Client
	Video, Picture chan string // subscribe to this channels to get the generated files path from motion
	watcherPort    string
}

// Configures the motion client to send commands and start the event watcher
func New(controlHost, watcherPort string) *Motion {
	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	m := &Motion{
		controlHost: controlHost,
		httpClient:  client,
		watcherPort: watcherPort,
		Video:       make(chan string, 100),
		Picture:     make(chan string, 100),
	}

	m.startWatcher()
	return m
}

// SendCommand Use the predefined command const values in this package to send commands to motion
func (m *Motion) SendCommand(command MotionCommand) (string, error) {
	resp, err := m.httpClient.Get(m.controlHost + "/0/" + string(command))
	if err != nil {
		log.Error().Err(err).Msgf("fail to send command:%s", command)
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		return string(bodyBytes), nil
	}
	log.Error().Msgf("fail to send command:%s", command)
	return "", fmt.Errorf("command failed")
}

// StartWatcher creates an http server on the selected port and listen to motion request
// Motion should be configured to generate a GET request as follows:
// 		# Motion config file
//		on_picture_save curl "http://telebot:8000/send?picture=%f"
// 		on_movie_end curl "http://telebot:8000/send?video=%f"
// %f adds the path to the generated file
func (m *Motion) startWatcher() {
	watcher_server := http.NewServeMux()
	watcher_server.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.URL.Query() {
			if k == "picture" {
				for _, path := range v {
					m.Picture <- path
				}
			}
			if k == "video" {
				for _, path := range v {
					m.Video <- path
				}
			}
		}
	})

	go func() {
		log.Fatal().Err(http.ListenAndServe(":"+m.watcherPort, watcher_server))
	}()
}
