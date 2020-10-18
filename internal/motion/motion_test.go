package motion

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatcher(t *testing.T) {
	m := New("", "9999")

	http.Get("http://localhost:9999/send?picture=test.png")
	http.Get("http://localhost:9999/send?video=test.mp4")

	require.Equal(t, 1, len(m.Video))
	require.Equal(t, 1, len(m.Picture))
	assert.Equal(t, "test.mp4", <-m.Video)
	assert.Equal(t, "test.png", <-m.Picture)
}

func TestMotion(t *testing.T) {
	motionMock := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sends the path queried into the body response.
		fmt.Fprintln(w, r.URL.Path)
	}))
	defer motionMock.Close()

	m := New(motionMock.URL, "9999")
	resp, err := m.SendCommand(DetectionStatus)

	assert.NoError(t, err)
	assert.Contains(t, resp, DetectionStatus)
}
