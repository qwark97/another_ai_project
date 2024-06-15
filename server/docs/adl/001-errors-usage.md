# Errors usage

Errors are returned by the methods / functions which depend on external services (DB, API requests etc.).

The exception are the methods / functions which are based on the AI - these calls, in case of error return the message (string) that informs about the error.

Initialisation of the components that happens at the beginning of the life of the service, in case of failure causes the server to panic
