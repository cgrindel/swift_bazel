package gazelle

import (
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const languageName = "swift"

const swiftLibraryRule = "swift_library"
const swiftBinaryRule = "swift_binary"
const swiftTestRule = "swift_test"

var kindShared = rule.KindInfo{
	NonEmptyAttrs:  map[string]bool{"srcs": true, "deps": true},
	MergeableAttrs: map[string]bool{"srcs": true},
}

var kinds = map[string]rule.KindInfo{
	swiftLibraryRule: kindShared,
	swiftBinaryRule:  kindShared,
	swiftTestRule:    kindShared,
}

var loads = []rule.LoadInfo{
	{
		Name: "@build_bazel_rules_swift//swift:swift.bzl",
		Symbols: []string{
			swiftLibraryRule,
			swiftBinaryRule,
			swiftTestRule,
		},
	},
}

type swiftLang struct {
	language.BaseLang
}

func NewLanguage() language.Language {
	return &swiftLang{}
}

func (*swiftLang) Name() string { return languageName }

func (*swiftLang) Kinds() map[string]rule.KindInfo {
	return kinds
}

func (*swiftLang) Loads() []rule.LoadInfo {
	return loads
}

// func (sl *swiftLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
// }

// type swiftModuleCollector struct {
// 	ModuleFiles map[string][]string
// }

// func (l *swiftLang) Resolve(
// 	c *config.Config,
// 	ix *resolve.RuleIndex,
// 	rc *repo.RemoteCache,
// 	r *rule.Rule,
// 	imports interface{},
// 	from label.Label) {
// }

// func (*swiftLang) Configure(c *config.Config, rel string, f *rule.File) {
// }

// // Imports returns a list of ImportSpecs that can be used to import the rule
// // r. This is used to populate RuleIndex.
// //
// // If nil is returned, the rule will not be indexed. If any non-nil slice is
// // returned, including an empty slice, the rule will be indexed.
// func (*swiftLang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
// 	srcs := r.AttrStrings("srcs")
// 	imports := make([]resolve.ImportSpec, 0, len(srcs))
//
// 	for _, src := range srcs {
// 		spec := resolve.ImportSpec{
// 			// Lang is the language in which the import string appears (this should
// 			// match Resolver.Name).
// 			Lang: languageName,
// 			// Imp is an import string for the library.
// 			Imp: fmt.Sprintf("//%s:%s", f.Pkg, src),
// 		}
//
// 		imports = append(imports, spec)
// 	}
//
// 	return imports
// }
