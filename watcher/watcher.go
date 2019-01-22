/* 2019-01-20 (cc) <paul4hough@gmail.com>
   application entry point
*/
package main

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	sc "syscall"
	"unsafe"

	"github.com/pahoughton/log-alert/config"

)

type DirWatchMap	map[string]WatchDir
type WatchNotifyMap	map[int]interface{}

type Notice struct {
	sc.InotifyEvent
	Name string
}

type Watcher struct {
	WatchFd		int
	DirMap		DirWatchMap
	WatchMap	WatchNotifyMap
}


func NewWatcher(c *config.Config) (*Watcher,error) {

	w := &Watcher{
		DirWatch:	make(DirWatchMap)
		WatchMap:	make(WatchNotifyMap)
	}

	w.Watchfd, err := sc.InotifyInit1(sc.IN_CLOEXEC)
	if err != nil {
		panic(err)
	}

	for _, clf := range c.LogFiles {

		wf,err := watcher.NewWatchFile(clf)
		if err != nil {
			panic(err)
		}

		dir := path.Dir(wf.Path)
		wdir, ok := w.DirWatch[dir]
		if ! ok {
			wdir := NewWatchDir(dir)
			if err := wdir.AddFile(wf); err != nil {
				panic(err)
			}
			w.DirWatch[dir] = wdir
			dinfo, err := os.Stat(wdir)

			if err != nil {
				panic(fmt.Errorf("stat %s - %s",dir,err.Error()))

			} else if dinfo.IsDir != true  {
				panic(fmt.Errorf("%s is not dir of path %s",dir,wf.Path))
			}
			wdir.Watch, err := sc.InotifyAddWatch(
				w.Watchfd,
				dir,
				sc.IN_MOVED_FROM|
					sc.IN_MOVED_TO|
					sc.IN_DELETE|
					sc.IN_CREATE|
					sc.IN_IGNORED|
					sc.IN_DELETE_SELF|
					sc.IN_MOVE_SELF)
			if err != nil {
				panic(err)
			}
			w.WatchMap[wdir.Watch] = wdir
		} else {
			if err := wdir.AddFile(wf); err != nil {
				panic(err)
			}
		}
		finfo, err := os.Stat(wf.Path)
		if err != nil {
			continue
		}
		if err := w.AddWatch(wf); err != nil {
			panic(err)
		}
	}
}

func (w *Watcher) AddWatch(wf *WatchFile) error {
	wf.Scan()
	wf.Watch, err := sc.InotifyAddWatch(
		w.Watchfd,
		wf.Path,
		sc.IN_MODIFY|
			sc.IN_CLOSE_WRITE|
			sc.IN_DELETE_SELF|
			sc.IN_MOVE_SELF|
			sc.IN_IGNORED)
	if err != nil {
		return(err)
	}
	w.WatchMap[wf.Watch] = wf
	return nil
}

func (w *Watcher) RmWatch(wf *WatchFile) {
	if wf.Watch > -1 {
		delete(w.WatchMap,wf.Watch)
		sc.InotifyRmWatch(
			w.Watchfd,
			wf.Watch,
			sc.IN_ALL_EVENTS)
		wf.Watch = -1
	}
	wf.Reset()
}

