package main

// This exists so we can have unique key to store and retrieve the authentication status
// from a request context without the risk of naming collisions
// as some third-party packages might use the same "isAuthenticated" key
type contextKey string

const isAuthenticatedContextKey = contextKey("isAuthenticated")
