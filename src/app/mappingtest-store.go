package app

import (
	"encoding/json"
	"fmt"
	"log/slog"

	_ "github.com/go-sql-driver/mysql"
)

func (s *Store) GetMappingtests1() ([]Mappingtest1, error) {
	slog.Debug("store::GetMappingtests1() - Executing")
	var tt []Mappingtest1
	sql := `SELECT id, anint, abigint, anintunsighed, abigintunsigned, adecimal, afloat, adouble, aboolean, avarchar, adatetime, ablob
		FROM mappingtest ORDER BY id`
	if err := s.db.Select(&tt, sql); err != nil {
		slog.Error("store::GetMappingtests1() - Failed to find mappingtests", "error", err)
		return []Mappingtest1{}, fmt.Errorf("Failed to find mappingtests: %w", err)
	}
	for i := range tt {
		t := tt[i]
		j, _ := json.MarshalIndent(t, "", "  ")
		fmt.Println("GetMappingtests1() -", "index", i, "value \n", string(j))
	}
	return tt, nil
}

func (s *Store) GetMappingtests2() ([]Mappingtest2, error) {
	slog.Debug("store::GetMappingtests2() - Executing")
	var tt []Mappingtest2
	sql := `SELECT id, anint, abigint, anintunsighed, abigintunsigned, adecimal, afloat, adouble, aboolean, avarchar, adatetime, ablob
		FROM mappingtest ORDER BY id`
	if err := s.db.Select(&tt, sql); err != nil {
		slog.Error("store::GetMappingtests2() - Failed to find mappingtests", "error", err)
		return []Mappingtest2{}, fmt.Errorf("Failed to find mappingtests: %w", err)
	}
	for i := range tt {
		t := tt[i]
		j, _ := json.MarshalIndent(t, "", "  ")
		fmt.Println("GetMappingtests2() -", "index", i, "value \n", string(j))
	}
	return tt, nil
}

/*
 */
