package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	logicentities "svalka-service/internal/logic/entities"
	logicInterfaces "svalka-service/internal/logic/interfaces"
	pgClient "svalka-service/internal/pgclient"
	repentities "svalka-service/internal/repository/entities"
	"svalka-service/pkg/pg"
	"time"

	"github.com/Masterminds/squirrel"
	sq "github.com/Masterminds/squirrel"
)

type repository struct {
	client pg.Client
}

func NewRepository(ctx context.Context) (logicInterfaces.Repository, error) {
	postgresClient, err := pgClient.GetPgClient(ctx)
	if err != nil {
		return nil, err
	}
	return repository{
		client: postgresClient,
	}, nil
}

// CreateUser ...
func (rep repository) CreateUser(ctx context.Context, user logicentities.User) error {
	userDB := repentities.UserCreate{
		Name: user.Name,
		Position: sql.NullString{
			String: user.Position,
			Valid:  true,
		},
		Phone: sql.NullString{
			String: user.Phone,
			Valid:  true,
		},
		Email:    user.Email,
		Password: user.Password,
	}
	colums, values := pg.ParseDbModel(userDB)
	builder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns(colums...).
		Values(values...).
		Suffix("RETURNING id")

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "user.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

// Login ...
func (rep repository) Login(ctx context.Context, login logicentities.Login) (bool, error) {
	user, err := rep.GetUserByEmail(ctx, login.Email)
	if err != nil {
		return false, err
	}
	if user.Password != login.Password {
		return false, errors.New("Login error: Invalid credentials.")
	}
	return true, nil
}

// CreateRegion ...
func (rep repository) CreateRegion(ctx context.Context, region logicentities.Region) error {
	regionDB := repentities.RegionCreate{
		Name: region.Title,
	}
	colums, values := pg.ParseDbModel(regionDB)
	builder := sq.Insert("regions").
		PlaceholderFormat(sq.Dollar).
		Columns(colums...).
		Values(values...).
		Suffix("RETURNING id")

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "region.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return err
	}

	return nil
}

// GetAllRegions ...
func (rep repository) GetAllRegions(ctx context.Context) ([]logicentities.Region, error) {
	colums, _ := pg.ParseDbModel(repentities.Region{})
	builder := sq.Select(colums...).
		From("regions").
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "regions.GetAll",
		QueryRaw: query,
	}

	var resultDB []*repentities.Region
	err = rep.client.PG().SelectContext(ctx, &resultDB, q, v...)
	if err != nil {
		return nil, err
	}
	var result []logicentities.Region
	for _, r := range resultDB {
		result = append(result, logicentities.Region{
			ID:    r.ID,
			Title: r.Name,
		})
	}
	return result, nil
}

// CreateLandfill ...
func (rep repository) CreateLandfill(ctx context.Context, landfill logicentities.Landfill) error {
	s := pg.StructToDBString(landfill.Manager)
	landfillDB := repentities.LandfillCreate{
		City: sql.NullString{
			String: landfill.City,
			Valid:  true,
		},
		RegionID: landfill.RegionID,
		Illegal:  landfill.Illegal,
		Address: sql.NullString{
			String: landfill.Address,
			Valid:  true,
		},
		Coordinates: landfill.Coordinates,
		PreviewImagePath: sql.NullString{
			String: landfill.PreviewImagePath,
			Valid:  true,
		},
		Manager: sql.NullString{
			String: s,
			Valid:  true,
		},
		Type: sql.NullString{
			String: landfill.Type,
			Valid:  true,
		},
	}
	colums, values := pg.ParseDbModel(landfillDB)
	builder := sq.Insert("landfills").
		PlaceholderFormat(sq.Dollar).
		Columns(colums...).
		Values(values...)

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "landfill.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		rows.Close()
		return err
	}

	rows.Next()
	err = rows.Err()
	if err != nil {
		return err
	}
	return nil
}

