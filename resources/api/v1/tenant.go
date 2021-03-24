package v1resources

type CreateNewTenantRequest struct {
	SubDomainIdentifier string `form:"subDomainIdentifier" json:"subDomainIdentifier" binding:"required"`
}
