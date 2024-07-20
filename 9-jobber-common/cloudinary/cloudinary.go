package cloudinary

import (
	"context"

	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type CloudinaryHandler struct {
	cld *cloudinary.Cloudinary
}

func New(cloud, key, secret string) (*CloudinaryHandler, error) {
	cld, err := cloudinary.NewFromParams(cloud, key, secret)
	if err != nil {
		return nil, err
	}

	cld.Config.URL.Secure = true

	return &CloudinaryHandler{cld}, nil
}

func (ch *CloudinaryHandler) UploadImage(ctx context.Context, file string, public_id string, overwrite bool, invalidate bool) (string, error) {
	// Upload the image.
	resp, err := ch.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       public_id,
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(overwrite),
		Invalidate:     api.Bool(invalidate),
		ResourceType:   "auto",
	})
	if err != nil {
		return "", nil
	}

	return resp.SecureURL, nil
}

func (ch *CloudinaryHandler) UploadVideo(ctx context.Context, file string, public_id string, overwrite bool, invalidate bool) (string, error) {
	// Upload the video.
	resp, err := ch.cld.Upload.Upload(ctx, file, uploader.UploadParams{
		PublicID:       public_id,
		UniqueFilename: api.Bool(false),
		Overwrite:      api.Bool(overwrite),
		Invalidate:     api.Bool(invalidate),
		ResourceType:   "video",
	})
	if err != nil {
		return "", nil
	}

	return resp.SecureURL, nil
}
