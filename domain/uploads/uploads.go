package uploads

// TODO: move this to an internal folder or something like it, or just outside the domain folder

import (
	"io"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/google/uuid"
)

const uploads_folder = "public/media/uploads"

// file_ext must include the dot separator, i.e. ".png". Just like the function filepath.Ext returns it
func SaveFile(bucket string, file_ext string, file io.Reader) (string, error) {
	var (
		file_name = uuid.NewString() + file_ext
		file_dir  = path.Join(uploads_folder, bucket)
		file_path = path.Join(file_dir, file_name)
	)

	if err := os.MkdirAll(file_dir, 0750); err != nil {
		log.Printf("failed to create dirs for bucket in %s: %s\n", file_dir, err.Error())
		return "", err
	}

	out, err := os.Create(file_path)

	if err != nil {
		log.Printf("failed to create file in %s: %s\n", file_path, err.Error())
		return "", err
	}

	if _, err = io.Copy(out, file); err != nil {
		log.Printf("failed to copy file contents to %s: %s\n", file_path, err.Error())
		return "", err
	}

	if err := out.Close(); err != nil {
		log.Printf("failed to close the new file in %s: %s\n", file_path, err.Error())
		return "", err
	}

	return file_name, nil
}

func DeleteFile(bucket, filename string) error {
	return os.Remove(GetFilePath(bucket, filename))
}

func GetFilePath(bucket, filename string) string {
	return filepath.Join(uploads_folder, bucket, filename)
}

func GetFilePublicUrl(bucket, filename string) string {
	// TODO: somehow connect the /public endpoint registered on the server with the
	// url generated here. Probably using something like a config value/ENV var
	return "/" + GetFilePath(bucket, filename)
}
