package model

import (
	"github.com/rsmarincu/billr/domain"
)

type AverageRateRecord struct {
	Id       string  `db:"id"`
	ClientId string  `db:"clientId"`
	Min      float64 `db:"min"`
	Max      float64 `db:"max"`
	Median   float64 `db:"median"`
	Average  float64 `db:"average"`
	Role     string  `db:"role"`
}

func (a AverageRateRecord) ToDomain() domain.AverageRate {
	return domain.AverageRate{
		Id:       a.Id,
		ClientId: a.ClientId,
		Min:      a.Min,
		Max:      a.Max,
		Median:   a.Median,
		Average:  a.Average,
		Role:     domain.Role(a.Role),
	}
}

func AverageRateToRecord(a domain.AverageRate) AverageRateRecord {
	return AverageRateRecord{
		Id:       a.Id,
		ClientId: a.ClientId,
		Min:      a.Min,
		Max:      a.Max,
		Median:   a.Median,
		Average:  a.Average,
		Role:     string(a.Role),
	}
}
