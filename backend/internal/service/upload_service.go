package service

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/openfga/go-sdk/client"
	"authority-app/backend/internal/openfga"
)

type UploadService struct {
	fgaClient *openfga.FGAClient
}

func NewUploadService(fgaClient *openfga.FGAClient) *UploadService {
	return &UploadService{fgaClient: fgaClient}
}

func (s *UploadService) UploadAndExtract(file *multipart.FileHeader, uploaderUsername string) error {
	if filepath.Ext(file.Filename) != ".zip" {
		return errors.New("only zip files are allowed")
	}

	// Save the uploaded zip file to a temporary location
	tempZipPath := filepath.Join(os.TempDir(), file.Filename)
	if err := saveUploadedFile(file, tempZipPath); err != nil {
		return fmt.Errorf("failed to save uploaded file: %w", err)
	}
	defer os.Remove(tempZipPath) // Clean up the temporary zip file

	// Create a temporary directory to extract the zip file
	extractDir, err := os.MkdirTemp("", "uploaded_dir_*")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(extractDir) // Clean up the extracted directory

	// Open the zip file
	r, err := zip.OpenReader(tempZipPath)
	if err != nil {
		return fmt.Errorf("failed to open zip file: %w", err)
	}
	defer r.Close()

	// Extract the files
	for _, f := range r.File {
		fpath := filepath.Join(extractDir, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory for extracted file: %w", err)
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("failed to create extracted file: %w", err)
		}
		defer outFile.Close()

		inFile, err := f.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in zip: %w", err)
		}
		defer inFile.Close()

		if _, err = io.Copy(outFile, inFile); err != nil {
			return fmt.Errorf("failed to copy file from zip: %w", err)
		}
	}

	// Create an owner tuple for the uploaded root directory in OpenFGA
	directoryId := filepath.Base(file.Filename) // Use filename as directory ID for simplicity

	ctx := context.Background()
	_, err = s.fgaClient.Client.Write(ctx).Body(client.ClientWriteRequest{
		Writes: []client.ClientTupleKey{
			{
				User:     "user:" + uploaderUsername,
				Relation: "owner",
				Object:   "directory:" + directoryId,
			},
		},
	}).Execute()
	if err != nil {
		return fmt.Errorf("failed to write OpenFGA tuple: %w", err)
	}

	return nil
}

// saveUploadedFile saves a multipart.FileHeader to disk.
// This is a helper function, similar to gin.Context.SaveUploadedFile.
func saveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

