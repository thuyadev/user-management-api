package feature_test

import (
	"net/http"
	"testing"

	"user-management-api/tests/testutil"
)

func TestLogStatsEndpoints(t *testing.T) {
	app := testutil.SetupTestApp(t)
	adminToken := app.LoginToken(t, app.Admin.Email, testutil.TestPassword)

	eventResp := app.Request(http.MethodGet, "/api/v1/admin/logs/stats/events?days=30", nil, adminToken)
	if eventResp.Code != http.StatusOK {
		t.Fatalf("event stats: status %d body %s", eventResp.Code, eventResp.Body.String())
	}

	dailyResp := app.Request(http.MethodGet, "/api/v1/admin/logs/stats/daily?days=7", nil, adminToken)
	if dailyResp.Code != http.StatusOK {
		t.Fatalf("daily stats: status %d body %s", dailyResp.Code, dailyResp.Body.String())
	}

	userToken := app.LoginToken(t, app.User.Email, testutil.TestPassword)
	forbidden := app.Request(http.MethodGet, "/api/v1/admin/logs/stats/events?days=30", nil, userToken)
	if forbidden.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for user role, got %d", forbidden.Code)
	}
}
