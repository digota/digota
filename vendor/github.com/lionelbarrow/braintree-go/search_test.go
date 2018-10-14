package braintree

import (
	"encoding/xml"
	"testing"
	"time"
)

func TestSearchXMLEncode(t *testing.T) {
	t.Parallel()

	s := new(SearchQuery)

	f := s.AddTextField("customer-first-name")
	f.Is = "A"
	f.IsNot = "B"
	f.StartsWith = "C"
	f.EndsWith = "D"
	f.Contains = "E"

	f2 := s.AddRangeField("amount")
	f2.Is = 15.01
	f2.Min = 10.01
	f2.Max = 20.01

	startDate := time.Date(2016, time.September, 11, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2016, time.September, 11, 23, 59, 59, 0, time.UTC)
	f3 := s.AddTimeField("settled-at")
	f3.Min = startDate
	f3.Max = endDate

	f4 := s.AddTimeField("created-at")
	f4.Min = startDate

	f5 := s.AddTimeField("authorization-expired-at")
	f5.Min = startDate

	f6 := s.AddTimeField("authorized-at")
	f6.Min = startDate

	f7 := s.AddTimeField("failed-at")
	f7.Min = startDate

	f8 := s.AddTimeField("gateway-rejected-at")
	f8.Min = startDate

	f9 := s.AddTimeField("processor-declined-at")
	f9.Min = startDate

	f10 := s.AddTimeField("submitted-for-settlement-at")
	f10.Min = startDate

	f11 := s.AddTimeField("voided-at")
	f11.Min = startDate

	f12 := s.AddTimeField("disbursement-date")
	f12.Min = startDate

	f13 := s.AddTimeField("dispute-date")
	f13.Min = startDate

	f14 := s.AddMultiField("status")
	f14.Items = []string{
		string(TransactionStatusAuthorized),
		string(TransactionStatusSubmittedForSettlement),
		string(TransactionStatusSettled),
	}

	b, err := xml.MarshalIndent(s, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	xmls := string(b)

	expect := `<search>
  <customer-first-name>
    <is>A</is>
    <is-not>B</is-not>
    <starts-with>C</starts-with>
    <ends-with>D</ends-with>
    <contains>E</contains>
  </customer-first-name>
  <amount>
    <is>15.01</is>
    <min>10.01</min>
    <max>20.01</max>
  </amount>
  <settled-at>
    <min type="datetime">2016-09-11T00:00:00Z</min>
    <max type="datetime">2016-09-11T23:59:59Z</max>
  </settled-at>
  <created-at>
    <min type="datetime">2016-09-11T00:00:00Z</min>
  </created-at>
  <authorization-expired-at>
    <min type="datetime">2016-09-11T00:00:00Z</min>
  </authorization-expired-at>
  <authorized-at>
    <min type="datetime">2016-09-11T00:00:00Z</min>
  </authorized-at>
  <failed-at>
    <min type="datetime">2016-09-11T00:00:00Z</min>
  </failed-at>
  <gateway-rejected-at>
    <min type="datetime">2016-09-11T00:00:00Z</min>
  </gateway-rejected-at>
  <processor-declined-at>
    <min type="datetime">2016-09-11T00:00:00Z</min>
  </processor-declined-at>
  <submitted-for-settlement-at>
    <min type="datetime">2016-09-11T00:00:00Z</min>
  </submitted-for-settlement-at>
  <voided-at>
    <min type="datetime">2016-09-11T00:00:00Z</min>
  </voided-at>
  <disbursement-date>
    <min type="datetime">2016-09-11T00:00:00Z</min>
  </disbursement-date>
  <dispute-date>
    <min type="datetime">2016-09-11T00:00:00Z</min>
  </dispute-date>
  <status type="array">
    <item>authorized</item>
    <item>submitted_for_settlement</item>
    <item>settled</item>
  </status>
</search>`

	if xmls != expect {
		t.Fatalf("got %#v, want %#v", xmls, expect)
	}
}

func TestSearchResultUnmarshal(t *testing.T) {
	t.Parallel()

	xmls := `<search-results>
  <page-size type="integer">50</page-size>
  <ids type="array">
      <item>k658ww</item>
      <item>fd2h96</item>
  </ids>
</search-results>`

	var v SearchResults
	err := xml.Unmarshal([]byte(xmls), &v)
	if err != nil {
		t.Fatal(err)
	}

	if len(v.Ids.Item) != 2 {
		t.Fatal(v.Ids)
	}
	if x := v.Ids.Item[0]; x != "k658ww" {
		t.Fatal(x)
	}
	if x := v.Ids.Item[1]; x != "fd2h96" {
		t.Fatal(x)
	}
}
