package feedbackend

// import (
// 	"errors"
// 	"strings"
// 	"sync"
// 	"testing"
// 	"time"

// 	"github.com/alicebob/miniredis/v2"
// 	"github.com/bluecolor/tractor/pkg/lib/msg"
// )

// func TestNew(t *testing.T) {
// 	t.Parallel()
// 	mr := miniredis.RunT(t)
// 	if _, err := New(mr.Addr()); err != nil {
// 		t.Error(err)
// 	}
// }

// func TestClose(t *testing.T) {
// 	t.Parallel()
// 	mr := miniredis.RunT(t)
// 	r, err := New(mr.Addr())
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if err := r.Close(); err != nil {
// 		t.Error(err)
// 	}
// }

// func TestStore(t *testing.T) {
// 	t.Parallel()
// 	mr := miniredis.RunT(t)
// 	r, err := New(mr.Addr())
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	sessionID := "sessionID"
// 	feedbacks := []*msg.Feedback{
// 		msg.NewInputProgress(10),
// 		msg.NewOutputProgress(5),
// 		msg.NewError(msg.InputConnector, errors.New("error")),
// 	}
// 	for _, f := range feedbacks {
// 		if err := r.Process(sessionID, f); err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	session, err := r.GetSession(sessionID)
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	if session[StatusKey] != strings.ToLower(msg.Error.String()) {
// 		t.Errorf("status is %s, want %s", session[StatusKey], msg.Error.String())
// 	}
// 	if session[InputProgressKey] != "10" {
// 		t.Errorf("input progress is %s, want %s", session[InputProgressKey], "10")
// 	}
// 	if session[OutputProgressKey] != "5" {
// 		t.Errorf("output progress is %s, want %s", session[OutputProgressKey], "5")
// 	}
// }

// func TestPubsub(t *testing.T) {
// 	t.Parallel()
// 	mr := miniredis.RunT(t)
// 	r, err := New(mr.Addr())
// 	if err != nil {
// 		t.Error(err)
// 	}
// 	sessionID := "sessionID"
// 	feedbacks := []*msg.Feedback{
// 		msg.NewInputProgress(10),
// 		msg.NewOutputProgress(5),
// 		msg.NewError(msg.InputConnector, errors.New("error")),
// 	}
// 	wg := &sync.WaitGroup{}
// 	wg.Add(1)
// 	go func(wg *sync.WaitGroup, t *testing.T) {
// 		defer wg.Done()
// 		pubsub, channel, err := r.Subscribe(sessionID)
// 		if err != nil {
// 			t.Error(err)
// 		}
// 		feedbacks := []*msg.Feedback{}
// 		for feedback := range channel {
// 			if feedback == nil {
// 				break
// 			}
// 			feedbacks = append(feedbacks, feedback)
// 		}
// 		if err := pubsub.Unsubscribe(); err != nil {
// 			t.Error(err)
// 		}
// 		if len(feedbacks) != 3 {
// 			t.Errorf("feedbacks length is %d, want %d", len(feedbacks), 3)
// 		}
// 		// todo: check feedbacks
// 	}(wg, t)

// 	// wait client connection to ensure ordering
// 	// https://stackoverflow.com/a/45076713/3401812
// 	time.Sleep(time.Millisecond * 100)
// 	for _, f := range feedbacks {
// 		if err := r.Publish(sessionID, f); err != nil {
// 			t.Error(err)
// 		}
// 	}
// 	if err := r.Publish(sessionID, nil); err != nil {
// 		t.Error(err)
// 	}
// 	wg.Wait()
// }
