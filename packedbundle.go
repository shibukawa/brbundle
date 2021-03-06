package brbundle

import (
	"archive/zip"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type packedBundle struct {
	baseBundle
	reader        *zip.Reader
	closer        io.Closer
	localFilePath string
}

func newPackedBundle(reader *zip.Reader, closer io.Closer, option Option) *packedBundle {
	reader.RegisterDecompressor(ZIPMethodSnappy, snappyDecompressor)
	sort.Slice(reader.File, func(i, j int) bool {
		return reader.File[i].Name < reader.File[j].Name
	})
	mountPoint := option.MountPoint
	if mountPoint != "" && !strings.HasSuffix(mountPoint, "/") {
		mountPoint = mountPoint + "/"
	}
	result := &packedBundle{
		baseBundle: baseBundle{
			mountPoint:    mountPoint,
			name:          option.Name,
			decryptorType: reader.Comment[:1],
		},
		reader: reader,
		closer: closer,
	}
	return result
}

func (p packedBundle) find(searchPath string) (FileEntry, error) {
	i := sort.Search(len(p.reader.File), func(i int) bool {
		return p.reader.File[i].Name >= searchPath
	})
	if i < len(p.reader.File) && p.reader.File[i].Name == searchPath {
		decryptor, err := p.baseBundle.getDecryptor()
		if err != nil {
			return nil, err
		}
		file := p.reader.File[i]
		var decompressor Decompressor
		switch file.Comment[0:1] {
		case UseBrotli:
			decompressor = brotliDecompressor
		case NotToCompress:
			decompressor = nullDecompressor
		}
		return &packedFileEntry{
			file:         file,
			decompressor: decompressor,
			decryptor:    decryptor,
			logicalPath:  path.Clean("/" + path.Join(p.baseBundle.mountPoint, file.Name)),
		}, nil
	} else {
		return nil, nil
	}
}

func (packedBundle) readdir(path string) []FileEntry {
	panic("implement me")
}

func (p *packedBundle) close() {
	if p.closer != nil {
		p.closer.Close()
	}
}

type packedFileEntry struct {
	file         *zip.File
	logicalPath  string
	decompressor Decompressor
	decryptor    Decryptor
}

func (p packedFileEntry) Reader() (io.ReadCloser, error) {
	reader, err := p.file.Open()
	if err != nil {
		return nil, err
	}
	decryptoReader, err := p.decryptor.Decrypto(reader)
	if err != nil {
		return nil, err
	}
	return NewReadCloser(p.decompressor(decryptoReader), reader), nil
}

func (p packedFileEntry) ReadAll() ([]byte, error) {
	reader, err := p.Reader()
	if err != nil {
		return nil, err
	}
	defer reader.Close()
	return ioutil.ReadAll(reader)
}

func (p packedFileEntry) BrotliReader() (io.ReadCloser, error) {
	if p.file.Comment[0:1] == UseBrotli {
		reader, err := p.file.Open()
		if err != nil {
			return nil, err
		}
		decryptoReader, err := p.decryptor.Decrypto(reader)
		if err != nil {
			return nil, err
		}
		return NewReadCloser(decryptoReader, reader), nil
	}
	return nil, errors.New("Source data is not compressed by brotli")
}

func (p packedFileEntry) CompressedSize() int64 {
	if p.file.Comment[0:1] == UseBrotli {
		return int64(p.file.CompressedSize64)
	}
	return -1
}

func (p packedFileEntry) Stat() os.FileInfo {
	_, etag, _ := ParseCommentString(p.file.Comment)
	sizePart := strings.Split(etag, "-")[0]
	originalSize, _ := strconv.ParseInt(sizePart, 16, 64)
	return &zipFileInfo{
		name:         p.Name(),
		modTime:      p.file.Modified,
		originalSize: originalSize,
	}
}

func (p packedFileEntry) Name() string {
	return path.Base(p.file.Name)
}

func (p packedFileEntry) Path() string {
	return p.logicalPath
}

func (p packedFileEntry) EtagAndContentType() (string, string) {
	_, etag, contentType := ParseCommentString(p.file.Comment)
	return etag, contentType
}

type zipFileInfo struct {
	name         string
	modTime      time.Time
	originalSize int64
}

func (z zipFileInfo) Name() string {
	return z.name
}

func (z zipFileInfo) Size() int64 {
	return z.originalSize
}

func (z zipFileInfo) Mode() os.FileMode {
	return 0444
}
func (z zipFileInfo) ModTime() time.Time {
	return z.modTime
}

func (z zipFileInfo) IsDir() bool {
	return false
}

func (z zipFileInfo) Sys() interface{} {
	return nil
}
