module github.com/CiscoDevNet/terraform-provider-dcnm

go 1.15

require (
	github.com/ciscoecosystem/dcnm-go-client v0.1.4
	github.com/hashicorp/terraform-plugin-sdk/v2 v2.4.3
)

replace github.com/ciscoecosystem/dcnm-go-client => ../../ciscoecosystem/dcnm-go-client
