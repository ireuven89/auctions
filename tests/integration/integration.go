package integration

/*
func setupTestServer() http.Handler {
	// Replace this with your actual server setup
	return auth.NewServer()
}

func TestSignupAndLoginFlow(t *testing.T) {
	server := setupTestServer()

	// 1️⃣ Signup request
	signupPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	signupBody, _ := json.Marshal(signupPayload)

	signupReq := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(signupBody))
	signupReq.Header.Set("Content-Type", "application/json")
	signupRec := httptest.NewRecorder()

	server.ServeHTTP(signupRec, signupReq)

	if signupRec.Code != http.StatusCreated {
		t.Fatalf("signup failed: status %d, body %s", signupRec.Code, signupRec.Body.String())
	}

	// 2️⃣ Login request
	loginPayload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	loginBody, _ := json.Marshal(loginPayload)

	loginReq := httptest.NewRequest(http.MethodPost, "/login", bytes.NewReader(loginBody))
	loginReq.Header.Set("Content-Type", "application/json")
	loginRec := httptest.NewRecorder()

	server.ServeHTTP(loginRec, loginReq)

	if loginRec.Code != http.StatusOK {
		t.Fatalf("login failed: status %d, body %s", loginRec.Code, loginRec.Body.String())
	}

	var loginResp map[string]interface{}
	if err := json.Unmarshal(loginRec.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("invalid login response JSON: %v", err)
	}

	token, ok := loginResp["token"].(string)
	if !ok || token == "" {
		t.Fatalf("missing or invalid token in response: %+v", loginResp)
	}

	t.Logf("received token: %s", token)
}*/
