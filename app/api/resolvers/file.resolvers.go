package resolvers

import (
	"context"
	"fmt"
	"io/ioutil"
	"orijinplus/app/models"

	"github.com/99designs/gqlgen/graphql"
	"github.com/h2non/filetype"
)

func (r *mutationResolver) FileUpload(ctx context.Context, file graphql.Upload) (*models.File, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.CreateSKU, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	obj, err := r.Upload(ctx, file, auther)
	if err != nil {
		return nil, err
	}

	return obj, nil
}

func (r *mutationResolver) FileUploadMultiple(ctx context.Context, files []graphql.Upload) ([]models.File, error) {
	auther, authErr := r.GetAuther(ctx)
	if authErr != nil {
		return nil, authErr
	}
	if err := r.services.AuthService.GrantPermission(ctx, auther, models.CreateSKU, true, false); err != nil {
		return nil, fmt.Errorf(err.Message)
	}

	objects := []models.File{}

	for _, file := range files {
		obj, err := r.Upload(ctx, file, auther)
		if err != nil {
			return nil, err
		}

		objects = append(objects, *obj)
	}

	return objects, nil
}

// Upload adds a new blob to the db
func (r *mutationResolver) Upload(ctx context.Context, file graphql.Upload, auther *models.Auther) (*models.File, error) {
	// Read file data
	fileObj, err := ioutil.ReadAll(file.File)
	if err != nil {
		return nil, fmt.Errorf("file upload - read file")
	}

	// get mime type
	kind, err := filetype.Match(fileObj)
	if err != nil {
		return nil, fmt.Errorf("file upload - get mime type")
	}

	if kind == filetype.Unknown {
		return nil, fmt.Errorf("file upload - image type is unknown")
	}

	// mimeType := kind.MIME.Value
	// extension := kind.Extension

	obj, uploadErr := r.filestore.UploadFile(file.Filename, fileObj)
	if uploadErr != nil {
		return nil, fmt.Errorf(uploadErr.Message)
	}

	return obj, nil
}
