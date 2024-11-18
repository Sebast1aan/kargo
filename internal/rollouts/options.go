package rollouts

import (
	"maps"

	"k8s.io/apimachinery/pkg/types"
)

const (
	// ulidLength is the length of a ulid.ULID string.
	ulidLength = 26
	// maxNameSuffixLength is the maximum length of the name suffix for an
	// AnalysisRun. It assumes that the suffix contains e.g. a SHA and can
	// be truncated to a smaller length than the maxNamePrefixLength.
	maxNameSuffixLength = 7
	// maxNamePrefixLength is the maximum length of the name prefix for an
	// AnalysisRun. It takes into account the maximum length of the name
	// field (253 characters), and the additional characters that will be
	// appended to the name (ULID, SHA, and period separators).
	maxNamePrefixLength = 253 - (1 + ulidLength) - (1 + maxNameSuffixLength)
)

// AnalysisRunOption is an option for configuring the build of an AnalysisRun.
type AnalysisRunOption interface {
	ApplyToAnalysisRun(*AnalysisRunOptions)
}

// AnalysisRunOptions holds the options for building an AnalysisRun.
type AnalysisRunOptions struct {
	NamePrefix  string
	NameSuffix  string
	ExtraLabels map[string]string
	Owners      []Owner
}

// Owner represents a reference to an owner object.
type Owner struct {
	APIVersion    string
	Kind          string
	Reference     types.NamespacedName
	BlockDeletion bool
}

// Apply applies the given options to the AnalysisRunOptions.
func (o *AnalysisRunOptions) Apply(opts ...AnalysisRunOption) {
	for _, opt := range opts {
		opt.ApplyToAnalysisRun(o)
	}
}

// WithNamePrefix sets the name prefix for the AnalysisRun. If it is longer
// than maxNamePrefixLength, it will be truncated.
type WithNamePrefix string

func (o WithNamePrefix) ApplyToAnalysisRun(opts *AnalysisRunOptions) {
	prefix := o
	if len(prefix) > maxNamePrefixLength {
		prefix = prefix[0:maxNamePrefixLength]
	}
	opts.NamePrefix = string(prefix)
}

// WithNameSuffix sets the name suffix for the AnalysisRun. If it is longer
// than maxNameSuffixLength, it will be truncated.
type WithNameSuffix string

func (o WithNameSuffix) ApplyToAnalysisRun(opts *AnalysisRunOptions) {
	suffix := o
	if len(suffix) > maxNameSuffixLength {
		suffix = suffix[0:maxNameSuffixLength]
	}
	opts.NameSuffix = string(suffix)
}

// WithExtraLabels sets the extra labels for the AnalysisRun. It can be passed
// multiple times to add more labels.
type WithExtraLabels map[string]string

func (o WithExtraLabels) ApplyToAnalysisRun(opts *AnalysisRunOptions) {
	if opts.ExtraLabels != nil {
		maps.Copy(opts.ExtraLabels, o)
		return
	}
	opts.ExtraLabels = o
}

// WithOwner sets the owner for the AnalysisRun. It can be passed multiple times
// to add more owners.
type WithOwner Owner

func (o WithOwner) ApplyToAnalysisRun(opts *AnalysisRunOptions) {
	opts.Owners = append(opts.Owners, Owner(o))
}