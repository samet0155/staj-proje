package rules

type Rule struct {
	Pattern     string
	Description string
}

var Rules = []Rule{
	{Pattern: "password", Description: "Hardcoded password/secret detected"},
	{Pattern: "secret", Description: "Hardcoded password/secret detected"},
	{Pattern: "permitrootlogin yes", Description: "PermitRootLogin is enabled"},
	{Pattern: "allowrootlogin yes", Description: "AllowRootLogin is enabled"},
}
