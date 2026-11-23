package ir

import "testing"

func TestSpanUsesHalfOpenOffsets(t *testing.T) {
	span := Span{File: FileID(1), Start: 4, End: 9}
	if !span.Valid() {
		t.Fatal("ordered span with a valid file must be valid")
	}
	if (Span{File: FileID(1), Start: 9, End: 4}).Valid() {
		t.Fatal("span with end before start must be invalid")
	}
	if (Span{File: FileID(1), Start: 4, End: 4}).Valid() != true {
		t.Fatal("empty half-open span should be valid")
	}
}

func TestOriginRequiresPrimarySpan(t *testing.T) {
	origin := Origin{Primary: Span{File: FileID(1), Start: 2, End: 3}}
	if !origin.Valid() {
		t.Fatal("origin with a valid primary span must be valid")
	}
	if (Origin{}).Valid() {
		t.Fatal("origin without a primary span must be invalid")
	}
}
