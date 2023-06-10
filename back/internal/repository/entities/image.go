package repentities

type Image struct {
	ID uint64 `db:"id"`
	ImageCreate
}

type ImageCreate struct {
	SurveyID  uint64 `db:"survey_id"`
	Path      string `db:"path"`
	RawStatus bool   `db:"raw_status"`
}

type ImageConverter struct {
}

// ToDB ...
func (ImageConverter) ToDB(imagesPaths []string, surveyId uint64, rawStatus bool) (images []ImageCreate) {
	for _, path := range imagesPaths {
		images = append(images, ImageCreate{
			SurveyID:  surveyId,
			Path:      path,
			RawStatus: rawStatus,
		})
	}
	return images
}
