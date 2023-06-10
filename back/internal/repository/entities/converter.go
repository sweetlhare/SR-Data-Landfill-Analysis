package repentities

type Converter struct {
	User      UserConverter
	Audit     AuditConverter
	Lanfill   LandfillConverter
	Violation ViolationConverter
	Survey    SurveyConverter
	Region    RegionConverter
	Image     ImageConverter
}
