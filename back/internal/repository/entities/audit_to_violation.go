package repentities

type AuditToViolation struct {
	ViolationID int64 `db:"violation_id"`
	AuditID     int64 `db:"audit_id"`
}
