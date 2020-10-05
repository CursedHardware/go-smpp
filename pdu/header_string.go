package pdu

import (
	"fmt"
	"strings"
)

//goland:noinspection SpellCheckingInspection
var commandStatusNames = map[CommandStatus]string{
	0x000: "ok",
	0x001: "invmsglen",
	0x002: "invcmdlen",
	0x003: "invcmdid",
	0x004: "invbndsts",
	0x005: "alybnd",
	0x006: "invprtflg",
	0x007: "invregdlvflg",
	0x008: "syserr",
	0x00A: "invsrcadr",
	0x00B: "invdstadr",
	0x00C: "invmsgid",
	0x00D: "bindfail",
	0x00E: "invpaswd",
	0x00F: "invsysid",
	0x011: "cancelfail",
	0x013: "replacefail",
	0x014: "msgqful",
	0x015: "invsertyp",
	0x033: "invnumdests",
	0x034: "invdlname",
	0x040: "invdestflag",
	0x042: "invsubrep",
	0x043: "invesmclass",
	0x044: "cntsubdl",
	0x045: "submitfail",
	0x048: "invsrcton",
	0x049: "invsrcnpi",
	0x050: "invdstton",
	0x051: "invdstnpi",
	0x053: "invsystyp",
	0x054: "invrepflag",
	0x055: "invnummsgs",
	0x058: "throttled",
	0x061: "invsched",
	0x062: "invexpiry",
	0x063: "invdftmsgid",
	0x064: "x_t_appn",
	0x065: "x_p_appn",
	0x066: "x_r_appn",
	0x067: "queryfail",
	0x0C0: "invtlvstream",
	0x0C1: "tlvnotallwd",
	0x0C2: "invtlvlen",
	0x0C3: "missingtlv",
	0x0C4: "invtlvval",
	0x0FE: "deliveryfailure",
	0x0FF: "unknownerr",
	0x100: "sertypunauth",
	0x101: "prohibited",
	0x102: "sertypunavail",
	0x103: "sertypdenied",
	0x104: "invdcs",
	0x105: "invsrcaddrsubunit",
	0x106: "invdstaddrsubunit",
	0x107: "invbcastfreqint",
	0x108: "invbcastalias_name",
	0x109: "invbcastareafmt",
	0x10A: "invnumbcast_areas",
	0x10B: "invbcastcnttype",
	0x10C: "invbcastmsgclass",
	0x10D: "bcastfail",
	0x10E: "bcastqueryfail",
	0x10F: "bcastcancelfail",
	0x110: "invbcast_rep",
	0x111: "invbcastsrvgrp",
	0x112: "invbcastchanind",
}

func (c CommandStatus) String() string {
	if name, ok := commandStatusNames[c]; ok {
		return fmt.Sprintf("ESME_R%s", strings.ToUpper(name))
	}
	return fmt.Sprintf("%08X", uint32(c))
}
