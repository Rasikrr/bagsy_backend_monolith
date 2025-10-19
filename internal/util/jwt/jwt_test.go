package jwt

// func TestValidateToken(t *testing.T) {
// 	secret := "test"
// 	params := &entity.PayloadParams{
// 		Phone:   "test",
// 		Role:    "test",
// 		Active:  true,
// 		Refresh: false,
// 	}
// 	token, err := GenerateAccessToken(params, secret)
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}
// 	valid, err := ValidateToken(token, secret)
// 	if err != nil {
// 		t.Errorf("expected no error, got %v", err)
// 	}
// 	if !valid {
// 		t.Errorf("expected valid, got %v", valid)
// 	}
// }

// func TestGenTokens(t *testing.T) {
// 	params := &entity.PayloadParams{
// 		Phone:  "77777777777",
// 		Role:   "admin",
// 		Active: true,
// 	}
// 	secret := "secret"
// 	access, err := GenerateAccessToken(params, secret)
// 	require.NoError(t, err)
// 	refresh, err := GenerateRefreshToken(params, secret)
// 	require.NoError(t, err)
// 	t.Log("Access token:", access)
// 	t.Log("Refresh token:", refresh)
// }
