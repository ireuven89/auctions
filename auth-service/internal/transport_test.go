package internal

import (
	"bytes"
	"context"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"

	user2 "github.com/ireuven89/auctions/auth-service/user"
	"github.com/ireuven89/auctions/shared/jwksprovider"
)

func TestDecodeRegisterUserRequest_Success(t *testing.T) {
	user := user2.User{Name: "testuser", Password: "testpass"}
	body, _ := json.Marshal(user)
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))

	req, err := decodeRegisterUserRequest(context.Background(), r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	rr, ok := req.(RegisterUserRequest)
	if !ok {
		t.Fatal("expected RegisterUserRequest type")
	}
	if rr.user.Name != user.Name || rr.user.Password != user.Password {
		t.Errorf("decoded user mismatch: got %+v, want %+v", rr.user, user)
	}
}

func TestDecodeRegisterUserRequest_InvalidJSON(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("notjson")))
	_, err := decodeRegisterUserRequest(context.Background(), r)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestEncodeRegisterUserResponse_Success(t *testing.T) {
	resp := RegisterUserResponse{AccessToken: "a", RefreshToken: "r"}
	w := httptest.NewRecorder()
	err := encodeRegisterUserResponse(context.Background(), w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected Content-Type json, got %s", ct)
	}
	var m map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["accessToken"] != "a" || m["refreshToken"] != "r" {
		t.Errorf("unexpected output: %+v", m)
	}
}

func TestEncodeRegisterUserResponse_BadType(t *testing.T) {
	w := httptest.NewRecorder()
	err := encodeRegisterUserResponse(context.Background(), w, struct{}{})
	if err == nil || err.Error() != "encodeRegisterUserResponse.failed casting response" {
		t.Errorf("expected type error, got %v", err)
	}
}

func TestDecodeLoginRequest_Success(t *testing.T) {
	login := LoginRequestModel{Identifier: "u", Password: "p"}
	body, _ := json.Marshal(login)
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req, err := decodeLoginRequest(context.Background(), r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := req.(LoginRequestModel)
	if got.Identifier != "u" || got.Password != "p" {
		t.Errorf("unexpected login: %+v", got)
	}
}

func TestDecodeLoginRequest_InvalidJSON(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{")))
	_, err := decodeLoginRequest(context.Background(), r)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestEncodeLoginUserResponse_Success(t *testing.T) {
	resp := LoginResponseModel{AccessToken: "tok", RefreshToken: "ref"}
	w := httptest.NewRecorder()
	err := encodeLoginUserResponse(context.Background(), w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]string
	if err = json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["token"] != "tok" || m["refreshToken"] != "ref" {
		t.Errorf("unexpected output: %+v", m)
	}
}

func TestEncodeLoginUserResponse_BadType(t *testing.T) {
	w := httptest.NewRecorder()
	err := encodeLoginUserResponse(context.Background(), w, struct{}{})
	if err == nil || err.Error() != "encodeLoginUserResponse.failed casting response" {
		t.Errorf("expected type error, got %v", err)
	}
}

func TestDecodeGetPublicRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	req, err := decodeGetPublicRequest(context.Background(), r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if req != nil {
		t.Errorf("expected nil, got %v", req)
	}
}

func TestEncodeGetPublicResponse_Success(t *testing.T) {
	resp := GetPublicKeyResponse{
		PublicKey: jwksprovider.JWKS{
			Keys: []json.RawMessage{json.RawMessage(`"k"`)},
		},
	}

	w := httptest.NewRecorder()
	err := encodeGetPublicResponse(context.Background(), w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var m map[string]interface{}
	if err = json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}

	jwks, ok := m["jwks"].([]interface{})
	if !ok {
		t.Errorf("jwks key missing or wrong type: %+v", m)
	}

	if len(jwks) != 1 {
		t.Errorf("expected 1 jwk, got %d", len(jwks))
	}

	// Check if the first key matches what we put in
	if jwks[0] != "k" {
		t.Errorf("unexpected jwk value: got %+v, want %q", jwks[0], "k")
	}
}

func TestEncodeGetPublicResponse_BadType(t *testing.T) {
	w := httptest.NewRecorder()
	err := encodeGetPublicResponse(context.Background(), w, struct{}{})
	if err == nil || err.Error() != "encodeGetPublicResponse failed casting response" {
		t.Errorf("expected type error, got %v", err)
	}
}

func TestDecodeRefreshRequest_Success(t *testing.T) {
	refresh := RefreshRequestModel{Refresh: "refresh"}
	body, _ := json.Marshal(refresh)
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(body))
	req, err := decodeRefreshRequest(context.Background(), r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got := req.(RefreshRequestModel)
	if got.Refresh != "refresh" {
		t.Errorf("unexpected refresh: %+v", got)
	}
}

func TestDecodeRefreshRequest_InvalidJSON(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader([]byte("{")))
	_, err := decodeRefreshRequest(context.Background(), r)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestEncodeRefreshResponse_Success(t *testing.T) {
	resp := RefreshResponseModel{AccessToken: "tok"}
	w := httptest.NewRecorder()
	err := encodeRefreshResponse(context.Background(), w, resp)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var m map[string]string
	if err := json.Unmarshal(w.Body.Bytes(), &m); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if m["token"] != "tok" {
		t.Errorf("unexpected output: %+v", m)
	}
}
