package entities

import (
	"github.com/google/uuid"
	"time"
)

type Doctor struct {
	Interface
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name string
	CreatedAt time.Time
	UpdateAt time.Time
	DoctorsPatients []DoctorPatient `gorm:"foreignkey:DoctorID;association_foreignkey:ID"`
}

func (d *Doctor) TableName() string {
	return "doctors"
}

type DoctorPatient struct {
	Interface
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
	PatientID uuid.UUID `sql:"type:uuid REFERENCES patients(id) ON DELETE CASCADE"`
	DoctorID uuid.UUID `sql:"type:uuid REFERENCES doctors(id) ON DELETE CASCADE"`
	CreatedAt time.Time
	UpdateAt time.Time
	Patient Patient `gorm:"foreignkey:PatientID;association_foreignkey:ID"`
	Doctor Doctor `gorm:"foreignkey:DoctorID;association_foreignkey:ID"`
}

func (dp *DoctorPatient) TableName() string {
	return "doctors_patients"
}

type Patient struct {
	Interface
	ID uuid.UUID `gorm:"type:uuid;primary_key;"`
	Name string
	CreatedAt time.Time
	UpdateAt time.Time
	DoctorsPatients []DoctorPatient `gorm:"foreignkey:PatientID;association_foreignkey:ID"`
}

func (p *Patient) TableName() string {
	return "patients"
}