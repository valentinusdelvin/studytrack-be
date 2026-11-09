package enums

import (
	"database/sql/driver"
	"fmt"
)

type TugasStatus string

const (
	TugasStatusSelesai      TugasStatus = "selesai"
	TugasStatusBelumSelesai TugasStatus = "belum"
)

func (t *TugasStatus) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		*t = TugasStatus(v)
	case string:
		*t = TugasStatus(v)
	case nil:
		return nil
	default:
		return fmt.Errorf("unsupported TugasStatus type: %T", value)
	}
	return nil
}

func (t TugasStatus) Value() (driver.Value, error) {
	return string(t), nil
}

// Every enum type should implement EnumType for migration generation
// Look at database/migrate.go for how this is used in migration generation
// After implementing EnumType, add the enum struct to the enums list in loadEnums function in database/migrate.go
type TugasStatusEnum struct{}

func (TugasStatusEnum) TypeName() string { return "tugas_status" }
func (TugasStatusEnum) Values() []string {
	return []string{
		string(TugasStatusSelesai),
		string(TugasStatusBelumSelesai),
	}
}
