package report

import (
	"context"
	"fmt"
	"time"

	"xygo/internal/library/reportpdf"
)

func (s *Service) ExportPDF(ctx context.Context, memberID, id uint64) (filename string, content []byte, err error) {
	item, err := s.repository.Get(ctx, memberID, id)
	if err != nil {
		return "", nil, err
	}

	companyName := "未命名企业"
	if s.companyLoader != nil {
		opc, err := s.companyLoader(ctx, memberID)
		if err != nil {
			return "", nil, err
		}
		if opc != nil && opc.CompanyName != "" {
			companyName = opc.CompanyName
		}
	}

	data := reportpdf.Data{
		CompanyName: companyName,
		PeriodKey:   item.PeriodKey,
		Version:     item.Version,
		Status:      item.Status,
		GeneratedAt: formatReportTime(item.CreatedAt),
	}
	if item.StructuredReport == nil {
		data.LegacyContent = item.Content
	} else {
		mapStructuredPDFData(&data, *item.StructuredReport)
	}
	content, err = s.pdfGenerator(data)
	if err != nil {
		return "", nil, err
	}

	safePeriod := item.PeriodKey
	if !validPeriod(safePeriod) {
		safePeriod = fmt.Sprintf("report-%d", item.ID)
	}
	return fmt.Sprintf("monthly-checkup-%s-v%d.pdf", safePeriod, item.Version), content, nil
}

func formatReportTime(timestamp uint64) string {
	if timestamp == 0 {
		return ""
	}
	return time.Unix(int64(timestamp), 0).Format("2006-01-02 15:04:05")
}

func mapStructuredPDFData(data *reportpdf.Data, structured StructuredReport) {
	data.Summary = reportpdf.Summary{
		Conclusion:       structured.Summary.Conclusion,
		CompletenessRate: structured.Summary.CompletenessRate,
		HighCount:        structured.Summary.HighCount,
		MediumCount:      structured.Summary.MediumCount,
		LowCount:         structured.Summary.LowCount,
		DataNotice:       structured.Summary.DataNotice,
	}
	data.Categories = make([]reportpdf.Category, 0, len(structured.Categories))
	for _, category := range structured.Categories {
		mapped := reportpdf.Category{
			Code: category.Code, Name: category.Name, Status: category.Status, Summary: category.Summary,
			Checks: make([]reportpdf.Check, 0, len(category.Checks)),
		}
		for _, check := range category.Checks {
			mapped.Checks = append(mapped.Checks, reportpdf.Check{Code: check.Code, Name: check.Name, Status: check.Status, Message: check.Message})
		}
		data.Categories = append(data.Categories, mapped)
	}
	data.Anomalies = make([]reportpdf.Anomaly, 0, len(structured.Anomalies))
	for _, anomaly := range structured.Anomalies {
		data.Anomalies = append(data.Anomalies, reportpdf.Anomaly{
			Code: anomaly.Code, CategoryCode: anomaly.CategoryCode, Title: anomaly.Title, Severity: anomaly.Severity,
			Facts: anomaly.Facts, Basis: anomaly.Basis, Impact: anomaly.Impact, Recommendation: anomaly.Recommendation,
			RequiredMaterials: append([]string(nil), anomaly.RequiredMaterials...), DueDate: anomaly.DueDate,
			RequiresManualReview: anomaly.RequiresManualReview, RuleVersion: anomaly.RuleVersion,
		})
	}
}
