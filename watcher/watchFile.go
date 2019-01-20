/* 2019-01-20 (cc) <paul4hough@gmail.com>
   a watch file is a file being monitored for alert conditions

*/
package watcher

type Pattern struct {
	Exp		*Regexp
	Labels	map[string]string
	Annots	map[string]string
}

type WatchFile struct {
	WatchId	int
	Path	string
	Offset	uint64
	Labels	map[string]string
	Annots	map[string]string
	Pats	[]Pattern
}

func NewWatchFile(fpath	string)
) {

	wf := &WatchFile{
		WatchId:	-1,
		Path:		fpath,
		Offset:		0,
	}
	return wf
}
