package sebaknode

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNodeInitState(t *testing.T) {
	require.Equal(t, NodeInitState, NodeStateNONE)
}

func TestNodeStateString(t *testing.T) {
	require.Equal(t, NodeStateNONE.String(), "NONE")
	require.Equal(t, NodeStateBOOTING.String(), "BOOTING")
	require.Equal(t, NodeStateCATCHUP.String(), "CATCHUP")
	require.Equal(t, NodeStateCONSENSUS.String(), "CONSENSUS")
	require.Equal(t, NodeStateTERMINATING.String(), "TERMINATING")
}

func TestNodeStateMarshalJSON(t *testing.T) {
	ret, err := NodeStateNONE.MarshalJSON()
	require.Equal(t, err, nil)
	require.Equal(t, "\"NONE\"", string(ret))

	ret, err = NodeStateBOOTING.MarshalJSON()
	require.Equal(t, err, nil)
	require.Equal(t, "\"BOOTING\"", string(ret))

	ret, err = NodeStateCATCHUP.MarshalJSON()
	require.Equal(t, err, nil)
	require.Equal(t, "\"CATCHUP\"", string(ret))

	ret, err = NodeStateCONSENSUS.MarshalJSON()
	require.Equal(t, err, nil)
	require.Equal(t, "\"CONSENSUS\"", string(ret))

	ret, err = NodeStateTERMINATING.MarshalJSON()
	require.Equal(t, err, nil)
	require.Equal(t, "\"TERMINATING\"", string(ret))
}

func TestNodeStateUnmarshalJSON(t *testing.T) {
	ns := NodeStateNONE
	nodeStateByteArray, _ := NodeStateNONE.MarshalJSON()
	ns.UnmarshalJSON(nodeStateByteArray)
	require.Equal(t, NodeStateNONE, ns)

	nodeStateByteArray, _ = NodeStateBOOTING.MarshalJSON()
	ns.UnmarshalJSON(nodeStateByteArray)
	require.Equal(t, NodeStateBOOTING, ns)

	nodeStateByteArray, _ = NodeStateCATCHUP.MarshalJSON()
	ns.UnmarshalJSON(nodeStateByteArray)
	require.Equal(t, NodeStateCATCHUP, ns)

	nodeStateByteArray, _ = NodeStateCONSENSUS.MarshalJSON()
	ns.UnmarshalJSON(nodeStateByteArray)
	require.Equal(t, NodeStateCONSENSUS, ns)

	nodeStateByteArray, _ = NodeStateTERMINATING.MarshalJSON()
	ns.UnmarshalJSON(nodeStateByteArray)
	require.Equal(t, NodeStateTERMINATING, ns)
}
