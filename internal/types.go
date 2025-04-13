package internal

type Commit struct {
	ID          string
	Message     string
	AuthoreDate string
}

type Config struct {
	APIKey          string `yaml:"api_key"`
	Provider        string `yaml:"provider"`
	Model           string `yaml:"model"`
	Mood            string `yaml:"mood"`
	ReleaseType     string `yaml:"release_type"`
	BulletStyle     string `yaml:"bullet_style"`
	IncludeSections bool   `yaml:"include_sections"`
	Language        string `yaml:"language"`
	Emojis          bool
	Copy            bool
}

type NotesPrompt struct {
	Context []string
	TagName string
	Config
}
