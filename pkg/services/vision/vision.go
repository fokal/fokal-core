package vision

import (
	"context"
	"encoding/base64"

	"github.com/cridenour/go-postgis"
	"github.com/devinmcgloin/clr/clr"
	"github.com/fokal/fokal-core/pkg/domain"
	"github.com/fokal/fokal-core/pkg/logger"
	"github.com/fokal/fokal-core/pkg/services/color"
	"github.com/jmoiron/sqlx"

	"image"

	"bytes"
	"image/jpeg"

	"github.com/nfnt/resize"
	"google.golang.org/api/vision/v1"
)

type VisionService struct {
	db     *sqlx.DB
	vision *vision.Service
}

func New(db *sqlx.DB, vision *vision.Service) *VisionService {
	return &VisionService{db: db, vision: vision}
}

func (vs VisionService) AnnotateImage(ctx context.Context, img image.Image) (*domain.ImageAnnotation, error) {
	m := resize.Resize(300, 0, img, resize.Bilinear)
	buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, m, nil)
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}
	// Construct a text request, encoding the image in base64.

	req := &vision.AnnotateImageRequest{
		// Apply image which is encoded by base64
		Image: &vision.Image{
			Content: base64.StdEncoding.EncodeToString(buf.Bytes()),
		},
		// Apply features to indicate what type of image detection
		Features: []*vision.Feature{
			{Type: "SAFE_SEARCH_DETECTION"},
			{Type: "LANDMARK_DETECTION"},
			{Type: "IMAGE_PROPERTIES"},
			{Type: "LABEL_DETECTION"},
		},
	}

	batch := &vision.BatchAnnotateImagesRequest{
		Requests: []*vision.AnnotateImageRequest{req},
	}

	res, err := vs.vision.Images.Annotate(batch).Do()
	if err != nil {
		logger.Error(ctx, err)
		return nil, err
	}

	r := res.Responses[0]
	rsp := &domain.ImageAnnotation{Safe: true}

	shade := color.NewWithType(vs.db, color.Shade)
	specific := color.NewWithType(vs.db, color.SpecificColor)

	for _, col := range r.ImagePropertiesAnnotation.DominantColors.Colors {
		sRGB := clr.RGB{
			R: uint8(col.Color.Red),
			G: uint8(col.Color.Green),
			B: uint8(col.Color.Blue)}

		h, s, v := sRGB.HSV()
		rsp.ColorProperties = append(rsp.ColorProperties, domain.Color{
			SRGB:          sRGB,
			PixelFraction: col.PixelFraction,
			Score:         col.Score,
			Hex:           sRGB.Hex(),
			HSV: clr.HSV{
				H: h, S: s, V: v,
			},
			Shade:     sRGB.ColorName(shade),
			ColorName: sRGB.ColorName(specific),
		})
	}

	for _, likelihood := range []string{"POSSIBLE", "LIKELY", "VERY_LIKELY"} {
		if r.SafeSearchAnnotation.Adult == likelihood {
			rsp.Safe = false
		}
		if r.SafeSearchAnnotation.Violence == likelihood {
			rsp.Safe = false
		}
		if r.SafeSearchAnnotation.Medical == likelihood {
			rsp.Safe = false
		}
		if r.SafeSearchAnnotation.Spoof == likelihood {
			rsp.Safe = false
		}
	}

	unique := make(map[string]bool, len(rsp.Labels))

	for _, label := range r.LabelAnnotations {
		if _, ok := unique[label.Description]; !ok {
			rsp.Labels = append(rsp.Labels, domain.Label{
				Description: label.Description,
				Score:       label.Score,
			})
			unique[label.Description] = true
		}
	}

	for _, landmark := range r.LandmarkAnnotations {
		rsp.Landmark = append(rsp.Landmark, domain.Landmark{
			Description: landmark.Description,
			Score:       landmark.Score,
			Location: postgis.PointS{
				SRID: 4326,
				X:    landmark.Locations[0].LatLng.Longitude,
				Y:    landmark.Locations[0].LatLng.Latitude,
			},
		})
	}
	return rsp, nil
}