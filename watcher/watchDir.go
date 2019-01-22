/* 2019-01-21 (cc) <paul4hough@gmail.com>
   watch directories
*/
package watcher

import (
	"github.com/pahoughton/log-alert/config"
)

type WatchFileMap map[string]WatchFile

type WatchDir struct {
	Watch	int
	Path	string
	Files	WatchFileMap
}

func New(dir string) *WatchDir {
	wd := &WatchDir{
		Watch:	-1,
		Path:	dir,
		Files:	make(WatchFileMap),
	}
	return wd
}

func (wd *WatchDir) AddFile(wf *WatchFile) error {
	if _, ok := wd.Files[path.Base(wf.Path)]; ok == true {
		return fmt.Errorf("dup file: "+fpath)
	}
	wd.Files[path.Base(wf.Path)] = wf
	return nil
}
