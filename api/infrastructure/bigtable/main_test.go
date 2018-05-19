package bigtable

import (
	"os"
	"sync"
	"testing"

	fixture "github.com/takashabe/bt-fixture"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

var (
	fixtureClient *fixture.Fixture
	fixtureOnce   sync.Once
)

func loadFixture(t *testing.T, f string) {
	fixtureOnce.Do(func() {
		fix, err := fixture.NewFixture("test-project", "test-instance")
		if err != nil {
			t.Fatalf("failed to initialize fixture client. %v", err)
		}
		fixtureClient = fix
	})
	err := fixtureClient.Load(f)
	if err != nil {
		t.Fatalf("failed to load fixture. %v", err)
	}
}
