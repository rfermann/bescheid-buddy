package reply_service

import (
	"encoding/json"
	"testing"
)

func TestGenerationInputJSONRoundTrip(t *testing.T) {
	original := GenerationInput{Instruction: "Bitte widersprechen."}

	encoded, err := json.Marshal(original)
	if err != nil {
		t.Fatal(err)
	}
	if string(encoded) != `{"instruction":"Bitte widersprechen."}` {
		t.Fatalf("unexpected JSON: %s", encoded)
	}

	var decoded GenerationInput
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatal(err)
	}
	if decoded != original {
		t.Fatalf("unexpected round trip: %#v", decoded)
	}
}
