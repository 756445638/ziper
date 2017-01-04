package ziper

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Targz struct {
	dest string
}

func (z *Targz) TarGz(path, dest string) error {
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
	return z.TargzFilesAndDirectory(tw, path, "")
}

func (z *Targz) TargzFilesAndDirectory(w *tar.Writer, path string, rel string) error {
	fs, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}

	for _, finfo := range fs {
		if finfo.IsDir() { //write directory
			h, err := tar.FileInfoHeader(finfo, "")
			if err != nil {
				return err
			}
			h.Name = filepath.Join(rel, h.Name)
			err = w.WriteHeader(h)
			if err != nil {
				return err
			}
			err = z.TargzFilesAndDirectory(w, filepath.Join(path, finfo.Name()), filepath.Join(rel, finfo.Name()))
			if err != nil {
				return err
			}
			continue
		}
		//files
		h, err := tar.FileInfoHeader(finfo, "")
		if err != nil {
			return err
		}
		h.Name = filepath.Join(rel, finfo.Name())
		err = w.WriteHeader(h)
		if err != nil {
			return err
		}
		//write content
		f, err := os.Open(filepath.Join(path, finfo.Name()))
		if err != nil {
			return err
		}
		_, err = io.Copy(w, f)
		f.Close()
		if err != nil {
			return err
		}
	}
	return nil
}

func (z *Targz) UnTarGzFromReader(reader io.Reader, dest string) error {
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

func (z *Targz) UnTarGz(src, dest string) error {
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
