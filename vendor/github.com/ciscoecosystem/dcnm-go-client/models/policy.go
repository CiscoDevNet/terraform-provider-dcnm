package models

type Policy struct {
	Id                  string      `json:"id,omitempty"`
	PolicyId            string      `json:"policyId,omitempty"`
	Source              string      `json:"source,omitempty"`
	Description         string      `json:"description,omitempty"`
	SerialNumber        string      `json:"serialNumber,omitempty"`
	EntityType          string      `json:"entityType,omitempty"`
	EntityName          string      `json:"entityName,omitempty"`
	Priority            string      `json:"priority,omitempty"`
	TemplateName        string      `json:"templateName,omitempty"`
	TemplateContentType string      `json:"templateContentType,omitempty"`
	NVPairs             interface{} `json:"nvpairs,omitempty"`
}

func (policy *Policy) ToMap() (map[string]interface{}, error) {
	policyMap := make(map[string]interface{})
	A(policyMap, "id", policy.Id)
	A(policyMap, "source", policy.Source)

	A(policyMap, "policyId", policy.PolicyId)

	A(policyMap, "serialNumber", policy.SerialNumber)

	A(policyMap, "entityType", policy.EntityType)

	A(policyMap, "entityName", policy.EntityName)

	A(policyMap, "templateName", policy.TemplateName)

	A(policyMap, "priority", policy.Priority)

	A(policyMap, "description", policy.Description)

	A(policyMap, "templateContentType", policy.TemplateContentType)

	if policy.NVPairs != nil {
		A(policyMap, "nvPairs", policy.NVPairs)
	}

	return policyMap, nil
}
