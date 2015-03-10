package setting

var (
	// Application settings.
	AppVer   string
	ProdMode bool
	HTTPPort int
	DataDir  = "./data"
)

func init() {
	// Dummy settings for now.
	HTTPPort = 4000
}
