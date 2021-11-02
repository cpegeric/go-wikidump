package dump

type Dump struct {
	BaseURL   string
	Ver       string
	Directory string
}

var base Dump = Dump{
	BaseURL:   "https://dumps.wikimedia.org",
	Ver:       "/enwiki/20210720/",
	Directory: "./wikipedia-dump/",
}

// Creates a new dump struct with the default paramters.
func New(params map[string]string) *Dump {
	config := base
	for k, v := range params {
		switch k {
		case "BaseURL":
			config.BaseURL = v
		case "Ver":
			config.Ver = v
		case "Directory":
			config.Directory = v
		}
	}
	return &config
}
