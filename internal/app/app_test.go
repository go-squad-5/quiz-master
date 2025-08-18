package app

import (
	"testing"

	"github.com/go-squad-5/quiz-master/internal/config"
)

func Test_NewApp(t *testing.T) {
	tests := []struct {
		name   string
		config config.Config
	}{
		{
			name: "valid config",
			config: config.Config{
				Port:        ":8090",
				DSN:         "jjj",
				WorkerCount: 9,
			},
		},
    {
      name: "invalid config"
    }
	}
	app := NewApp()
}
