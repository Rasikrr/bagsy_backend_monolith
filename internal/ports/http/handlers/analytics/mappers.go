package analytics

import (
	domainAnalytics "github.com/Rasikrr/bagsy_backend_monolith/internal/domain/analytics"
)

func toKpiValueDTO(v domainAnalytics.KpiValue) kpiValueDTO {
	return kpiValueDTO{Value: v.Value, Prev: v.Prev, DeltaPercent: v.DeltaPercent}
}

func toKpiBlockDTO(b domainAnalytics.KPIBlock) kpiBlockDTO {
	return kpiBlockDTO{
		Revenue:             toKpiValueDTO(b.Revenue),
		Bookings:            toKpiValueDTO(b.Bookings),
		Clients:             toKpiValueDTO(b.Clients),
		AvgCheck:            toKpiValueDTO(b.AvgCheck),
		LoadPercent:         toKpiValueDTO(b.LoadPercent),
		CancellationPercent: toKpiValueDTO(b.CancellationPercent),
	}
}

func toDailyPointDTOs(points []domainAnalytics.DailyPoint) []dailyPointDTO {
	out := make([]dailyPointDTO, 0, len(points))
	for _, p := range points {
		out = append(out, dailyPointDTO{
			Date:      p.Date.Format(dateLayout),
			Value:     p.Value,
			PrevValue: p.PrevValue,
		})
	}
	return out
}

func toTopItemDTOs(items []domainAnalytics.TopItem) []topItemDTO {
	out := make([]topItemDTO, 0, len(items))
	for _, it := range items {
		out = append(out, topItemDTO{
			ID:      it.ID.String(),
			Name:    it.Name,
			Revenue: it.Revenue,
			Share:   it.Share,
		})
	}
	return out
}

func toFunnelDTOs(stages []domainAnalytics.FunnelStage) []funnelStageDTO {
	out := make([]funnelStageDTO, 0, len(stages))
	for _, s := range stages {
		out = append(out, funnelStageDTO{Key: s.Key, Count: s.Count, Conversion: s.Conversion})
	}
	return out
}

func toHeatmapDTOs(cells []domainAnalytics.HeatmapCell) []heatmapCellDTO {
	out := make([]heatmapCellDTO, 0, len(cells))
	for _, c := range cells {
		out = append(out, heatmapCellDTO{Weekday: c.Weekday, Hour: c.Hour, Value: c.Value})
	}
	return out
}

func toInsightDTOs(insights []domainAnalytics.Insight) []insightDTO {
	out := make([]insightDTO, 0, len(insights))
	for _, i := range insights {
		out = append(out, insightDTO{Key: i.Key, Level: i.Level, Params: i.Params})
	}
	return out
}

func toClientsBreakdownDTO(b domainAnalytics.ClientsBreakdown) clientsBreakdownDTO {
	return clientsBreakdownDTO{New: b.New, Returning: b.Returning}
}

func toOverviewResponse(r *domainAnalytics.OverviewReport) overviewResponse {
	return overviewResponse{
		KPI:          toKpiBlockDTO(r.KPI),
		RevenueByDay: toDailyPointDTOs(r.RevenueByDay),
		TopEmployees: toTopItemDTOs(r.TopEmployees),
		TopServices:  toTopItemDTOs(r.TopServices),
		Funnel:       toFunnelDTOs(r.Funnel),
		Heatmap:      toHeatmapDTOs(r.Heatmap),
		Insights:     toInsightDTOs(r.Insights),
	}
}

func toMeResponse(r *domainAnalytics.MeReport) meResponse {
	return meResponse{
		KPI:              toKpiBlockDTO(r.KPI),
		RevenueByDay:     toDailyPointDTOs(r.RevenueByDay),
		TopServices:      toTopItemDTOs(r.TopServices),
		Heatmap:          toHeatmapDTOs(r.Heatmap),
		ClientsBreakdown: toClientsBreakdownDTO(r.ClientsBreakdown),
	}
}

