package namespace

import (
	"context"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/metadata"
)

const (
	DefaultNamespace = "go.micro"
	// NamespaceKey is used to set/get the namespace from the context
	NamespaceKey = "Micro-Namespace"
)

// FromContext gets the namespace from the context
func FromContext(ctx context.Context) string {
	// get the namespace which is set at ingress by micro web / api / proxy etc
	ns, ok := metadata.Get(ctx, NamespaceKey)
	if !ok || ns == DefaultNamespace {
		return DefaultNamespace
	}

	// get the account making the request. if there is no account then we return the namespace
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return ns
	}

	// if the account has the scope, return the namespace
	if acc.HasScope("namespace", ns) {
		return ns
	}

	// allow the runtime access to all namespaces.
	// TODO: grant runtime services elevated privelages and validate them here instead of assuming all
	// services in the default namespace are the runtime.
	if acc.HasScope("namespace", DefaultNamespace) {
		return ns
	}

	// a forbidden cross namespace request was made, return the default instead of the one requested
	return DefaultNamespace
}
