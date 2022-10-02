package adb

import "os"

// The status information extracted from man 2 stat

// S_IFMT   0170000  /* type of file */

// S_IFIFO  0010000  /* named pipe (fifo) */
// S_IFCHR  0020000  /* character special */
// S_IFDIR  0040000  /* directory */
// S_IFBLK  0060000  /* block special */

// S_IFREG  0100000  /* regular */
// S_IFLNK  0120000  /* symbolic link */
// S_IFSOCK 0140000  /* socket */
// S_IFWHT  0160000  /* whiteout */

// S_ISUID  0004000  /* set user id on execution */
// S_ISGID  0002000  /* set group id on execution */
// S_ISVTX  0001000  /* save swapped text even after use */
// S_IRUSR  0000400  /* read permission, owner */
// S_IWUSR  0000200  /* write permission, owner */
// S_IXUSR  0000100  /* execute/search permission, owner */

/**
 * @link https://github.com/codeskyblue/fa
 */

const (
	ModeDir        uint32 = 0040000
	ModeSocket            = 0140000
	ModeSymlink           = 0120000
	ModeRegular           = 0100000
	ModeNamedPipe         = 0010000
	ModeCharDevice        = 0020000
	ModeSetuid            = 0004000
	ModeSetgid            = 0002000
	ModePerm              = 0000777
)

func maskMatch(m uint32, mask uint32) bool {
	return m&mask == mask
}

var _modeMatches = map[uint32]os.FileMode{
	ModeDir:        os.ModeDir,
	ModeSocket:     os.ModeSocket,
	ModeSymlink:    os.ModeSymlink,
	ModeNamedPipe:  os.ModeNamedPipe,
	ModeCharDevice: os.ModeCharDevice,
	ModeSetuid:     os.ModeSetuid,
	ModeSetgid:     os.ModeSetgid,
}

func fileModeFromAdb(m uint32) os.FileMode {
	mode := os.FileMode(m & ModePerm)
	for statMask, modeMask := range _modeMatches {
		if m&statMask == statMask {
			mode |= modeMask
		}
	}
	return os.FileMode(mode)
}

func fileModeToAdb(mode os.FileMode) uint32 {
	m := uint32(mode) & ModePerm
	if mode.IsRegular() {
		m |= ModeRegular
	}
	for statMask, mask := range _modeMatches {
		if mode&mask == mask {
			m |= statMask
		}
	}
	return m
}
