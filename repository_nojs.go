// +build !js

package brbundle

import (
	"archive/zip"
	"fmt"
	"github.com/shibukawa/zipsection"
	"os"
)

func (r *Repository) registerBundle(path string, option Option) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	i, _ := f.Stat()
	reader, err := zip.NewReader(f, i.Size())
	if err != nil {
		fmt.Println(err)
		return err
	}
	b := newPackedBundle(reader, f, option)
	err = b.setDecryptionKey(option.DecryptoKey)
	if err != nil {
		return err
	}
	r.bundles[PackedBundleType] = append(r.bundles[PackedBundleType], b)
	return nil
}

func (r *Repository) registerFolder(path string, encrypted bool, option Option) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if !s.IsDir() {
		return fmt.Errorf("path '%s' for Folder Bundle is not", path)
	}
	f := newFolderBundle(path, encrypted, option)
	err = f.setDecryptionKey(option.DecryptoKey)
	if err != nil {
		return err
	}
	r.bundles[FolderBundleType] = append(r.bundles[FolderBundleType], f)
	return nil
}

func (r *Repository) initFolderBundleByEnvVar() {
}

func (r *Repository) initExeBundle() error {
	var err error
	filepath, err := os.Executable()
	if err != nil {
		return err
	}
	reader, closer, err := zipsection.Open(filepath)
	if err == nil {
		b := newPackedBundle(reader, closer, Option{})
		r.bundles[ExeBundleType] = append(r.bundles[ExeBundleType], b)
	}
	return nil
}
