package service

import (
	"context"

	"xygo/internal/logic/compliance/profile"
	"xygo/internal/logic/compliance/ruleengine"
)

type IComplianceProfile interface {
	Save(ctx context.Context, in profile.SaveInput) (*profile.Profile, bool, error)
}

type IComplianceRuleVersion interface {
	Transition(ctx context.Context, in ruleengine.TransitionInput) (*ruleengine.RuleVersion, error)
	Simulate(ctx context.Context, versionID uint64, facts map[string]any) (bool, error)
}

var localComplianceProfile IComplianceProfile
var localComplianceRuleVersion IComplianceRuleVersion

func ComplianceProfile() IComplianceProfile {
	if localComplianceProfile == nil {
		panic("implement not found for interface IComplianceProfile, forgot register?")
	}
	return localComplianceProfile
}

func RegisterComplianceProfile(i IComplianceProfile) {
	localComplianceProfile = i
}

func ComplianceRuleVersion() IComplianceRuleVersion {
	if localComplianceRuleVersion == nil {
		panic("implement not found for interface IComplianceRuleVersion, forgot register?")
	}
	return localComplianceRuleVersion
}

func RegisterComplianceRuleVersion(i IComplianceRuleVersion) {
	localComplianceRuleVersion = i
}
