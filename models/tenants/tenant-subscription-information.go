package models

import "go-multitenancy-boilerplate/models"

type TenantSubscriptionInformation struct {
	models.Model
	TenantId         uint
	SubscriptionType uint // This is linked to the TenantSubscriptionType Table
}
