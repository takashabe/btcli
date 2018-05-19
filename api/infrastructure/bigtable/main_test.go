package bigtable

import (
	"log"
	"os"
	"sync"
	"testing"

	fixture "github.com/takashabe/bt-fixture"
)

func TestMain(m *testing.M) {
	project := getEnvWithDefault("BTCLI_PROJECT", "test-project")
	instance := getEnvWithDefault("BTCLI_INSTANCE", "test-instance")
	connect(project, instance)

	os.Exit(m.Run())
}

func getEnvWithDefault(env, def string) string {
	act := os.Getenv(env)
	if len(act) == 0 {
		return def
	}
	return act
}

var (
	fixtureClient *fixture.Fixture
	fixtureOnce   sync.Once
)

func connect(project, instance string) {
	fixtureOnce.Do(func() {
		fix, err := fixture.NewFixture(project, instance)
		if err != nil {
			log.Fatalf("failed to initialize fixture client. %v", err)
		}
		fixtureClient = fix
	})
}

func loadFixture(t *testing.T, f string) {
	err := fixtureClient.Load(f)
	if err != nil {
		t.Fatalf("failed to load fixture. %v", err)
	}
}
