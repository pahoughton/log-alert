/* 2019-01-20 (cc) <paul4hough@gmail.com>
   application entry point
*/
package main

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"
)

type WatchFileMap map[string]WatchFile

type WatchDir struct {
	Watch		int
	Path		string
	WatchFile	WatchFileMap
}

type DirWatchMap	map[string]WatchDir
type WatchMap		map[int]interface{}

type Watcher struct {
	WatchFd		int
	DirWatch	DirWatchMap
	DirMissing	map[stirng]WatchDir
	Watches		WatchMap
}

func (w *WatchFile) AlertString(s string) {
	fmt.Println("ALERT: %s %s",w.Path,s)
}

func NewWatcher() (*Watcher,error) {

	w := &Watcher{
		DirWatch:	DirWatchMap{}
		DirMissing:	[]string{}
		Watch		WatchMap{}
	}

	watchfd, err := syscall.InotifyInit1(syscall.IN_CLOEXEC)
	if err != nil {
		panic(err)
	}

	watchPathList := []string{
		"/var/log/auth.log",
		"/var/log/dpkg.log",
		"/var/log/kern.log",
		"/var/log/mail.log",
		"/var/log/syslog",
		"/home/paul/test/log-alert/test.log",
		"/home/paul/test/log-alert/subdir/sub-test.log",
	}

	for _, fpath := range watchPathList {

		wf := &WatchFile{
			Path:	fpath,
			Offset:	0,
			Watch:	-1,
		}
		// watch directory for changes
		dir := path.Dir(filepath.Abs(fpath))
		wdir, ok := w.DirWatch[dir]
		if ! ok {
			wdir := &WatchDir{
				Path:		wdir,
				Watch:		-1,
				WatchFile:	WatchFileMap{},
			}
			wdir.WatchFiles[path.Base(fpath)] = wf
			w.DirWatch[dir] = wdir
			dinfo, err := os.Stat(wdir)
			if err != nil {
				w.MissDirs[dir] = wdir
				wf.AlertString(fmt.Srintf("path: %s dir %s: %s",
					fpath,dir,err.Error()))
				continue
			} else if dinfo.IsDir != true  {
				w.MissDirs[dir] = wdir
				wf.AlertString(fmt.Sprintf("path %s dir %s: not dir",
					fpath,dir))
				continue
			} else {
				watch, err := syscall.InotifyAddWatch(
					watchfd,
					wdir,
					syscall.IN_MODIFY|
						syscall.IN_MOVED_FROM|
						syscall.IN_MOVED_TO|
						syscall.IN_DELETE|
						syscall.IN_CREATE)
				if err != nil {
					panic(err)
				}
				wdir.Watch = watch
				WatchMap[watch] = wdir
			}
		} else {
			wdir.WatchFiles[path.Base(fpath)] = wf
		}
		if wdir.Watch < 0 {
			wf.AlertString("dir not watched for: "+fpath)
			continue
		}
		finfo, err := os.Stat(fpath)
		if err != nil {
			wf.AlertString(fpath+" stat "+err.Error())
			continue
		}
		wf.Scan()
		watch, err := syscall.InotifyAddWatch(
			watchfd,
			fpath,
			syscall.IN_MODIFY|
				syscall.IN_CLOSE_WRITE|
				syscall.IN_DELETE_SELF|
				syscall.IN_MOVE_SELF)
		if err != nil {
			wf.AlertString("InotifyAddWatch "+err.Error())
		}
		wf.Watch = watch
		WatchMap[watch] = wf
	}
}






					}

		}



				}
			}

			if err != nil {
				continue
			}


		if err != nil {

		wfinfo, err := os.Stat(wpath)



}