// GetLandfills ...
func (rep repository) GetLandfills(ctx context.Context, regionID int64) ([]logicentities.LandfillAllResponse, error) {
	colums, _ := pg.ParseDbModel(repentities.Landfill{})
	builder := sq.Select(colums...).
		From("landfills").
		Where("region_id = ?", regionID).
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "landfills.GetByRegion",
		QueryRaw: query,
	}

	var resultDB []*repentities.Landfill
	err = rep.client.PG().SelectContext(ctx, &resultDB, q, v...)
	if err != nil {
		return nil, err
	}
	var result []logicentities.LandfillAllResponse
	for _, r := range resultDB {
		result = append(result, logicentities.LandfillAllResponse{
			LandfillInfo: logicentities.LandfillInfo{
				ID:               r.ID,
				RegionID:         r.RegionID,
				Illegal:          r.Illegal,
				City:             r.City.String,
				Address:          r.Address.String,
				Type:             r.Type.String,
				Coordinates:      r.Coordinates,
				PreviewImagePath: r.PreviewImagePath.String,
			},
			ViolationsCount: 0, //TODO
		})
	}
	return result, nil
}

// CreateViolation ...
func (rep repository) CreateViolation(ctx context.Context, violation logicentities.Violation) error {
	violationDB := repentities.ViolationCreate{
		Description:   violation.Title,
		DefaultStatus: violation.Default,
	}
	colums, values := pg.ParseDbModel(violationDB)
	builder := sq.Insert("violations").
		PlaceholderFormat(sq.Dollar).
		Columns(colums...).
		Values(values...)

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "violation.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	err = rows.Err()
	if err != nil {
		return err
	}

	return nil

}

// GetAllViolations ...
func (rep repository) GetAllViolations(ctx context.Context) ([]logicentities.Violation, error) {
	colums, _ := pg.ParseDbModel(repentities.Violation{})
	builder := sq.Select(colums...).
		From("violations").
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "violations.GetAll",
		QueryRaw: query,
	}

	var resultDB []*repentities.Violation
	err = rep.client.PG().SelectContext(ctx, &resultDB, q, v...)
	if err != nil {
		return nil, err
	}
	var result []logicentities.Violation
	for _, r := range resultDB {
		result = append(result, logicentities.Violation{
			ID:      r.ID,
			Title:   r.Description,
			Default: r.DefaultStatus,
		})
	}
	return result, nil
}

