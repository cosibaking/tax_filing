package complianceassistantin

import "xygo/internal/logic/compliance/profile"

type ProfileSaveInp = profile.SaveInput
type ProfileData = profile.Data

type ProfileSaveModel struct {
	Profile *profile.Profile `json:"profile"`
	Changed bool             `json:"changed"`
}
