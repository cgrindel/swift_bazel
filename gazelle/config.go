package gazelle

import (
	"flag"
	"fmt"
	"log"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/cgrindel/swift_bazel/gazelle/internal/swift"
	"github.com/cgrindel/swift_bazel/gazelle/internal/swiftbin"
	"github.com/cgrindel/swift_bazel/gazelle/internal/wspace"
)

const swiftConfigName = "swift"

type swiftConfig struct {
	moduleIndex  *swift.ModuleIndex
	swiftBinPath string
}

func setStringToPtr(ptr *string, val string) {
	*ptr = val
}

func getSwiftConfig(c *config.Config) *swiftConfig {
	return c.Exts[swiftConfigName].(*swiftConfig)
}

func (*swiftLang) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	// Initialize location for custom configuration
	sc := &swiftConfig{
		moduleIndex: swift.NewModuleIndex(),
	}
	c.Exts[swiftConfigName] = sc
}

func (sl *swiftLang) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	sc := getSwiftConfig(c)

	// TODO(chuck): Add flag so that the client can tell use which Swift to use.

	// Find the Swift executable
	swiftBinPath, err := swiftbin.FindSwiftBinPath()
	if err != nil {
		return err
	}
	setStringToPtr(&sc.swiftBinPath, swiftBinPath)
	// DEBUG BEGIN
	log.Printf("*** CHUCK: CheckFlags swiftBinPath: %+#v", swiftBinPath)
	// DEBUG END

	// Look for http_archive declarations with Swift declarations.
	wkspFilePath := wspace.FindWORKSPACEFile(c.RepoRoot)
	wkspFile, err := rule.LoadWorkspaceFile(wkspFilePath, "")
	if err != nil {
		return fmt.Errorf("failed to load WORKSPACE file %v: %w", wkspFilePath, err)
	}
	archives, err := swift.NewHTTPArchivesFromWkspFile(wkspFile)
	if err != nil {
		return fmt.Errorf(
			"failed to retrieve archives from workspace file %v: %w",
			wkspFilePath,
			err,
		)
	}
	for _, archive := range archives {
		sc.moduleIndex.AddModules(archive.Modules...)
	}

	return nil
}
