package delta

import (
	"crypto/md5"
	"strings"
	"testing"

	"github.com/OllieBM/eiger/pkg/mocks"
	"github.com/OllieBM/eiger/pkg/rolling_checksum"
	"github.com/OllieBM/eiger/pkg/signature"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
)

func TestCalculate(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	source := "Hello"
	target := "HelloWorld"

	rc := rolling_checksum.New()
	sig, err := signature.New(strings.NewReader(source), 5, md5.New(), rc)
	require.NoError(t, err)
	require.NotNil(t, sig)

	mockWriter := mocks.NewMockDiffWriter(ctrl)
	mockWriter.EXPECT().AddMatch(uint64(0)).Times(1)
	mockWriter.EXPECT().AddMiss(byte('W')).Times(1)
	mockWriter.EXPECT().AddMiss(byte('o')).Times(1)
	mockWriter.EXPECT().AddMiss(byte('r')).Times(1)
	mockWriter.EXPECT().AddMiss(byte('l')).Times(1)
	mockWriter.EXPECT().AddMiss(byte('d')).Times(1)
	mockWriter.EXPECT().AddMissingIndexes([]uint64{}).Times(1)
	err = Calculate(strings.NewReader(target), sig, rc, mockWriter)
	require.NoError(t, err)

}
