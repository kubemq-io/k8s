package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpConfig_SetConfig_Empty(t *testing.T) {
	cfg := newTestConfig()
	(&HttpConfig{}).SetConfig(cfg)
	assert.Len(t, vars(cfg), 0)
}

func TestHttpConfig_SetConfig_EmptyCors(t *testing.T) {
	cfg := newTestConfig()
	// An empty CORS list must never be emitted (server keeps its non-empty default).
	(&HttpConfig{Cors: &CorsConfig{AllowOrigins: []string{}}}).SetConfig(cfg)
	assert.Len(t, vars(cfg), 0)
}

func TestHttpConfig_SetConfig_AllFields(t *testing.T) {
	cfg := newTestConfig()
	(&HttpConfig{
		Port:        ptr32(8080),
		ReadTimeout: ptr32(30),
		BodyLimit:   strptr("4M"),
		BaseURL:     strptr("https://kubemq.example.com"),
		Cors: &CorsConfig{
			AllowOrigins:     []string{"https://a", "https://b"},
			AllowMethods:     []string{"GET", "POST"},
			AllowHeaders:     []string{"Content-Type", "Authorization"},
			AllowCredentials: ptrBool(true),
			ExposeHeaders:    []string{"X-Total-Count"},
			MaxAge:           ptr32(3600),
		},
	}).SetConfig(cfg)
	v := vars(cfg)
	require.Len(t, v, 10)
	assert.Equal(t, "8080", v["CONNECTORS_HTTP_PORT"])
	assert.Equal(t, "30", v["CONNECTORS_HTTP_READ_TIMEOUT"])
	assert.Equal(t, "4M", v["CONNECTORS_HTTP_BODY_LIMIT"])
	assert.Equal(t, "https://kubemq.example.com", v["CONNECTORS_HTTP_BASE_URL"])
	assert.Equal(t, "https://a,https://b", v["CONNECTORS_HTTP_CORS_ALLOW_ORIGINS"])
	assert.Equal(t, "GET,POST", v["CONNECTORS_HTTP_CORS_ALLOW_METHODS"])
	assert.Equal(t, "Content-Type,Authorization", v["CONNECTORS_HTTP_CORS_ALLOW_HEADERS"])
	assert.Equal(t, "true", v["CONNECTORS_HTTP_CORS_ALLOW_CREDENTIALS"])
	assert.Equal(t, "X-Total-Count", v["CONNECTORS_HTTP_CORS_EXPOSE_HEADERS"])
	assert.Equal(t, "3600", v["CONNECTORS_HTTP_CORS_MAX_AGE"])
}

func TestHttpConfig_DeepCopy_Independent(t *testing.T) {
	src := &HttpConfig{
		Port:        ptr32(8080),
		ReadTimeout: ptr32(30),
		BodyLimit:   strptr("4M"),
		BaseURL:     strptr("https://x"),
		Cors: &CorsConfig{
			AllowOrigins:     []string{"https://a"},
			AllowCredentials: ptrBool(true),
			MaxAge:           ptr32(3600),
		},
	}
	cp := src.DeepCopy()
	*src.Port = 1
	*src.BodyLimit = "mutated"
	src.Cors.AllowOrigins[0] = "mutated"
	*src.Cors.AllowCredentials = false
	*src.Cors.MaxAge = 1
	src.Cors = nil

	assert.Equal(t, int32(8080), *cp.Port)
	assert.Equal(t, "4M", *cp.BodyLimit)
	require.NotNil(t, cp.Cors)
	assert.Equal(t, "https://a", cp.Cors.AllowOrigins[0])
	assert.Equal(t, true, *cp.Cors.AllowCredentials)
	assert.Equal(t, int32(3600), *cp.Cors.MaxAge)
}

// CorsConfig is also reused by the REST connector; verify the prefix wiring there.
func TestRestConfig_SetConfig_Cors(t *testing.T) {
	cfg := newTestConfig()
	(&RestConfig{
		ReadTimeout:  ptr32(15),
		WriteTimeout: ptr32(15),
		Cors: &CorsConfig{
			AllowOrigins: []string{"*"},
			MaxAge:       ptr32(600),
		},
	}).SetConfig(cfg)
	v := vars(cfg)
	assert.Equal(t, "15", v["CONNECTORS_REST_READ_TIMEOUT"])
	assert.Equal(t, "15", v["CONNECTORS_REST_WRITE_TIMEOUT"])
	assert.Equal(t, "*", v["CONNECTORS_REST_CORS_ALLOW_ORIGINS"])
	assert.Equal(t, "600", v["CONNECTORS_REST_CORS_MAX_AGE"])
}
