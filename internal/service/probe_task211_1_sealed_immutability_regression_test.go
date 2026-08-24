package service

import (
	"testing"

	"task211-quenchprop/internal/model"
	"task211-quenchprop/internal/store"
)

func TestBug01_SealedExperimentRemainsImmutable(t *testing.T) {
	db, err := store.Open(t.TempDir() + "/sealed.db")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	app := NewApp(db, DefaultConfig())
	if err := app.AddTopology(&model.CoilTopology{ID: "topo-sealed", Name: "topology"}); err != nil {
		t.Fatal(err)
	}
	exp, err := app.CreateExperiment(ExperimentInput{
		Name: "sealed", MagnetSN: "mag-1", Operator: "op", TopologyID: "topo-sealed", SampleRate: 1000,
	})
	if err != nil {
		t.Fatal(err)
	}
	ch := &model.Channel{Index: 1, SegmentID: "seg-1", Status: model.ChanOnline}
	if err := app.AddChannel(exp.ID, ch); err != nil {
		t.Fatal(err)
	}
	for _, next := range []string{model.ExpStatusCollecting, model.ExpStatusAnalyzing, model.ExpStatusConfirmed, model.ExpStatusSealed} {
		if _, err := app.TransitionExperiment(exp.ID, next); err != nil {
			t.Fatalf("transition to %s: %v", next, err)
		}
	}
	if _, err := app.IngestWaveform(exp.ID, ch.ID, 1000, []model.Sample{{RawT: 0, V: 0}, {RawT: 0.001, V: 1}}); err != model.ErrSealed {
		t.Fatalf("sealed waveform write error = %v, want %v", err, model.ErrSealed)
	}
	if err := app.AddChannel(exp.ID, &model.Channel{Index: 2, SegmentID: "seg-2"}); err != model.ErrSealed {
		t.Fatalf("sealed channel write error = %v, want %v", err, model.ErrSealed)
	}
	sealed, err := app.GetExperiment(exp.ID)
	if err != nil {
		t.Fatal(err)
	}
	sealed.TopologyID = "topo-other"
	if err := app.UpdateExperimentTopology(sealed); err != model.ErrSealed {
		t.Fatalf("sealed topology update error = %v, want %v", err, model.ErrSealed)
	}
}
