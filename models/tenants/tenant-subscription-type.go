package models

import "go-multitenancy-boilerplate/models"

type TenantSubscriptionType struct {
	models.Model
	SubscriptionName    string
	SubscriptionPrice   uint
	SubscriptionPeriod  uint // renewal period denoted as 1-24
	SubscriptionRenewal bool
}