func toStaffResponse(r *domainAnalytics.StaffReport) staffResponse {
	rows := make([]staffRowDTO, 0, len(r.Rows))
	for _, row := range r.Rows {
		rows = append(rows, staffRowDTO{
			EmployeeID:  row.EmployeeID.String(),
			FullName:    row.FullName,
			Revenue:     row.Revenue,
			Bookings:    row.Bookings,
			AvgCheck:    row.AvgCheck,
			LoadPercent: row.LoadPercent,
			Cancellations: cancellationsDTO{
				Count:   row.CancellationsCount,
				Percent: row.CancellationsPercent,
			},
			Rating: row.Rating,
		})
	}

	weekday := make([]weekdayLoadDTO, 0, len(r.WeekdayLoad))
	for _, w := range r.WeekdayLoad {
		weekday = append(weekday, weekdayLoadDTO{
			EmployeeID: w.EmployeeID.String(),
			Weekday:    w.Weekday,
			Value:      w.Value,
		})
	}

	return staffResponse{
		Rows:        rows,
		WeekdayLoad: weekday,
		Insights:    toInsightDTOs(r.Insights),
	}
}

func toStaffDetailResponse(r *domainAnalytics.StaffDetailReport) staffDetailResponse {
	hourly := make([]hourlyLoadDTO, 0, len(r.HourlyLoad))
	for _, h := range r.HourlyLoad {
		hourly = append(hourly, hourlyLoadDTO{Hour: h.Hour, Value: h.Value})
	}

	return staffDetailResponse{
		Employee: employeeRefDTO{
			ID:       r.Employee.ID.String(),
			FullName: r.Employee.FullName,
			// avatar_url пока не отдаётся (нужна генерация presigned URL из media UC).
		},
		KPI:              toKpiBlockDTO(r.KPI),
		RevenueByDay:     toDailyPointDTOs(r.RevenueByDay),
		TopServices:      toTopItemDTOs(r.TopServices),
		HourlyLoad:       hourly,
		ClientsBreakdown: toClientsBreakdownDTO(r.ClientsBreakdown),
	}
}

func toFinanceResponse(r *domainAnalytics.FinanceReport) financeResponse {
	payroll := make([]payrollRowDTO, 0, len(r.Payroll))
	for _, p := range r.Payroll {
		payroll = append(payroll, payrollRowDTO{
			EmployeeID:        p.EmployeeID.String(),
			FullName:          p.FullName,
			CommissionPercent: p.CommissionPercent,
			Amount:            p.Amount,
		})
	}
	return financeResponse{
		Revenue: financeRevenueDTO{
			Services: r.RevenueServices,
			Products: r.RevenueProducts,
			Total:    r.RevenueTotal,
		},
		Payroll:       payroll,
		PayrollTotal:  r.PayrollTotal,
		GrossProfit:   r.GrossProfit,
		MarginPercent: r.MarginPercent,
	}
}

func toClientsResponse(r *domainAnalytics.ClientsReport) clientsResponse {
	segments := make([]segmentDTO, 0, len(r.Segments))
	for _, s := range r.Segments {
		segments = append(segments, segmentDTO{Key: string(s.Key), Count: s.Count, Share: s.Share})
	}
	cohorts := make([]cohortDTO, 0, len(r.Cohorts))
	for _, c := range r.Cohorts {
		cohorts = append(cohorts, cohortDTO{Month: c.Month, NewCount: c.NewCount, ActivePercent: c.ActivePercent})
	}
	return clientsResponse{
		KPI: clientKpiDTO{
			Total:     toKpiValueDTO(r.KPI.Total),
			New:       toKpiValueDTO(r.KPI.New),
			Returning: toKpiValueDTO(r.KPI.Returning),
			Lost:      toKpiValueDTO(r.KPI.Lost),
		},
		Segments: segments,
		Retention: retentionDTO{
			After1: r.Retention.After1,
			After2: r.Retention.After2,
			After3: r.Retention.After3,
		},
		Cohorts: cohorts,
	}
}