// CreateSurvey ...
func (rep repository) CreateSurvey(ctx context.Context, surveyRequest logicentities.SurveyRequest) (surveyID int64, err error) {
	surveyDB := repentities.SurveyCreate{
		Date:       surveyRequest.Date,
		LandfillID: surveyRequest.LandfillID,
	}
	colums, values := pg.ParseDbModel(surveyDB)
	builder := sq.Insert("surveys").
		PlaceholderFormat(sq.Dollar).
		Columns(colums...).
		Values(values...).
		Suffix("RETURNING id")

	query, v, err := builder.ToSql()
	if err != nil {
		return 0, err
	}

	q := pg.Query{
		Name:     "survey.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	rows.Next()
	var id int64
	err = rows.Scan(&id)
	if err != nil {
		return 0, nil
	}

	return id, nil
}

// GetSurvey ...
func (rep repository) GetSurvey(ctx context.Context, surveyID int64) (*logicentities.Survey, error) {
	// parse
	columns, _ := pg.ParseDbModel(repentities.Survey{})
	builder := squirrel.
		Select(columns...).
		From("surveys").
		Where("id = ?", surveyID).
		PlaceholderFormat(squirrel.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "surveys.Get",
		QueryRaw: query,
	}
	respDB := repentities.Survey{}
	err = rep.client.PG().GetContext(ctx, &respDB, q, v...)
	if err != nil {
		return nil, err
	}
	resp := logicentities.Survey{
		ID:   respDB.ID,
		Date: respDB.Date,
	}

	// get images
	resp.RawImages, resp.AiImages, err = rep.GetImages(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	// get audits
	resp.Audits, err = rep.GetAudits(ctx, surveyID)
	if err != nil {
		return nil, err
	}

	return &resp, nil
}

// GetUser ...
func (rep repository) GetUser(ctx context.Context, userID int64) (*logicentities.UserInfo, error) {
	// parse
	columns, _ := pg.ParseDbModel(repentities.User{})
	builder := squirrel.
		Select(columns...).
		From("users").
		Where("id = ?", userID).
		PlaceholderFormat(squirrel.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "users.Get",
		QueryRaw: query,
	}
	respDB := repentities.User{}
	err = rep.client.PG().GetContext(ctx, &respDB, q, v...)
	if err != nil {
		return nil, err
	}
	resp := logicentities.UserInfo{
		ID:       respDB.ID,
		Name:     respDB.Name,
		Position: respDB.Position.String,
		Phone:    respDB.Phone.String,
		Email:    respDB.Email,
	}

	return &resp, nil
}

// GetAudits ...
func (rep repository) GetAudits(ctx context.Context, surveyID int64) (audits []logicentities.Audit, err error) {
	// parse
	columns, _ := pg.ParseDbModel(repentities.Audit{})
	builder := squirrel.
		Select(columns...).
		From("audits").
		Where("survey_id = ?", surveyID).
		PlaceholderFormat(squirrel.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "surveys.Get",
		QueryRaw: query,
	}

	var resultDB []*repentities.Audit
	err = rep.client.PG().SelectContext(ctx, &resultDB, q, v...)
	if err != nil {
		return nil, err
	}

	for _, r := range resultDB {
		user, err := rep.GetUser(ctx, r.UserID)
		if err != nil {
			return nil, err
		}
		if user == nil {
			return nil, errors.New("user = nil")
		}

		vs, count, err := rep.GetViolations(ctx, r.ID)
		if err != nil {
			return nil, err
		}
		audits = append(audits, logicentities.Audit{
			ID:                  r.ID,
			Date:                r.Date,
			AiGeneratedStatus:   r.AiGeneratedStatus,
			Auditor:             *user,
			Violations:          vs,
			TrueViolationsCount: count,
		})
	}

	return audits, nil
}

// GetViolations ...
func (rep repository) GetViolations(ctx context.Context, auditID int64) (violations []logicentities.Violation, trueViolationCount int, err error) {
	// getAllviolations
	violations, err = rep.GetAllViolations(ctx)
	if err != nil {
		return nil, 0, err
	}

	// parse
	columns, _ := pg.ParseDbModel(repentities.AuditToViolation{})
	builder := squirrel.
		Select(columns...).
		From("audits_to_violations").
		Where("audit_id = ?", auditID).
		PlaceholderFormat(squirrel.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, 0, err
	}

	q := pg.Query{
		Name:     "surveys.Get",
		QueryRaw: query,
	}

	var resultDB []*repentities.AuditToViolation
	err = rep.client.PG().SelectContext(ctx, &resultDB, q, v...)
	if err != nil {
		return nil, 0, err
	}
	trueViolations := make(map[int64]interface{}, len(resultDB))
	for _, r := range resultDB {
		trueViolations[r.ViolationID] = 0
	}

	for i, v := range violations {
		_, ok := trueViolations[v.ID]
		if ok {
			v := violations[i]
			v.Status = true
			violations[i] = v
		}
	}

	return violations, len(trueViolations), nil
}

// GetImages ...
func (rep repository) GetImages(ctx context.Context, surveyID int64) (rootImages []string, aiImages []string, err error) {
	// parse
	columns, _ := pg.ParseDbModel(repentities.Image{})
	builder := squirrel.
		Select(columns...).
		From("images").
		Where("survey_id = ?", surveyID).
		PlaceholderFormat(squirrel.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, nil, err
	}

	q := pg.Query{
		Name:     "surveys.Get",
		QueryRaw: query,
	}

	var resultDB []*repentities.Image
	err = rep.client.PG().SelectContext(ctx, &resultDB, q, v...)
	if err != nil {
		return nil, nil, err
	}

	for _, r := range resultDB {
		if r.RawStatus {
			rootImages = append(rootImages, r.Path)
			continue
		}
		aiImages = append(aiImages, r.Path)
	}
	return rootImages, aiImages, nil
}

// SaveSurveyResult ...
func (rep repository) SaveSurveyResult(ctx context.Context, surveyResult logicentities.SurveyResult) error {
	// parse
	columns, _ := pg.ParseDbModel(repentities.ImageCreate{})
	builder := squirrel.
		Insert("images").
		Columns(columns...).
		PlaceholderFormat(squirrel.Dollar)

	var images []repentities.ImageCreate
	for _, path := range surveyResult.CdnResult.RawImagesPaths {
		images = append(images, repentities.ImageCreate{
			SurveyID:  surveyResult.SurveyID,
			Path:      path,
			RawStatus: true,
		})
	}
	for _, path := range surveyResult.AiResult.AiImagesPaths {
		images = append(images, repentities.ImageCreate{
			SurveyID:  surveyResult.SurveyID,
			Path:      path,
			RawStatus: false,
		})
	}
	for _, i := range images {
		_, values := pg.ParseDbModel(i)
		builder = builder.Values(values...)
	}

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "image.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return err
	}
	rows.Next()
	if rows.Err() != nil {
		return rows.Err()
	}
	defer rows.Close()

	// CreateAuditToViolations
	err = rep.CreateAudit(ctx, logicentities.AuditRequest{
		Date:              time.Now(),
		SurveyID:          surveyResult.SurveyID,
		AiGeneratedStatus: true,
		Violations:        surveyResult.AiResult.Violations,
	})
	if err != nil {
		return err
	}

	return nil
}

// CreateAudit ...
func (rep repository) CreateAudit(ctx context.Context, auditRequest logicentities.AuditRequest) error {
	var userID int64
	if auditRequest.Email != "" {
		// GetUser ...
		user, err := rep.GetUserByEmail(ctx, auditRequest.Email)
		if err != nil {
			return err
		}
		userID = user.ID
	}

	// Audit ...
	surveyDB := repentities.AuditCreate{
		Date:              auditRequest.Date,
		UserID:            userID,
		SurveyID:          auditRequest.SurveyID,
		AiGeneratedStatus: auditRequest.AiGeneratedStatus,
	}
	colums, values := pg.ParseDbModel(surveyDB)
	builder := sq.Insert("audits").
		PlaceholderFormat(sq.Dollar).
		Columns(colums...).
		Values(values...).
		Suffix("RETURNING id")

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "audit.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return err
	}
	defer rows.Close()

	rows.Next()
	var auditID int64
	err = rows.Scan(&auditID)
	if err != nil {
		return err
	}
	// CreateAuditToViolations
	var auditsToViolations []repentities.AuditToViolation
	for _, v := range auditRequest.Violations {
		auditsToViolations = append(auditsToViolations, repentities.AuditToViolation{
			ViolationID: v.ID,
			AuditID:     auditID,
		})
	}
	err = rep.CreateAuditToViolations(ctx, auditsToViolations)
	if err != nil {
		return err
	}

	return nil
}

