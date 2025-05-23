package internal

/*func TestService_Register_Success(t *testing.T) {
	ctx := context.Background()
	mockRepository := &mocks.MockRepo{
		CreateUserFunc: func(ctx context.Context, u user.User) error { return nil },
	}
	svc := &service{
		repository:   mockRepository,
		hashPassword: func(pw string) (string, error) { return "hashedpw", nil },
		generateID:   func() string { return "user123" },
		SignToken: func(ctx context.Context, u user.User) (string, error) {
			return "token123", nil
		},
		GenerateRefreshToken: func(ctx context.Context, id string) (string, error) {
			return "refresh456", nil
		},
	}
	userInput := user.User{Password: "plainpw"}
	token, refreshToken, err := svc.Register(ctx, userInput)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if token != "token123" {
		t.Errorf("expected token123, got %v", token)
	}
	if refreshToken != "refresh456" {
		t.Errorf("expected refresh456, got %v", refreshToken)
	}
}

func TestService_Register_HashPasswordError(t *testing.T) {
	ctx := context.Background()
	svc := &service{
		repository:           &mocks.MockRepo{},
		hashPassword: func(pw string) (string, error) { return "", errors.New("hash error") },
		generateID:   func() string { return "user123" },
		SignToken:    func(ctx context.Context, u user.User) (string, error) { return "", nil },
		GenerateRefreshToken: func(ctx context.Context, id string) (string, error) { return "", nil },
	}
	}
	userInput := user.User{Password: "plainpw"}
	_, _, err := svc.Register(ctx, userInput)
	if err == nil || err.Error() == "" {
		t.Fatal("expected hash error, got nil")
	}
}
*/
