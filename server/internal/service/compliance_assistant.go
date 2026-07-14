package service

import (
	"context"

	"xygo/internal/logic/compliance/profile"
)

type IComplianceProfile interface {
	Save(ctx context.Context, in profile.SaveInput) (*profile.Profile, bool, error)
}

var localComplianceProfile IComplianceProfile

func ComplianceProfile() IComplianceProfile {
	if localComplianceProfile == nil {
		panic("implement not found for interface IComplianceProfile, forgot register?")
	}
	return localComplianceProfile
}

func RegisterComplianceProfile(i IComplianceProfile) {
	localComplianceProfile = i
}
