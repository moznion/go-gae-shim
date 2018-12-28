package gaeshim

// ref: https://cloud.google.com/appengine/docs/standard/python/config/appref#handlers_element
type handler struct {
	URL         string            `yaml:"url"`
	Script      string            `yaml:"script"`
	StaticDir   string            `yaml:"static_dir"`
	MimeType    string            `yaml:"mime_type"`
	HTTPHeaders map[string]string `yaml:"http_headers"`
}

type configuration struct {
	EnvVariables map[string]string `yaml:"env_variables,omitempty"`
	Handlers     []*handler        `yaml:"handlers"`
}