// GetUserIDByEmail ...
func (rep repository) GetUserByEmail(ctx context.Context, userEmail string) (*logicentities.User, error) {
	// parse
	columns, _ := pg.ParseDbModel(repentities.User{})
	builder := squirrel.
		Select(columns...).
		From("users").
		Where("email = ?", userEmail).
		PlaceholderFormat(squirrel.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "users.GetByEmail",
		QueryRaw: query,
	}
	respDB := repentities.User{}
	err = rep.client.PG().GetContext(ctx, &respDB, q, v...)
	if err != nil {
		return nil, err
	}
	resp := logicentities.User{
		UserInfo: logicentities.UserInfo{
			ID:       respDB.ID,
			Name:     respDB.Name,
			Position: respDB.Position.String,
			Phone:    respDB.Phone.String,
			Email:    respDB.Email,
		},
		Password: respDB.Password,
	}

	return &resp, nil
}

// CreateAuditToViolations ...
func (rep repository) CreateAuditToViolations(ctx context.Context, auditsToViolations []repentities.AuditToViolation) error {
	colums, _ := pg.ParseDbModel(repentities.AuditToViolation{})
	builder := sq.Insert("audits_to_violations").
		PlaceholderFormat(sq.Dollar).
		Columns(colums...)

	for _, auditToViolation := range auditsToViolations {
		_, volumes := pg.ParseDbModel(auditToViolation)
		builder = builder.Values(volumes...)
	}

	query, v, err := builder.ToSql()
	if err != nil {
		return err
	}

	q := pg.Query{
		Name:     "audits_to_violations.Create",
		QueryRaw: query,
	}

	rows, err := rep.client.PG().QueryContext(ctx, q, v...)
	if err != nil {
		return err
	}
	defer rows.Close()

	return nil
}

