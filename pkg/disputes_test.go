package moov

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestDisputesMarshal(t *testing.T) {
	input := []byte(`{
			"amount": {
				"currency": "USD",
				"value": 1204
			},
			"createdOn": null,
			"disputeID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43",
			"networkReasonCode": null,
			"networkReasonDescription": null,
			"respondBy": null,
			"status": "response-needed",
			"transfer": {
				"transferID": "ec7e1848-dc80-4ab0-8827-dd7fc0737b43"
			}}`)

	dispute := new(Dispute)

	dec := json.NewDecoder(bytes.NewReader(input))
	dec.DisallowUnknownFields()

	err := dec.Decode(&dispute)
	if err != nil {
		require.NoError(t, err)
	}
	assert.Equal(t, "ec7e1848-dc80-4ab0-8827-dd7fc0737b43", dispute.DisputeID)
}

type DisputesTestSuite struct {
	suite.Suite
	// values for testing will be set in init()
	DisputeID string
}

// listen for 'go test' command --> run test methods
func TestDisputesSuite(t *testing.T) {
	suite.Run(t, new(DisputesTestSuite))
}

func (s *DisputesTestSuite) SetupSuite() {
	//mc, err := NewClient()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//disputes, err := mc.listdi()
	//if err != nil {
	//	log.Fatal(err)
	//}
}

func (s *DisputesTestSuite) TearDownSuite() {
}

func (s *DisputesTestSuite) TestListDisputes() {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	zeroTime := time.Time{}

	disputes, err := mc.ListDisputes(100, 0, zeroTime, zeroTime, "", "", "", zeroTime, zeroTime, "")
	if err != nil {
		s.T().Fatal(err)
	}
	fmt.Println(len(disputes))
	assert.NotNil(s.T(), disputes)

	if len(disputes) > 0 {
		s.DisputeID = disputes[0].DisputeID
	}
}

func (s *DisputesTestSuite) TestGetDispute() {
	mc, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	disputeID := s.DisputeID
	if disputeID == "" {
		disputeID = "2ce45e4e-8d96-45e4-8658-5767423e098d"
	}

	dispute, err := mc.GetDispute(disputeID)
	if err != nil {
		s.T().Fatal(err)
	}
	assert.Equal(s.T(), disputeID, dispute.DisputeID)
}
