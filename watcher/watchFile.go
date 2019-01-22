/* 2019-01-20 (cc) <paul4hough@gmail.com>
   a watch file is a file being monitored for alert conditions

*/
package watcher

import (
	"path/filepath"
	"regexp"

	"github.com/pahoughton/log-alert/config"
)

type Pattern struct {
	Exp		*Regexp
	Labels	map[string]string
	Annots	map[string]string
}

type WatchFile struct {
	WatchId	int
	Path	string
	Offset	uint64
	Line	uint64
	Labels	map[string]string
	Annots	map[string]string
	Pats	[]Pattern
}

func NewWatchFile(clf *config.LogFile) (*WatchFile, error) {

	wf := &WatchFile{
		WatchId:	-1,
		Path:		filepath.Abs(clf.Path),
		Offset:		0,
		Line:		0,
		Labels:		clf.Labels,
		Annots:		clf.Annots,
		Pats:		make([]Pattern),
	}

	for _, pat := range clf.Pats {
		re, err := regexp.MustCompile(pat.Regex)
		if err != nil {
			return nil, fmt.Errorf("pat compile %s for %s",clf.Path,pat.Regex)
		}
		p := &Pattern{
			Exp:	re,
			Labels:	pat.Labels,
			Annots:	pat.Annots,
		}
		wf.Pats = append(wf.Pats,p)
	}
	return wf, nil
}

const (
	LineSep = []byte{'\n'}
)
func (wf *WatchFile) Scan() {

	f, err := os.Open(wf.Path)
	if err != nil {
		amgr.AlertFile(wf,err.Error())
		return
	}
	finfo, err := f.Stat()
	if err != nil {
		amgr.AlertFile(wf,err.Error())
		return
	}
	if wf.Offset == finfo.Size() {
		return
	} else if wf.Offset > finfo.Size() {
		wf.Offset = 0
		wf.Lines = 0
	} else if wf.Offset > 0 && wf.Offset < finfo.Size() {
		f.Seek(wf.Offset, 0)
	}

	/* fixme
found by pattern map
alert only at end of searches
include matched text
one alert per ????
match label set
label set merge

*/
	flines := make([]uint64)
	buf := make([]byte, 4096 * 8)
	pos := 0

	for n, err := os.Read(buf[pos:]); err == nil && n > 0 {

		for lpos := bytes.Index(buf[pos:n],LineSep); lpos > -1 {
			wf.Lines += 1
			for _, p := range wf.Pats {
				if p.Find(buf[pos:lpos]) != nil {
					if len(p.Labels) > 0 || len(p.Annots) > 0 {
						amgr.AlertPat(wf,p)
					} else {
						flines = append(flines,wf.Lines)
					}
				}
			}
			pos += lpos+1
		}
		if pos == n {
			pos = 0
		} else {
			copy(buf, buf[pos:n])
			pos = n - pos
		}
	}
	if len(flines) > 0 {
		amgr.AlertLines(wf,flines)
	}
}
