package repentities

type Image struct {
	ID        int64  `db:"id"`
	SurveyID  int64  `db:"survey_id"`
	Path      string `db:"path"`
	RawStatus bool   `db:"raw_status"`
}

type ImageCreate struct {
	SurveyID  int64  `db:"survey_id"`
	Path      string `db:"path"`
	RawStatus bool   `db:"raw_status"`
}
