package setting

var (
	// Application settings.
	AppVer   string
	ProdMode bool
	HTTPPort int
)

func init() {
	// Dummy settings for now.
	HTTPPort = 4000
}
