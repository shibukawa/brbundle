package main

import (
	"io"
	"os"
	"path/filepath"
	"runtime"

	"github.com/fatih/color"
	"github.com/shibukawa/brbundle"
)

func copyWorker(compressor *Compressor, encryptor *Encryptor, destPath, srcDirPath string, jobs <-chan Entry, wait chan<- struct{}) {
	for entry := range jobs {
		outputPath := filepath.Join(destPath, entry.Path)
		output, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, entry.Info.Mode())
		if err != nil {
			color.Red("write file creation error: %s\n", entry.Path, err.Error())
			continue
		}

		err = processInput(compressor, encryptor, srcDirPath, entry, func(writerTo io.WriterTo, etag string) {
			writerTo.WriteTo(output)
			output.Close()
		})

		if err != nil {
			continue
		}

		os.Chtimes(outputPath, entry.Info.ModTime(), entry.Info.ModTime())
	}
	wait <- struct{}{}
}

func createContentFolder(ctype brbundle.CompressionType, etype brbundle.EncryptionType, ekey []byte, destPath, srcDirPath string) {
	color.Cyan("Content Folder Mode (Compression: %s, Encyrption: %s)\n\n", ctype, etype)

	os.MkdirAll(destPath, 0755)
	paths, dirs, ignored := Traverse(srcDirPath)

	for _, dir := range dirs {
		os.MkdirAll(filepath.Join(destPath, dir.Path), 0755)
	}

	wait := make(chan struct{})
	for i := 0; i < runtime.NumCPU(); i++ {
		go copyWorker(NewCompressor(ctype), NewEncryptor(etype, ekey), destPath, srcDirPath, paths, wait)
	}

	close(paths)

	for i := 0; i < runtime.NumCPU(); i++ {
		<-wait
	}

	for _, path := range ignored {
		color.Yellow("  ignored: %s\n", path)
	}
	color.Cyan("\nDone\n\n")
}
