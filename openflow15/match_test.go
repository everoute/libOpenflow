package openflow15

import (
	"bytes"
	"fmt"
	"net"
	"testing"
)

func TestMatchEthAddresses(t *testing.T) {
	ethSrcAddress, _ := net.ParseMAC("aa:aa:aa:aa:aa:aa")
	ethDstAddress, _ := net.ParseMAC("ff:ff:ff:ff:ff:ff")

	ofMatch := NewMatch()
	{
		macSrcField := NewEthSrcField(ethSrcAddress, nil)
		ofMatch.AddField(*macSrcField)

		macDstField := NewEthDstField(ethDstAddress, nil)
		ofMatch.AddField(*macDstField)
	}

	if err := checkMatchSerializationConsistency(ofMatch); err != nil {
		t.Fatal(err)
	}
}

func TestMatchIcmpFields(t *testing.T) {
	ofMatch := NewMatch()
	{
		icmpTypeField := NewIcmpTypeField(8)
		if icmpTypeField.Class != OXM_CLASS_OPENFLOW_BASIC || icmpTypeField.Field != OXM_FIELD_ICMPV4_TYPE {
			t.Fatalf("unexpected icmp type field header: class=%d field=%d", icmpTypeField.Class, icmpTypeField.Field)
		}
		ofMatch.AddField(*icmpTypeField)

		icmpCodeField := NewIcmpCodeField(0)
		if icmpCodeField.Class != OXM_CLASS_OPENFLOW_BASIC || icmpCodeField.Field != OXM_FIELD_ICMPV4_CODE {
			t.Fatalf("unexpected icmp code field header: class=%d field=%d", icmpCodeField.Class, icmpCodeField.Field)
		}
		ofMatch.AddField(*icmpCodeField)
	}

	if err := checkMatchSerializationConsistency(ofMatch); err != nil {
		t.Fatal(err)
	}
}

func checkMatchSerializationConsistency(ofMatch *Match) error {
	// Serialize the original match
	ofMatchRaw, err := ofMatch.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal match: %w", err)
	}

	// Deserialize into a new match object
	ofMatchRecovered := NewMatch()
	if err := ofMatchRecovered.UnmarshalBinary(ofMatchRaw); err != nil {
		return fmt.Errorf("failed to unmarshal match: %w", err)
	}

	// Serialize the recovered match for comparison
	ofMatchRecoveredRaw, err := ofMatchRecovered.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal recovered match: %w", err)
	}

	// Check for serialization consistency
	if !bytes.Equal(ofMatchRaw, ofMatchRecoveredRaw) {
		return fmt.Errorf("initial and recovered match structures do not match")
	}

	return nil
}
