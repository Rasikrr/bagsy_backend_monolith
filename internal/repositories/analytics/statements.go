package analytics

// Все запросы аналитики — простые агрегаты (GROUP BY + SUM/COUNT/MIN/MAX/AVG).
// Бизнес-логика (классификация, перцентили, нормализация) — в domain/analytics.
//
// Единый набор плейсхолдеров для запросов по записям:
//   $1 — organization_id
//   $2 — location_ids (uuid[], пустой = все локации организации)
//   $3 — from (дата начала, включительно)
//   $4 — to   (дата конца, включительно; в SQL добавляем +1 день для exclusive верхней границы)
//   $5 — employee_id (uuid, NULL = все сотрудники)

const (
	// scopeFilter — общий фильтр по организации, локациям, периоду и сотруднику.
	scopeFilter = `
		organization_id = $1
		AND (CARDINALITY($2::uuid[]) = 0 OR location_id = ANY($2))
		AND start_at >= $3 AND start_at < $4::date + INTERVAL '1 day'
		AND ($5::uuid IS NULL OR employee_id = $5)
	`

	getKPIRaw = `
		SELECT
			COALESCE(SUM(price) FILTER (WHERE status = 'completed'), 0)::float8           AS revenue,
			COUNT(*) FILTER (WHERE status = 'completed')                                  AS bookings,
			COUNT(DISTINCT customer_id) FILTER (WHERE status = 'completed')               AS clients,
			COALESCE(SUM(duration_minutes) FILTER (WHERE status = 'completed'), 0)::float8 AS duration_minutes,
			COUNT(*)                                                                      AS created,
			COUNT(*) FILTER (WHERE status = 'cancelled')                                  AS cancelled
		FROM appointments
		WHERE ` + scopeFilter

	getRevenueByDay = `
		SELECT
			(start_at AT TIME ZONE 'UTC')::date AS date,
			COALESCE(SUM(price), 0)::float8     AS revenue
		FROM appointments
		WHERE ` + scopeFilter + ` AND status = 'completed'
		GROUP BY 1
		ORDER BY 1
	`

	getEmployeeRevenue = `
		SELECT
			a.employee_id                   AS id,
			e.full_name                     AS name,
			COALESCE(SUM(a.price), 0)::float8 AS revenue
		FROM appointments a
		JOIN employees e ON e.id = a.employee_id
		WHERE a.organization_id = $1
			AND (CARDINALITY($2::uuid[]) = 0 OR a.location_id = ANY($2))
			AND a.start_at >= $3 AND a.start_at < $4::date + INTERVAL '1 day'
			AND ($5::uuid IS NULL OR a.employee_id = $5)
			AND a.status = 'completed'
		GROUP BY a.employee_id, e.full_name
	`

	getServiceRevenue = `
		SELECT
			a.service_id                    AS id,
			s.name                          AS name,
			COALESCE(SUM(a.price), 0)::float8 AS revenue
		FROM appointments a
		JOIN services s ON s.id = a.service_id
		WHERE a.organization_id = $1
			AND (CARDINALITY($2::uuid[]) = 0 OR a.location_id = ANY($2))
			AND a.start_at >= $3 AND a.start_at < $4::date + INTERVAL '1 day'
			AND ($5::uuid IS NULL OR a.employee_id = $5)
			AND a.status = 'completed'
		GROUP BY a.service_id, s.name
	`

	getFunnel = `
		SELECT
			COUNT(*)                                                                            AS created,
			COUNT(*) FILTER (WHERE status IN ('confirmed', 'in_progress', 'completed'))         AS confirmed,
			COUNT(*) FILTER (WHERE status = 'completed')                                        AS completed
		FROM appointments
		WHERE ` + scopeFilter

	getHeatmap = `
		SELECT
			EXTRACT(DOW  FROM start_at AT TIME ZONE 'UTC')::int AS weekday,
			EXTRACT(HOUR FROM start_at AT TIME ZONE 'UTC')::int AS hour,
			COUNT(*)                                            AS cnt
		FROM appointments
		WHERE ` + scopeFilter + ` AND status = 'completed'
		GROUP BY 1, 2
	`

	getScheduleMinutes = `
		SELECT COALESCE(SUM(EXTRACT(EPOCH FROM (ls.end_time - ls.start_time)) / 60), 0)::float8
		FROM location_schedules ls
		JOIN locations l ON l.id = ls.location_id
		WHERE l.organization_id = $1
			AND (CARDINALITY($2::uuid[]) = 0 OR ls.location_id = ANY($2))
			AND ls.date >= $3::date AND ls.date <= $4::date
			AND ls.type = 'work'
	`

	getEmployeeScheduleMinutes = `
		SELECT COALESCE(SUM(EXTRACT(EPOCH FROM (end_time - start_time)) / 60), 0)::float8
		FROM employee_schedules
		WHERE employee_id = $1 AND date >= $2::date AND date <= $3::date AND type = 'work'
	`

	getEmployeeScheduleMinutesByEmployee = `
		SELECT
			es.employee_id AS employee_id,
			COALESCE(SUM(EXTRACT(EPOCH FROM (es.end_time - es.start_time)) / 60), 0)::float8 AS minutes
		FROM employee_schedules es
		JOIN employees e ON e.id = es.employee_id
		WHERE e.organization_id = $1
			AND (CARDINALITY($2::uuid[]) = 0 OR e.location_id = ANY($2))
			AND es.date >= $3::date AND es.date <= $4::date AND es.type = 'work'
		GROUP BY es.employee_id
	`

	// Полная история по клиентам (без периода) — основа для сегментов/retention/когорт.
	getCustomerStats = `
		SELECT
			a.customer_id              AS customer_id,
			MIN(a.start_at)            AS first_visit,
			MAX(a.start_at)            AS last_visit,
			COUNT(*)                   AS total_visits,
			COALESCE(AVG(a.price), 0)::float8 AS avg_check
		FROM appointments a
		WHERE a.organization_id = $1
			AND (CARDINALITY($2::uuid[]) = 0 OR a.location_id = ANY($2))
			AND ($3::uuid IS NULL OR a.employee_id = $3)
			AND a.status = 'completed'
			AND a.customer_id IS NOT NULL
		GROUP BY a.customer_id
	`

	// Клиенты, посетившие (завершённые записи) в периоде.
	getCustomersVisited = `
		SELECT DISTINCT customer_id
		FROM appointments
		WHERE ` + scopeFilter + ` AND status = 'completed' AND customer_id IS NOT NULL
	`

	// Строки отчёта по мастерам: включает мастеров без активности (LEFT JOIN).
	getStaffRows = `
		SELECT
			e.id        AS employee_id,
			e.full_name AS full_name,
			COALESCE(SUM(a.price) FILTER (WHERE a.status = 'completed'), 0)::float8           AS revenue,
			COUNT(a.id) FILTER (WHERE a.status = 'completed')                                 AS bookings,
			COALESCE(SUM(a.duration_minutes) FILTER (WHERE a.status = 'completed'), 0)::float8 AS duration_minutes,
			COUNT(a.id) FILTER (WHERE a.status = 'cancelled')                                 AS cancelled,
			COUNT(a.id)                                                                       AS created
		FROM employees e
		LEFT JOIN appointments a
			ON a.employee_id = e.id
			AND a.start_at >= $3 AND a.start_at < $4::date + INTERVAL '1 day'
		WHERE e.organization_id = $1
			AND e.deleted_at IS NULL
			AND e.can_provide_services = true
			AND (CARDINALITY($2::uuid[]) = 0 OR e.location_id = ANY($2))
		GROUP BY e.id, e.full_name
		ORDER BY revenue DESC
	`

	getStaffWeekdayLoad = `
		SELECT
			a.employee_id                                       AS employee_id,
			EXTRACT(DOW FROM a.start_at AT TIME ZONE 'UTC')::int AS weekday,
			COUNT(*)                                            AS cnt
		FROM appointments a
		JOIN employees e ON e.id = a.employee_id
		WHERE a.organization_id = $1
			AND (CARDINALITY($2::uuid[]) = 0 OR e.location_id = ANY($2))
			AND a.start_at >= $3 AND a.start_at < $4::date + INTERVAL '1 day'
			AND a.status = 'completed'
			AND e.can_provide_services = true
		GROUP BY a.employee_id, weekday
	`

	getEmployeeInfo = `
		SELECT id AS id, full_name AS full_name, avatar_id AS avatar_id
		FROM employees
		WHERE id = $1 AND organization_id = $2 AND deleted_at IS NULL
	`

	getEmployeeFinance = `
		SELECT
			e.id                AS employee_id,
			e.full_name         AS full_name,
			e.commission_percent AS commission_percent,
			COALESCE(SUM(a.price) FILTER (WHERE a.status = 'completed'), 0)::float8 AS revenue
		FROM employees e
		LEFT JOIN appointments a
			ON a.employee_id = e.id
			AND a.start_at >= $3 AND a.start_at < $4::date + INTERVAL '1 day'
			AND (CARDINALITY($2::uuid[]) = 0 OR a.location_id = ANY($2))
		WHERE e.organization_id = $1
			AND e.deleted_at IS NULL
			AND e.can_provide_services = true
			AND (CARDINALITY($2::uuid[]) = 0 OR e.location_id = ANY($2))
		GROUP BY e.id, e.full_name, e.commission_percent
		ORDER BY revenue DESC
	`

	// Проверка принадлежности локации организации.
	getLocationBelongsToOrg = `
		SELECT EXISTS(
			SELECT 1 FROM locations WHERE id = $1 AND organization_id = $2
		)
	`
)
