package enclave

import (
	enclaveTrustRequirement "github.com/enclave-networks/go-enclaveapi/data/trustrequirement"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func toTrustRequirementSlice(ids []types.Int64) []enclaveTrustRequirement.TrustRequirementId {
	output := make([]enclaveTrustRequirement.TrustRequirementId, len(ids))
	for i, id := range ids {
		output[i] = enclaveTrustRequirement.TrustRequirementId(id.Value)
	}

	return output
}
