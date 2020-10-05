package coding

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

//goland:noinspection SpellCheckingInspection
var multipartList = [][]string{
	{ // 7 Bits
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Praesent molestie eros ut ex dapibus sollicitudin in ut eros. Pellentesque venenatis vitae est e",
		"u porttitor. Nullam facilisis euismod felis, consectetur vehicula lorem aliquet at. Sed sit amet auctor lorem. Pellentesque euismod, orci non iaculis ull",
		"amcorper, massa ligula commodo sem, ac dictum lorem nulla vel tortor.",
	},
	{ // 1 byte
		"Лорем ипсум долор сит амет, еа яуи ностер елигенди. Перпетуа ассентиор ех нам. Ан молестие торяуатос вис, яуи виси тота трацтатос те. ",
		"Алии дебитис ин усу. Алии ерат тимеам дуо цу, пурто ерос иус ид, ет малис инвенире иус. Агам солет семпер яуо цу, граецо аперири витуп",
		"ерата еа цум. Ат сед дебет вениам сигниферумяуе, пер но миним фацете интеллегат, волумус демоцритум либерависсе еам цу.",
	},
	{ // 2 byte - Shift JIS
		"田ルマエ不効最ミラ報重マウタイ政身ロ朝連ドスど施康ルフオム象62訓誘30価ぽイ問需ニネキヲ指能トょそら問界亭ユムコス月投う続読ッ夕催ぶぼよ",
		"ょ京品状票変でッゅ。中執ぐでむ禁紀購ナアヒ中汁ノヘミオ行活リ和金ヒサテカ大情りイすち気康研ン瞬因ノ誕全げれフ働明えんざい最急僧かあ。入",
		"聴ほらゆで鉄42首ロ力録フタ都題月ロヨ真度づ見光ヨキ追5章否ご目3聞ぴぽ氏光ぽど価南べり回挑め。",
	},
	{ // 2 byte - EUC-KR
		"국회는 회계연도 개시 30일전까지 이를 의결하여야 한다, 동일한 범죄에 대하여 거듭 처벌받지 아니한다, 감사위원은 원장의 제청으로 대통령이",
		" 임명하고, 모든 국민은 신속한 재판을 받을 권리를 가진다.",
	},
	{ // UTF-16
		"👋🌌🎁🍛👵🐉🌁🐅🌜👰💧🕑👣🎤🌶🌐👺🍯🌶🌍 🎈🔶🐹👎🐁🔬🍌🔥💒🔴🌄🌁🕙",
		"📑🎒👤💫🔓🍜📎🎢🏭🐗🔊📶🌱🍣👞👫🔭📻📙🏢🍕🐻👼📍🎴🔪🌎🏫📭🍦🎃🐀🏮",
		"🕓🔐🏪🐈💐💧👷🏧🐸🍇💦🔧🔜👃📋💃🐒🔽🌃🔪🎲📇🍸",
	},
}

func TestSplit(t *testing.T) {
	limit := 134
	for _, multipart := range multipartList {
		expected := strings.Join(multipart, "")
		coding := BestCoding(expected)
		splitter := coding.Splitter()
		segments := splitter.Split(expected, limit)
		require.Equal(t, multipart, segments)
		encoder := coding.Encoding().NewEncoder()
		for _, segment := range segments {
			require.LessOrEqual(t, splitter.Len(segment), limit)
			encoded, err := encoder.Bytes([]byte(segment))
			require.NoError(t, err)
			require.Equal(t, splitter.Len(segment), len(encoded), segment)
		}
	}
}

func TestDataCoding(t *testing.T) {
	require.Nil(t, DataCoding(0b11111111).Encoding())
	require.Equal(t, DataCoding(0b11111111).GoString(), "11111111")
	require.NotEmpty(t, UCS2Coding.GoString())
	require.NotNil(t, UCS2Coding.Encoding())
}

//goland:noinspection SpellCheckingInspection
func TestBestCoding(t *testing.T) {
	mapping := map[DataCoding]string{
		GSM7BitCoding:  "ΨΠΦ",
		ASCIICoding:    "\x00abc",
		Latin1Coding:   "\u0100",
		ShiftJISCoding: "日本に行きたい。",
		CyrillicCoding: "\u0410",
		HebrewCoding:   "\u05B0",
		EUCKRCoding:    "안녕",
		UCS2Coding:     "💊",
	}
	for coding, input := range mapping {
		require.Equal(t, coding.String(), BestCoding(input).String())
	}
}

func TestBestSplitter(t *testing.T) {
	require.Nil(t, DataCoding(0b11111111).Splitter())
}
