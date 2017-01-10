package ziper

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func TarGz(path []string, dest string) error {
	// file write
	fw, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer fw.Close()
	// gzip write
	gw := gzip.NewWriter(fw)
	defer gw.Close()
	// tar write
	tw := tar.NewWriter(gw)
	defer tw.Close()
	for _, v := range path {
		if e := targzFilesAndDirectory(tw, v, ""); e != nil {
			return e
		}
	}
	return nil
}

func targzFilesAndDirectory(w *tar.Writer, path string, rel string) error {
	//self
	finfo, err := os.Stat(path)
	if err != nil {
		return err
	}
	h, err := tar.FileInfoHeader(finfo, "")
	if err != nil {
		return err
	}
	h.Name = filepath.Join(rel, h.Name)
	err = w.WriteHeader(h)
	if err != nil {
		return err
	}
	if !finfo.IsDir() {
		f, err := os.Open(path)
		if err != nil {
			return err
		}
		_, err = io.Copy(w, f)
		f.Close()
		if err != nil {
			return err
		}
		return nil
	}
	rel = filepath.Join(rel, finfo.Name())
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, finfo := range fs {
		err = targzFilesAndDirectory(w, filepath.Join(path, finfo.Name()), rel)
		if err != nil {
			return err
		}
	}
	return nil
}

func UnTarGzFromReader(reader io.Reader, dest string) error {
	gzreader, err := gzip.NewReader(reader)
	if err != nil {
		return err
	}
	defer gzreader.Close()
	tarreader := tar.NewReader(gzreader)
	for {
		h, err := tarreader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if h.FileInfo().IsDir() {
			os.MkdirAll(filepath.Join(dest, h.Name), h.FileInfo().Mode())
			continue
		}
		fd, err := os.OpenFile(filepath.Join(dest, h.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, h.FileInfo().Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(fd, tarreader)
		fd.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func UnTarGz(src, dest string) error {
	fd, err := os.Open(src)
	if err != nil {
		return err
	}
	defer fd.Close()
	gzreader, err := gzip.NewReader(fd)
	if err != nil {
		return err
	}
	defer gzreader.Close()
	tarreader := tar.NewReader(gzreader)
	for {
		h, err := tarreader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if h.FileInfo().IsDir() {
			os.MkdirAll(filepath.Join(dest, h.Name), h.FileInfo().Mode())
			continue
		}
		fd, err := os.OpenFile(filepath.Join(dest, h.Name), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, h.FileInfo().Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(fd, tarreader)
		fd.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
