package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"storgage/config"

	"github.com/Azure/azure-storage-blob-go/azblob"
	"github.com/gorilla/mux"
)

type StorageHandler struct {
	containerURL azblob.ContainerURL
	credential   *azblob.SharedKeyCredential
}

func NewStorageHandler(config *config.Config) (*StorageHandler, error) {
	credential, err := azblob.NewSharedKeyCredential(
		config.AzureStorageAccount,
		config.AzureStorageKey,
	)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials: %v", err)
	}

	pipeline := azblob.NewPipeline(credential, azblob.PipelineOptions{})

	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s",
		config.AzureStorageAccount,
		config.AzureStorageContainer))

	containerURL := azblob.NewContainerURL(*URL, pipeline)

	return &StorageHandler{
		containerURL: containerURL,
		credential:   credential,
	}, nil
}

func (h *StorageHandler) generateSasURL(blobURL azblob.BlockBlobURL) (string, error) {
	startTime := time.Now().UTC().Add(-1 * time.Minute)
	expiryTime := startTime.Add(1 * time.Hour)

	permissions := azblob.BlobSASPermissions{Read: true}.String()

	u := blobURL.URL()
	urlParts := azblob.NewBlobURLParts(u)

	sasQueryParams, err := azblob.BlobSASSignatureValues{
		Protocol:      azblob.SASProtocolHTTPS,
		StartTime:     startTime,
		ExpiryTime:    expiryTime,
		Permissions:   permissions,
		ContainerName: urlParts.ContainerName,
		BlobName:      urlParts.BlobName,
	}.NewSASQueryParameters(h.credential)

	if err != nil {
		return "", err
	}

	baseURL := u.String()
	return fmt.Sprintf("%s?%s", baseURL, sasQueryParams.Encode()), nil
}

func (h *StorageHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	blobURL := h.containerURL.NewBlockBlobURL(header.Filename)

	_, err = azblob.UploadStreamToBlockBlob(
		context.Background(),
		file,
		blobURL,
		azblob.UploadStreamToBlockBlobOptions{
			BufferSize: 4 * 1024 * 1024,
			MaxBuffers: 16,
		})

	if err != nil {
		http.Error(w, "Error uploading file", http.StatusInternalServerError)
		return
	}

	sasURL, err := h.generateSasURL(blobURL)
	if err != nil {
		http.Error(w, "Error generating file access URL", http.StatusInternalServerError)
		return
	}

	fileInfo := map[string]interface{}{
		"name":       header.Filename,
		"url":        sasURL,
		"size":       header.Size,
		"uploadDate": time.Now(),
	}

	respondJSON(w, fileInfo)
}

func (h *StorageHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	var files []map[string]interface{}

	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := h.containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			http.Error(w, "Error listing files", http.StatusInternalServerError)
			return
		}

		marker = listBlob.NextMarker

		for _, blobInfo := range listBlob.Segment.BlobItems {
			blobURL := h.containerURL.NewBlobURL(blobInfo.Name)

			sasURL, err := h.generateSasURL(blobURL.ToBlockBlobURL())
			if err != nil {
				continue
			}

			files = append(files, map[string]interface{}{
				"name":       blobInfo.Name,
				"url":        sasURL,
				"size":       *blobInfo.Properties.ContentLength,
				"uploadDate": blobInfo.Properties.LastModified,
			})
		}
	}

	respondJSON(w, files)
}

func (h *StorageHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileName := vars["filename"]
	if fileName == "" {
		http.Error(w, "Filename is required", http.StatusBadRequest)
		return
	}

	blobURL := h.containerURL.NewBlockBlobURL(fileName)

	ctx := context.Background()
	_, err := blobURL.Delete(ctx, azblob.DeleteSnapshotsOptionNone, azblob.BlobAccessConditions{})
	if err != nil {
		http.Error(w, "Could not delete file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "File deleted"})
}

func respondJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
