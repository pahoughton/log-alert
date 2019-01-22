/* 2019-01-21 (cc) <paul4hough@gmail.com>
   load config yml
*/
package config

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type AmgrSConfig struct {
	Targets	[]string	`yaml:"targets"`
}

type Amgr struct {
	Scheme		string		`yaml:"scheme"`
	SConfigs	AmgrSConfig	`yaml:"static-configs"`
}

type GlobalConfig struct {
	ListenAddr	string				`yaml:"listen-addr"`
	Labels		map[string]string	`yaml:"labels,omitempty"`
	Annots		map[string]string	`yaml:"annotations,omitempty"`
}

type Pattern struct {
	Regex	string				`yaml:"regex"`
	Labels	map[string]string	`yaml:"labels,omitempty"`
	Annots	map[string]string	`yaml:"annotations,omitempty"`
}

type LogFile struct {
	Path	string				`yaml:"path"`
	Labels	map[string]string	`yaml:"labels,omitempty"`
	Annots	map[string]string	`yaml:"annotations,omitempty"`
	Pats	[]Pattern			`yaml:"patterns"`
}

type Config struct {
	Global		GlobalConfig	`yaml:"global"`
	LogFiles	[]LogFile		`yaml:"log-files"`
	Amgrs		[]Amgr			`yaml:"alertmanagers"`
}

func LoadFile(fn string) (*Config, error) {

	dat, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.UnmarshalStrict(dat, cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
