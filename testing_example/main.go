package main

type EmailSender interface {
	Send(to, msg string) error
}

func NotifyUser(sender EmailSender, user string) error {
	return sender.Send(user, "Welcome !")
}

// Dummy
type DummyEmailSender struct{}

func (d DummyEmailSender) Send(to, msg string) error {
	return nil // never actually used
}

// Stub
type StubEmailSender struct{}

func (s StubEmailSender) Send(to, msg string) error {
	return nil // can be anythin fixed thing depending on return type
}

// Mocks
type MockEmailSender struct {
	Called bool
	To     string
	Msg    string
}

func (m MockEmailSender) Send(to, msg string) error {
	m.Called = true // specific set of steps being followed in a order
	m.To = to
	m.Msg = msg
	return nil
}

// fake
type FakeEmailSender struct {
	sent []string
}

func (f FakeEmailSender) Send(to, msg string) error {
	f.sent = append(f.sent, msg) // simplified implementation of in memory email backup
	return nil
}

// spies
type SpiesEmailSender struct {
	Real   EmailSender
	Called bool
}

func (s SpiesEmailSender) Send(to, msg string) error {
	s.Called = true // observe flow
	return s.Real.Send(to, msg)
}
