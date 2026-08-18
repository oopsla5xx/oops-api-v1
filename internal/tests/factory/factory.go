// Package factory provides FactoryBot-style test data builders.
// Use gofakeit for generating realistic fake values, and the functional
// overrides pattern to customize individual fields.
//
// # Convention for adding entity factories
//
// When a new domain module is added, create a corresponding factory file:
//
//	internal/tests/factory/<entity>.go
//
// Each factory file exposes a New<Entity> function using the overrides pattern:
//
//	type UserInput struct {
//	    Email string
//	    Name  string
//	}
//
//	func NewUser(overrides ...func(*UserInput)) UserInput {
//	    u := UserInput{
//	        Email: gofakeit.Email(),
//	        Name:  gofakeit.Name(),
//	    }
//	    for _, o := range overrides {
//	        o(&u)
//	    }
//	    return u
//	}
//
// Usage in tests:
//
//	user := factory.NewUser()                               // all defaults
//	admin := factory.NewUser(func(u *factory.UserInput) {   // with override
//	    u.Email = "admin@example.com"
//	})
//
// # Determinism
//
// Call factory.Seed in TestMain when reproducible data is required:
//
//	func TestMain(m *testing.M) {
//	    factory.Seed(42)
//	    os.Exit(m.Run())
//	}
package factory

import "github.com/brianvoe/gofakeit/v6"

// Seed initializes the global random seed for reproducible test data.
// Omit (or pass 0) for random data on each run.
func Seed(n int64) {
	gofakeit.Seed(n)
}
