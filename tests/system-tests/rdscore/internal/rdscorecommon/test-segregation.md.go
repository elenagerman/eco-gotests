package rdscorecommon

import (
	"fmt"

	"github.com/golang/glog"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/internal/remote"
	. "github.com/openshift-kni/eco-gotests/tests/system-tests/rdscore/internal/rdscoreinittools"
	"github.com/openshift-kni/eco-gotests/tests/system-tests/rdscore/internal/rdscoreparams"
)

// TestFRRroute test MetalLB FRR segregationi.
func TestFRRroute(ctx SpecContext) {
	By("Asserting if test URL is provided")

	if RDSCoreConfig.MetalLBFRRTestURLIPv4 == "" && RDSCoreConfig.MetalLBFRRTestURLIPv6 == "" {
		glog.V(rdscoreparams.RDSCoreLogLevel).Infof(
			"Test URLs for MetalLB FRR testing not specified or are empty. Skipping...")
		Skip("Test URL for MetalLB FRR testing not specified or are empty")
	}

	execCmd := fmt.Sprintf("sudo podman container inspect --format \"{{.NetworkSettings.SandboxKey}}\" %s | xargs basename", RDSCoreConfig.MetallbFRRContainerName)
	lbOamNs, err := remote.ExecCmdOnHost(RDSCoreConfig.MetallbFRRHostName, RDSCoreConfig.MetallbFRRHostUser, RDSCoreConfig.MetallbFRRHostPass, execCmd)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to retrieve netns from the container %s; %v",
			RDSCoreConfig.MetallbFRRContainerName, err))

	execCmd = fmt.Sprintf("sudo ip netns exec %s curl -Lv %s", lbOamNs, RDSCoreConfig.MetalLBFRRTestURLIPv4)
	output, err := remote.ExecCmdOnHost(RDSCoreConfig.MetallbFRRHostName, RDSCoreConfig.MetallbFRRHostUser, RDSCoreConfig.MetallbFRRHostPass, execCmd)
	Expect(err).ToNot(HaveOccurred(),
		fmt.Sprintf("Failed to execute %s command due to: %v. \noutput: %v", execCmd, err, output))

	glog.V(rdscoreparams.RDSCoreLogLevel).Infof("Reseved response: %s", output)
}
