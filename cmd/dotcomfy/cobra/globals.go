package cobra

import (
	"dotcomfy/internal/config"

	"github.com/charmbracelet/log"
)

var (
	BRANCH    string
	COMMIT    string
	CONFIRM   bool
	CFG_FILE  string
	CONFIG    config.Config
	REPO      string
	VERBOSITY int
	LOGGER    *log.Logger
)
