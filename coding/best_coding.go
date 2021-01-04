package coding

func BestCoding(input string) DataCoding {
	codings := []DataCoding{
		GSM7BitCoding, ASCIICoding, Latin1Coding,
		CyrillicCoding, HebrewCoding, ShiftJISCoding,
		EUCKRCoding,
	}
	for _, coding := range codings {
		if coding.Validate(input) {
			return coding
		}
	}
	return UCS2Coding
}

func BestSafeCoding(input string) DataCoding {
	if GSM7BitCoding.Validate(input) {
		return GSM7BitCoding
	}
	return UCS2Coding
}
