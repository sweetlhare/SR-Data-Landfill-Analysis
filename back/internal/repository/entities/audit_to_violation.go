package repentities

type AuditToViolation struct {
	ViolationID uint64 `db:"violation_id"`
	AuditID     uint64 `db:"audit_id"`
}
