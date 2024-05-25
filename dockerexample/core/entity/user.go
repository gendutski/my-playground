package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"strings"
	"time"
)

type GenderType int

const (
	InvalidGender GenderType = iota
	Male
	Female
)

const (
	maleGender   string = "male"
	femaleGender string = "female"
)

type User struct {
	ID        int        `gorm:"primaryKey"`
	Name      string     `gorm:"size:255" json:"name"`
	Gender    GenderType `gorm:"type:enum('male','female');default:'male'" json:"gender"`
	CreatedAt time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdtedAt  time.Time  `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
}

// Implement custom unmarshaler for GenderType
func (g *GenderType) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	switch s {
	case maleGender:
		*g = Male
	case femaleGender:
		*g = Female
	default:
		return errors.New("invalid gender value")
	}
	return nil
}

// Implement custom marshaler for GenderType (optional)
func (g GenderType) MarshalJSON() ([]byte, error) {
	var s string
	switch g {
	case Male:
		s = maleGender
	case Female:
		s = femaleGender
	default:
		return nil, errors.New("invalid gender value")
	}
	return json.Marshal(s)
}

// gorm scan value
func (e *GenderType) Scan(value interface{}) error {
	src, ok := value.([]uint8)
	if !ok {
		return errors.New("invalid gender type")
	}

	switch strings.ToLower(string(src)) {
	case maleGender:
		*e = Male
	case femaleGender:
		*e = Female
	}
	return nil
}

// gorm set value
func (e GenderType) Value() (driver.Value, error) {
	var result string
	switch e {
	case Male:
		result = maleGender
	case Female:
		result = femaleGender
	}
	return result, nil
}

// to string
func (e GenderType) String() string {
	switch e {
	case Male:
		return maleGender
	case Female:
		return femaleGender
	}
	return ""
}

type GetRequest struct {
	Offset int `query:"offset" json:"offset"`
	Limit  int `query:"limit" json:"limit"`
}

type GetResponse struct {
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	CreatedAt string `json:"createdAt"`
	UpdtedAt  string `json:"updatedAt"`
}

func ToGetResponse(e []*User) []*GetResponse {
	result := make([]*GetResponse, len(e))
	for i := range e {
		result[i] = &GetResponse{
			Name:      e[i].Name,
			Gender:    e[i].Gender.String(),
			CreatedAt: e[i].CreatedAt.Format(time.RFC3339),
			UpdtedAt:  e[i].UpdtedAt.Format(time.RFC3339),
		}
	}
	return result
}