// GetLandfill ...
func (rep repository) GetLandfill(ctx context.Context, landfillID int64) (*logicentities.LandfillExtended, error) {
	colums, _ := pg.ParseDbModel(repentities.Landfill{})
	builder := sq.Select(colums...).
		From("landfills").
		Where("id = ?", landfillID).
		PlaceholderFormat(sq.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "landfills.GetByID",
		QueryRaw: query,
	}

	var resultDB repentities.Landfill
	err = rep.client.PG().GetContext(ctx, &resultDB, q, v...)
	if err != nil {
		return nil, err
	}
	// get survey
	surveys, err := rep.GetSurveysByLandfilId(ctx, landfillID)
	if err != nil {
		return nil, err
	}
	manager := logicentities.Manager{}
	err = json.Unmarshal([]byte(resultDB.Manager.String), &manager)
	if err != nil {
		return nil, err
	}
	result := logicentities.LandfillExtended{
		LandfillInfo: logicentities.LandfillInfo{
			ID:               resultDB.ID,
			RegionID:         resultDB.RegionID,
			Illegal:          resultDB.Illegal,
			City:             resultDB.City.String,
			Address:          resultDB.Address.String,
			Type:             resultDB.Type.String,
			Coordinates:      resultDB.Coordinates,
			PreviewImagePath: resultDB.PreviewImagePath.String,
		},
		Manager:         manager,
		Surveys:         surveys,
		ViolationsCount: 0, //TODO
	}

	return &result, nil
}

// GetSurveysByLandfilId ...
func (rep repository) GetSurveysByLandfilId(ctx context.Context, landfillID int64) ([]logicentities.Survey, error) {
	// parse
	columns, _ := pg.ParseDbModel(repentities.Survey{})
	builder := squirrel.
		Select(columns...).
		From("surveys").
		Where("landfill_id = ?", landfillID).
		PlaceholderFormat(squirrel.Dollar)

	query, v, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	q := pg.Query{
		Name:     "surveys.Get",
		QueryRaw: query,
	}
	respDB := []*repentities.Survey{}
	err = rep.client.PG().SelectContext(ctx, &respDB, q, v...)
	if err != nil {
		return nil, err
	}

	var resp []logicentities.Survey
	for _, survey := range respDB {
		// get images
		rawImages, aiImages, err := rep.GetImages(ctx, int64(survey.ID))
		if err != nil {
			return nil, err
		}

		// get audits
		audits, err := rep.GetAudits(ctx, int64(survey.ID))
		if err != nil {
			return nil, err
		}

		resp = append(resp, logicentities.Survey{
			ID:        survey.ID,
			Date:      survey.Date,
			RawImages: rawImages,
			AiImages:  aiImages,
			Audits:    audits,
		})
	}

	return resp, nil
}
