package persist

import "flag"

type Config struct {
	EtcdConfigPath string
	Delay          int
}

var (
	Conf *Config
)

func init() {
	Conf = parseOptions()
}

func parseOptions() *Config {
	delay := flag.Int("delay", 10, "Delay time to extract etcd data")
	etcdConfigPath := flag.String("etcd-config", "etcd-config.json", "Location of the etcd config file")
	flag.Parse()
	return &Config{
		Delay:          *delay,
		EtcdConfigPath: *etcdConfigPath,
	}
}
