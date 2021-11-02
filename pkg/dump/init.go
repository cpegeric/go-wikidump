package dump

type DumpConfig struct {
	BaseURL   string
	Ver       string
	Directory string
}

type dump struct {
	BaseURL   string
	Ver       string
	Directory string
}

var base dump = dump{
	BaseURL:   "https://dumps.wikimedia.org",
	Ver:       "/enwiki/20210720/",
	Directory: "./wikipedia-dump/",
}

// Creates a new dump struct. If no paramters are passed the dump is created with the default parameters. Only
// Pass config paramters that you want to change from the default.
// Default values:
// BaseURL   https://dumps.wikimedia.org
// Ver       /enwiki/20210720/
// Directory ./wikipedia-dump/
func New(dumpConfig DumpConfig) *dump {
	result := base
	if dumpConfig.BaseURL != "" {
		result.BaseURL = dumpConfig.BaseURL
	}
	if dumpConfig.Ver != "" {
		result.Ver = dumpConfig.Ver
	}
	if dumpConfig.Directory != "" {
		result.Directory = dumpConfig.Directory
	}
	return &result
}