func (w *Watcher) Run() {

	buf := make([]byte, (sc.SizeofInotifyEvent+sc.NAME_MAX+1))

	for {
		rn, err := sc.Read(w.WatchFd, buf)
		if err != nil {
			panic(err)
		}
		if rn < sc.SizeofInotifyEvent {
			panic(fmt.Errorf("read size: %d",rn))
		}
		notice := Notice{*(*sc.InotifyEvent)(unsafe.Pointer(buf)),""}

		if notice.Len > 0 {
			b :=(*[sc.NAME_MAX]byte)(
				unsafe.Pointer(&buf[sc.SizeofInotifyEvent]))
			notice.Name = strings.TrimRight(string(b[0:notice.Len]),"\000")
		} else {
			notice.Name = ""
		}

		if w.Debug {
			fmt.Println(notice.String())
		}
		switch reflect.TypeOf(w.WatchMap[notice.wd]) {

		case reflect.TypeOf((*WatchDir)(nil)):
			wdir := (*WatchDir)(w.WatchMap[notice.wd])

			if notice.Mask & sc.IN_DELETE == sc.IN_DELETE ||
				notice.Mask & sc.IN_MOVED_FROM == sc.IN_MOVED_FROM {
				wf, ok := wdir.WatchFileMap[notice.Name]
				if ok {
					w.RmWatch(wf)
				}
			}

			if notice.Mask & sc.IN_CREATE == sc.IN_CREATE ||
				notice.Mask & sc.IN_MOVED_TO == sc.IN_MOVED_TO {
				wf, ok := wdir.WatchFileMap[notice.Name]
				if ok {
					w.RmWatch(wf)
					w.AddWatch(wf)
				}
			}
			if notice.Mask & sc.IN_IGNORED == sc.IN_IGNORED ||
				notice.Mask & sc.IN_DELETE_SELF == sc.IN_DELETE_SELF ||
				notice.Mask & sc.IN_MOVE_SELF == sc.IN_MOVE_SELF {
				panic(fmt.Errorf("%s dir deleted or moved",wdir.Path))
			}

		case reflect.TypeOf((*WatchFile)(nil)):
			wf := (*WatchFile)(w.WatchMap[notice.wd])
			if notice.Mask & sc.IN_DELETE_SELF == sc.IN_DELETE ||
				notice.Mask & sc.IN_MOVED_SELF == sc.IN_MOVED_SELF ||
				notice.Mask & sc.IN_IGNORED == sc.IN_IGNORED {

				w.RmWatch(wf)
			} else {
				wf.Scan()
			}

		default:
			panic(fmt.Errorf("ERROR: unknown watch type"))
		}
	}
}

func (n *Notice) String() string {
	result := "event"
	if len(n.Name) > 0 {
		result = fmt.Sprintf("%v with path %v and mask", result, n.Name)
	} else {
		result = fmt.Sprintf("%v with mask", result)
	}
	if n.Mask & sc.IN_ACCESS == sc.IN_ACCESS {
		result = fmt.Sprintf("%v IN_ACCESS", result)
	}
	if n.Mask & sc.IN_ATTRIB == sc.IN_ATTRIB {
		result = fmt.Sprintf("%v IN_ATTRIB", result)
	}
	if n.Mask & sc.IN_CLOSE_WRITE == sc.IN_CLOSE_WRITE {
		result = fmt.Sprintf("%v IN_CLOSE_WRITE", result)
	}
	if n.Mask & sc.IN_CLOSE_NOWRITE == sc.IN_CLOSE_NOWRITE {
		result = fmt.Sprintf("%v IN_CLOSE_NOWRITE", result)
	}
	if n.Mask & sc.IN_CREATE == sc.IN_CREATE {
		result = fmt.Sprintf("%v IN_CREATE", result)
	}
	if n.Mask & sc.IN_DELETE == sc.IN_DELETE {
		result = fmt.Sprintf("%v IN_DELETE", result)
	}
	if n.Mask & sc.IN_DELETE_SELF == sc.IN_DELETE_SELF {
		result = fmt.Sprintf("%v IN_DELETE_SELF", result)
	}
	if n.Mask & sc.IN_MODIFY == sc.IN_MODIFY {
		result = fmt.Sprintf("%v IN_MODIFY", result)
	}
	if n.Mask & sc.IN_MOVE_SELF == sc.IN_MOVE_SELF {
		result = fmt.Sprintf("%v IN_MOVE_SELF", result)
	}
	if n.Mask & sc.IN_MOVED_FROM == sc.IN_MOVED_FROM {
		result = fmt.Sprintf("%v IN_MOVED_FROM", result)
	}
	if n.Mask & sc.IN_MOVED_TO == sc.IN_MOVED_TO {
		result = fmt.Sprintf("%v IN_MOVED_TO", result)
	}
	if n.Mask & sc.IN_OPEN == sc.IN_OPEN {
		result = fmt.Sprintf("%v IN_OPEN", result)
	}
	if n.Mask & sc.IN_IGNORED == sc.IN_IGNORED {
		result = fmt.Sprintf("%v IN_IGNORED", result)
	}
	if n.Mask & sc.IN_ISDIR == sc.IN_ISDIR {
		result = fmt.Sprintf("%v IN_ISDIR", result)
	}
	if n.Mask & sc.IN_Q_OVERFLOW == sc.IN_Q_OVERFLOW {
		result = fmt.Sprintf("%v IN_Q_OVERFLOW", result)
	}
	if n.Mask & sc.IN_UNMOUNT == sc.IN_UNMOUNT {
		result = fmt.Sprintf("%v IN_UNMOUNT", result)
	}
	return result
}
