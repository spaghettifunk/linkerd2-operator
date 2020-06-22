package v1alpha1

// ConfigState describes the state of the operator
type ConfigState string

const (
	Created         ConfigState = "Created"
	ReconcileFailed ConfigState = "ReconcileFailed"
	Reconciling     ConfigState = "Reconciling"
	Available       ConfigState = "Available"
	Unmanaged       ConfigState = "Unmanaged"
)
