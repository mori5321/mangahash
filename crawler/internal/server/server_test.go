package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var conf AppConfig = AppConfig{
	AppPort: 5051,
}

func TestMain(m *testing.M) {
	ctx := context.Background()

	readyCh := make(chan bool)
	go func() {
		App(ctx, conf, readyCh)
	}()

	<-readyCh

	code := m.Run()

	os.Exit(code)
}

var appHost string = fmt.Sprintf("http://localhost:%d", conf.AppPort)

func TestHealthzHandler(t *testing.T) {
	url := appHost + "/healthz"

	res, err := http.Get(url)
	require.NoError(t, err)

	require.Equal(t, http.StatusOK, res.StatusCode)
}
