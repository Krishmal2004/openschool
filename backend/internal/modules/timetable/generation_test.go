package timetable

import (
	"testing"

	"github.com/google/uuid"
)

func TestGenerationSlotsExcludeBreaksAndBuildAdjacentPairs(t *testing.T) {
	first, second, breakNumber := int32(1), int32(2), int32(3)
	slots, doubles := generationSlots([]generationPeriod{
		{Number: &first, SlotType: "period"},
		{Number: &second, SlotType: "period"},
		{Number: &breakNumber, SlotType: "break"},
	})
	if len(slots) != 10 {
		t.Fatalf("slots = %d, want 10 weekday slots", len(slots))
	}
	if len(doubles) != 5 || doubles[0][0].Period != 1 || doubles[0][1].Period != 2 {
		t.Fatalf("double slots = %+v, want five 1-2 pairs", doubles)
	}
}

func TestGenerationShuffleIsStableForSameClassAndSubject(t *testing.T) {
	classID, subjectID := uuid.New(), uuid.New()
	slots := []generationSlot{{Day: 1, Period: 1}, {Day: 2, Period: 1}, {Day: 3, Period: 1}}
	first := shuffledGenerationSlots(slots, classID, subjectID)
	second := shuffledGenerationSlots(slots, classID, subjectID)
	for i := range first {
		if first[i] != second[i] {
			t.Fatalf("shuffle is not deterministic: first=%v second=%v", first, second)
		}
	}
}
