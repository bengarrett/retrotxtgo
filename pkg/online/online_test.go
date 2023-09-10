package online_test

import (
	"encoding/json"
	"testing"

	"github.com/bengarrett/retrotxtgo/pkg/online"
	"github.com/stretchr/testify/assert"
)

func TestEndpoint(t *testing.T) {
	t.Parallel()
	_, p, _ := online.Endpoint("https://demozoo.org/api/v1/productions/126496/",
		`W/"0708012ac3fb439a46dd5156195901b4"`)
	assert.Equal(t, float64(126496), p["id"])
}

func TestPingOk(t *testing.T) {
	t.Parallel()
	ok, _ := online.Ping("https://example.org")
	assert.True(t, ok)
}

func TestPingBad(t *testing.T) {
	t.Parallel()
	bad, _ := online.Ping("https://example.com/this/url/does/not/exist")
	assert.False(t, bad)
}

func TestGetJSON(t *testing.T) {
	t.Parallel()
	_, d, _ := online.Get("https://demozoo.org/api/v1/productions/126496/", "")
	assert.True(t, json.Valid(d), true)
}

func TestGetAPI(t *testing.T) {
	t.Parallel()
	_, d, _ := online.Get(online.ReleaseAPI, "")
	assert.True(t, json.Valid(d), true)
}
