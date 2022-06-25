package obj

import (
	"testing"
	"time"
)

func Test_ParseAlphaMessage(t *testing.T) {
	data := []string{
		`POCSAG512: Address:  780777  Function: 0  Alpha:   PUTNAM FIRE BE ADVISED ENGINE 278 IS BACK IN SERVICE.<EOT>`,
		`POCSAG512: Address:   56896  Function: 1 `,
		`POCSAG512: Address:  750777  Function: 0  Alpha:   STA75,STA576,WOOD1 33.88 PRI 1 | 80F | Weak | Diarrhea | Universal Precautions |  Sick Person * 166 WOODSTOCK RD, Woodstock *   CADY LN / COUNTY RD 6/17/2022 10:48:36<EOT><NUL>`,
		`POCSAG512: Address:  550511  Function: 0  Alpha:   GENERAL SICKNESS<LF>10 ROLLING HILL RD<LF>X: CLARK HILL RD <LF><EOT><NUL><NUL>`,
		`POCSAG512: Address: 1605376  Function: 3 `,
		`POCSAG512: Address: 1089908  Function: 0 `,
		`POCSAG512: Address:   70777  Function: 1  Alpha:   1 SILO CIR, Mansfield B204 Juniper Hill Village ALS 89/F DIZZY LOW BACK PAIN CARDIAC HX  Sta7,A507,WCMH-M33  33.44 Cross Street DEAD END / ALDER LN 6/17/2022 11:28:50 2022-00001031 (03070-<EOT><NUL><NUL>`,
		`POCSAG512: Address:  159961  Function: 1  Alpha:   1 SILO CIR, Mansfield B204 Juniper Hill Village ALS 89/F DIZZY LOW BACK PAIN CARDIAC HX  Sta7,A507,WCMH-M33  33.44 Cross Street DEAD END / ALDER LN 6/17/2022 11:28:50 2022-00001031 (030<EOT><NUL><NUL>`,
		`POCSAG512: Address:   40777  Function: 3  Alpha:   579 N WINDHAM RD, North Windham  Berkshire Bank Training Center Fire Alarm-Commercial GENERAL FIRE ALARM Sta2,Sta3,Sta4  Cross Street BOSTON POST RD, N WINDHAM RD EXT / BEAVER HILL RD 6/17<EOT><NUL><NUL>`,
		`POCSAG512: Address:  541511  Function: 1  Alpha:   435 HARTFORD TPKE, Vernon INTEGRATED REHAB PARKING LOT   BLS 91/F VERTIGO VOMITING Sta641,A541  33.48 Cross Street MERLINE RD / EXIT 65 6/17/2022 11:32:49 2022-00002823 (03120-41)<EOT>`,
	}

	for _, d := range data {
		ts := time.Now()
		alpha, err := ParseAlphaMessage(ts, d)
		if err != nil {
			t.Fatal(err)
		}
		if alpha.Valid {
			t.Logf("Alpha: %#v", alpha)
		}
		time.Sleep(time.Second)
	}
}
