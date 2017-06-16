package archiver

import (
	"archive/tar"
	"fmt"
	"os"
	"strings"

	"github.com/dsnet/compress/bzip2"
)

// TarBz2 is for TarBz2 format
var TarBz2 tarBz2Format

func init() {
	RegisterFormat("TarBz2", TarBz2)
}

type tarBz2Format struct{}

func (tarBz2Format) Match(filename string) bool {
	return strings.HasSuffix(strings.ToLower(filename), ".tar.bz2") ||
		strings.HasSuffix(strings.ToLower(filename), ".tbz2") ||
		isTarBz2(filename)
}

// isTarBz2 checks the file has the bzip2 compressed Tar format header by
// reading its beginning block.
func isTarBz2(tarbz2Path string) bool {
	f, err := os.Open(tarbz2Path)
	if err != nil {
		return false
	}
	defer f.Close()

	bz2r, err := bzip2.NewReader(f, nil)
	if err != nil {
		return false
	}
	defer bz2r.Close()

	buf := make([]byte, tarBlockSize)
	n, err := bz2r.Read(buf)
	if err != nil || n < tarBlockSize {
		return false
	}

	return hasTarHeader(buf)
}

// Make creates a .tar.bz2 file at tarbz2Path containing
// the contents of files listed in filePaths. File paths
// can be those of regular files or directories. Regular
// files are stored at the 'root' of the archive, and
// directories are recursively added.
func (tarBz2Format) Make(tarbz2Path string, filePaths []string, opts ...OpOption) error {
	ret := Op{verbose: false}
	ret.applyOpts(opts)

	out, err := os.Create(tarbz2Path)
	if err != nil {
		return fmt.Errorf("error creating %s: %v", tarbz2Path, err)
	}
	defer out.Close()

	bz2Writer, err := bzip2.NewWriter(out, nil)
	if err != nil {
		return fmt.Errorf("error compressing %s: %v", tarbz2Path, err)
	}
	defer bz2Writer.Close()

	tarWriter := tar.NewWriter(bz2Writer)
	defer tarWriter.Close()

	return tarball(filePaths, tarWriter, tarbz2Path, ret.verbose)
}

// Open untars source and decompresses the contents into destination.
func (tarBz2Format) Open(source, destination string, opts ...OpOption) error {
	ret := Op{verbose: false}
	ret.applyOpts(opts)

	f, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("%s: failed to open archive: %v", source, err)
	}
	defer f.Close()

	bz2r, err := bzip2.NewReader(f, nil)
	if err != nil {
		return fmt.Errorf("error decompressing %s: %v", source, err)
	}
	defer bz2r.Close()

	return untar(tar.NewReader(bz2r), destination, ret.verbose)
}