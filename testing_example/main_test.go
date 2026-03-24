package main

import (
	"errors"
	"testing"
)

// Dummy Test
func TestNotifyUser_WithDummy(t *testing.T) {
	sender := DummyEmailSender{}

	err := NotifyUser(&sender, "test@example.com")

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

// Stub Test
type ErrorStubEmailSender struct{}

func (s ErrorStubEmailSender) Send(to, msg string) error {
	return errors.New("stub error")
}

func TestNotifyUser_WithStub(t *testing.T) {
	sender := ErrorStubEmailSender{}

	err := NotifyUser(sender, "test@example.com")

	if err == nil {
		t.Errorf("expected error, got nil")
	}
}

// Mock Test
func TestNotifyUser_WithMock(t *testing.T) {
	mock := &MockEmailSender{}

	_ = NotifyUser(mock, "test@example.com")

	if !mock.Called {
		t.Errorf("expected Send to be called")
	}

	if mock.To != "test@example.com" {
		t.Errorf("expected recipient to be test@example.com, got %s", mock.To)
	}

	if mock.Msg != "Welcome !" {
		t.Errorf("unexpected message: %s", mock.Msg)
	}
}

// Fake Test
func TestNotifyUser_WithFake(t *testing.T) {
	fake := &FakeEmailSender{}

	_ = NotifyUser(fake, "test@example.com")

	if len(fake.sent) != 1 {
		t.Errorf("expected 1 email sent, got %d", len(fake.sent))
	}

	if fake.sent[0] != "Welcome !" {
		t.Errorf("unexpected message: %s", fake.sent[0])
	}
}

// Spy Test
func TestNotifyUser_WithSpy(t *testing.T) {
	fake := &FakeEmailSender{}
	spy := &SpiesEmailSender{
		Real: fake,
	}

	_ = NotifyUser(spy, "test@example.com")

	if !spy.Called {
		t.Errorf("expected spy to record call")
	}

	if len(fake.sent) != 1 {
		t.Errorf("expected underlying sender to be called")
	}
}
