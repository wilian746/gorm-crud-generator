package product

import (
	"encoding/json"
	"github.com/wilian746/gorm-crud-generator/internal/entities"
)

type Product struct {
	entities.Base
	Name string
}

func InterfaceToModel(data interface{}) (instance *Product, err error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return instance, err
	}

	return instance, json.Unmarshal(bytes, instance)
}

func (p *Product) TableName() string {
	return "products"
}

func (p *Product) Bytes() ([]byte, error) {
	return json.Marshal(p)
}
