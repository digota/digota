package braintree

import (
	"encoding/xml"
	"testing"
)

func TestModificationsXMLEmpty(t *testing.T) {
	t.Parallel()

	m := ModificationsRequest{}
	output, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatalf("got error %#v", err)
	}
	expectedOutput := `<ModificationsRequest></ModificationsRequest>`
	if string(output) != expectedOutput {
		t.Fatalf("got xml %#v, want %#v", string(output), expectedOutput)
	}
}

func TestModificationsXMLMinimalFields(t *testing.T) {
	t.Parallel()

	m := ModificationsRequest{
		Add: []AddModificationRequest{
			{
				InheritedFromID: "1",
			},
			{
				InheritedFromID: "2",
			},
		},
		Update: []UpdateModificationRequest{
			{
				ExistingID: "3",
			},
			{
				ExistingID: "4",
			},
		},
		RemoveExistingIDs: []string{
			"5",
			"6",
		},
	}
	output, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatalf("got error %#v", err)
	}
	expectedOutput := `<ModificationsRequest>
  <add type="array">
    <modification>
      <inherited-from-id>1</inherited-from-id>
    </modification>
    <modification>
      <inherited-from-id>2</inherited-from-id>
    </modification>
  </add>
  <update type="array">
    <modification>
      <existing-id>3</existing-id>
    </modification>
    <modification>
      <existing-id>4</existing-id>
    </modification>
  </update>
  <remove type="array">
    <modification>5</modification>
    <modification>6</modification>
  </remove>
</ModificationsRequest>`
	if string(output) != expectedOutput {
		t.Fatalf("got xml %#v, want %#v", string(output), expectedOutput)
	}
}

func TestModificationsXMLAllFields(t *testing.T) {
	t.Parallel()

	m := ModificationsRequest{
		Add: []AddModificationRequest{
			{
				InheritedFromID: "1",
				ModificationRequest: ModificationRequest{
					Amount:                NewDecimal(100, 2),
					NumberOfBillingCycles: 1,
					Quantity:              1,
					NeverExpires:          true,
				},
			},
			{
				InheritedFromID: "2",
				ModificationRequest: ModificationRequest{
					Amount:                NewDecimal(200, 2),
					NumberOfBillingCycles: 2,
					Quantity:              2,
					NeverExpires:          true,
				},
			},
		},
		Update: []UpdateModificationRequest{
			{
				ExistingID: "3",
				ModificationRequest: ModificationRequest{
					Amount:                NewDecimal(300, 2),
					NumberOfBillingCycles: 3,
					Quantity:              3,
					NeverExpires:          true,
				},
			},
			{
				ExistingID: "4",
				ModificationRequest: ModificationRequest{
					Amount:                NewDecimal(400, 2),
					NumberOfBillingCycles: 4,
					Quantity:              4,
					NeverExpires:          true,
				},
			},
		},
		RemoveExistingIDs: []string{
			"5",
			"6",
		},
	}
	output, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		t.Fatalf("got error %#v", err)
	}
	expectedOutput := `<ModificationsRequest>
  <add type="array">
    <modification>
      <amount>1.00</amount>
      <number-of-billing-cycles>1</number-of-billing-cycles>
      <quantity>1</quantity>
      <never-expires>true</never-expires>
      <inherited-from-id>1</inherited-from-id>
    </modification>
    <modification>
      <amount>2.00</amount>
      <number-of-billing-cycles>2</number-of-billing-cycles>
      <quantity>2</quantity>
      <never-expires>true</never-expires>
      <inherited-from-id>2</inherited-from-id>
    </modification>
  </add>
  <update type="array">
    <modification>
      <amount>3.00</amount>
      <number-of-billing-cycles>3</number-of-billing-cycles>
      <quantity>3</quantity>
      <never-expires>true</never-expires>
      <existing-id>3</existing-id>
    </modification>
    <modification>
      <amount>4.00</amount>
      <number-of-billing-cycles>4</number-of-billing-cycles>
      <quantity>4</quantity>
      <never-expires>true</never-expires>
      <existing-id>4</existing-id>
    </modification>
  </update>
  <remove type="array">
    <modification>5</modification>
    <modification>6</modification>
  </remove>
</ModificationsRequest>`
	if string(output) != expectedOutput {
		t.Fatalf("got xml %#v, want %#v", string(output), expectedOutput)
	}
}
