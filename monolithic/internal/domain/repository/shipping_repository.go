package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// ShippingMethodRepository defines the interface for shipping method data access
type ShippingMethodRepository interface {
	Create(method *entity.ShippingMethod) error
	GetByID(methodID uint) (*entity.ShippingMethod, error)
	List(active bool) ([]*entity.ShippingMethod, error)
	Update(method *entity.ShippingMethod) error
	Delete(methodID uint) error
}

// ShippingZoneRepository defines the interface for shipping zone data access
type ShippingZoneRepository interface {
	Create(zone *entity.ShippingZone) error
	GetByID(zoneID uint) (*entity.ShippingZone, error)
	List(active bool) ([]*entity.ShippingZone, error)
	Update(zone *entity.ShippingZone) error
	Delete(zoneID uint) error
}

// ShippingRateRepository defines the interface for shipping rate data access
type ShippingRateRepository interface {
	Create(rate *entity.ShippingRate) error
	GetByID(rateID uint) (*entity.ShippingRate, error)
	GetByMethodID(methodID uint) ([]*entity.ShippingRate, error)
	GetByZoneID(zoneID uint) ([]*entity.ShippingRate, error)
	GetAvailableRatesForAddress(address entity.Address, orderValue int64) ([]*entity.ShippingRate, error)
	CreateWeightBasedRate(weightRate *entity.WeightBasedRate) error
	CreateValueBasedRate(valueRate *entity.ValueBasedRate) error
	GetWeightBasedRates(rateID uint) ([]entity.WeightBasedRate, error)
	GetValueBasedRates(rateID uint) ([]entity.ValueBasedRate, error)
	Update(rate *entity.ShippingRate) error
	Delete(rateID uint) error
}
