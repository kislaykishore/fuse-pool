package fuse

import (
	"syscall"
	"unsafe"
)

func writev(fd int, packet [][]byte) (n int, err error) {
	var stackIovecs [128]syscall.Iovec
	var iovecs []syscall.Iovec
	if len(packet) <= len(stackIovecs) {
		iovecs = stackIovecs[:0]
	} else {
		iovecs = make([]syscall.Iovec, 0, len(packet))
	}
	for _, v := range packet {
		if len(v) == 0 {
			continue
		}
		vec := syscall.Iovec{
			Base: &v[0],
		}
		vec.SetLen(len(v))
		iovecs = append(iovecs, vec)
	}
	if len(iovecs) == 0 {
		return 0, nil
	}
	n1, _, e1 := syscall.Syscall(
		syscall.SYS_WRITEV,
		uintptr(fd), uintptr(unsafe.Pointer(&iovecs[0])), uintptr(len(iovecs)),
	)
	n = int(n1)
	if e1 != 0 {
		err = syscall.Errno(e1)
	}
	return
}
