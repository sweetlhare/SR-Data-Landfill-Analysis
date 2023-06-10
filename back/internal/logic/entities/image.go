package logicentities

import "mime/multipart"

type File struct {
	*multipart.FileHeader
}

type CdnResult struct {
	RawImagesPaths []string
}

type AiResult struct {
	AiImagesPaths []string
	Violations    []Violation
}

type ImageCreate struct {
	SurveyID  uint64 `json:"survey_id"`
	Path      string `json:"path"`
	RawStatus bool   `json:"raw_status"`
}
