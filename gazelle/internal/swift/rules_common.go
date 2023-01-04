package swift

import (
	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

// Rule Creation

func rulesForLibraryModule(
	defaultModuleName string,
	srcs []string,
	swiftImports []string,
	shouldSetVis bool,
	buildFile *rule.File,
) []*rule.Rule {
	var r *rule.Rule
	existingLibs := findRulesByKind(buildFile.Rules, LibraryRuleKind)
	// If we found a single swift_library, use its name. Otherwise, just use the module name.
	var name, moduleName string
	if len(existingLibs) == 1 {
		first := existingLibs[0]
		name = first.Name()
		moduleName = first.AttrString(ModuleNameAttrName)
	} else {
		name = defaultModuleName
		moduleName = defaultModuleName
	}
	r = rule.NewRule(LibraryRuleKind, name)
	setCommonSwiftAttrs(r, moduleName, srcs, swiftImports)
	setVisibilityAttr(r, shouldSetVis, []string{"//visibility:public"})
	return []*rule.Rule{r}
}

func rulesForBinaryModule(
	defaultModuleName string,
	srcs []string,
	swiftImports []string,
	shouldSetVis bool,
	buildFile *rule.File,
) []*rule.Rule {
	var r *rule.Rule
	existingLibs := findRulesByKind(buildFile.Rules, BinaryRuleKind)
	// If we found a single swift_binary, use its name. Otherwise, just use the module name.
	var name, moduleName string
	if len(existingLibs) == 1 {
		first := existingLibs[0]
		name = first.Name()
		moduleName = first.AttrString(ModuleNameAttrName)
	} else {
		name = defaultModuleName
		moduleName = defaultModuleName
	}
	r = rule.NewRule(BinaryRuleKind, name)
	setCommonSwiftAttrs(r, moduleName, srcs, swiftImports)
	setVisibilityAttr(r, shouldSetVis, []string{"//visibility:public"})
	return []*rule.Rule{r}
}

func rulesForTestModule(
	defaultModuleName string,
	srcs []string,
	swiftImports []string,
	shouldSetVis bool,
	buildFile *rule.File,
) []*rule.Rule {
	// Detect the type of rule that should be used to build the Swift sources.
	r := buildRuleForTestSrcs(buildFile, defaultModuleName)
	setCommonSwiftAttrs(r, defaultModuleName, srcs, swiftImports)
	return []*rule.Rule{r}
}

func setCommonSwiftAttrs(r *rule.Rule, moduleName string, srcs []string, swiftImports []string) {
	if moduleName != "" {
		r.SetAttr(ModuleNameAttrName, moduleName)
	}
	if srcs != nil {
		r.SetAttr("srcs", srcs)
	}
	r.SetPrivateAttr(config.GazelleImportsKey, swiftImports)
}

// Alias

func aliasRule(name, actual string) *rule.Rule {
	r := rule.NewRule(AliasRuleKind, name)
	r.SetAttr("actual", actual)
	return r
}

// Visibility

func shouldSetVisibility(args language.GenerateArgs) bool {
	// If the package has a default visibility set, do not set visibility
	if args.File != nil && args.File.HasDefaultVisibility() {
		return false
	}
	return true
}

func setVisibilityAttr(r *rule.Rule, shouldSetVis bool, visibility []string) {
	if !shouldSetVis || visibility == nil {
		return
	}
	r.SetAttr("visibility", visibility)
}
