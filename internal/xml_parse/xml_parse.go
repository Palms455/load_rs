package xml_parse

import (
	"bytes"
	"encoding/xml"
	"golang.org/x/net/html/charset"
)

type Zglv struct {
	VERSION, DATA, FILENAME, SD_Z string `xml:",omitempty" json:",omitempty"`
}

type Zglv_p struct {
	VERSION, DATA, FILENAME, FILENAME1 string `xml:",omitempty" json:",omitempty"`
}

type Schet struct {
	CODE, CODE_MO, YEAR, MONTH, NSCHET, DSCHET, PLAT, SUMMAV, COMENTS, SUMMAP, SANK_MEK, SANK_MEE, SANK_EKMP, DISP string `xml:",omitempty" json:",omitempty"`

}


type Pacient struct {
	ID_PAC string
	VPOLIS string
	SPOLIS string `xml:",omitempty" json:",omitempty"`
	NPOLIS string `xml:",omitempty" json:",omitempty"`
	ST_OKATO string `xml:",omitempty" json:",omitempty"`
	SMO string `xml:",omitempty" json:",omitempty"`
	SMO_OGRN string `xml:",omitempty" json:",omitempty"`
	SMO_OK string `xml:",omitempty" json:",omitempty"`
	SMO_NAM string `xml:",omitempty" json:",omitempty"`
	INV string `xml:",omitempty" json:",omitempty"`
	MSE string `xml:",omitempty" json:",omitempty"`
	NOVOR string `xml:",omitempty" json:",omitempty"`
	VNOV_S string `xml:",omitempty" json:",omitempty"`
}



type Rs struct {
	ZGLV Zglv
	SCHET Schet
	ZAP []struct {
		N_ZAP, PR_NOV string
		PACIENT Pacient
		Z_SL []struct {
			IDCASE, USL_OK, VIDPOM, FOR_POM, NPR_MO, NPR_DATE, LPU, VBR, DATE_Z_1, DATE_Z_2, KD_Z string `xml:",omitempty" json:",omitempty"`
			VNOV_M []string `xml:",omitempty" json:",omitempty"`
			RSLT, ISHOD, P_OTK, RSLT_D string `xml:",omitempty" json:",omitempty"`
			OS_SLUCH []string `xml:",omitempty" json:",omitempty"`
			VB_P string `xml:",omitempty" json:",omitempty"`
			SL []struct {
				SL_ID string
				VOD_HMP, METOD_HMP, LPU_1, PODR, PROFIL, PROFIL_K, DET, TAL_D, TAL_NUM, TAL_P, P_CEL, NHISTORY, P_PER, DATE_1, DATE_2, KD, DS0, DS1 string `xml:",omitempty" json:",omitempty"`
				DS2, DS3 []string `xml:",omitempty" json:",omitempty"`
				C_ZAB, DS1_PR, DS_ONK, PR_D_N, DN string `xml:",omitempty" json:",omitempty"`
				CODE_MES1 []string `xml:",omitempty" json:",omitempty"`
				CODE_MES2 string `xml:",omitempty" json:",omitempty"`
				DS2_N []struct{
					DS2, DS2_PR, PR_DS2_N string `xml:",omitempty" json:",omitempty"`
				} `xml:",omitempty" json:",omitempty"`
				NAZ []struct {
					NAZ_N, NAZ_R, NAZ_SP, NAZ_V, NAZ_USL, NAPR_DATE, NAPR_MO, NAZ_PMP, NAZ_PK string `xml:",omitempty" json:",omitempty"`
				} `xml:",omitempty" json:",omitempty"`
				NAPR []struct{
					NAPR_DATE, NAPR_MO, NAPR_V, MET_ISSL, NAPR_USL string `xml:",omitempty" json:",omitempty"`
				} `xml:",omitempty" json:",omitempty"`
				CONS []struct{
					PR_CONS, DT_CONS string `xml:",omitempty" json:",omitempty"`
				} `xml:",omitempty" json:",omitempty"`
				ONK_SL []struct{
					DS1_T, STAD, ONK_T, ONK_N, ONK_M, MTSTZ, SOD, K_FR, WEI, HEI, BSA string `xml:",omitempty" json:",omitempty"`
					B_DIAG []struct{
						DIAG_DATE, DIAG_TIP, DIAG_CODE, DIAG_RSLT, REC_RSLT string `xml:",omitempty" json:",omitempty"`
					} `xml:",omitempty" json:",omitempty"`
					B_PROT []struct{
						PROT, D_PROT string `xml:",omitempty" json:",omitempty"`
					} `xml:",omitempty" json:",omitempty"`
					ONK_USL []struct{
						USL_TIP, HIR_TIP, LEK_TIP_L, LEK_TIP_V string `xml:",omitempty" json:",omitempty"`
						LEK_PR []struct{
							REGNUM, CODE_SH string `xml:",omitempty" json:",omitempty"`
							DATE_INJ []string
						} `xml:",omitempty" json:",omitempty"`
					} `xml:",omitempty" json:",omitempty"`
					PPTR, LUCH_TIP string `xml:",omitempty" json:",omitempty"`

				} `xml:",omitempty" json:",omitempty"`
				KSG_KPG []struct{
					N_KSG, VER_KSG, KSG_PG, N_KPG, KOEF_Z, KOEF_UP, BZTSZ, KOEF_D, KOEF_U string `xml:",omitempty" json:",omitempty"`
					CRIT []string
					SL_K, IT_SL string `xml:",omitempty" json:",omitempty"`
					SL_KOEF []struct{
						IDSL, Z_SL string `xml:",omitempty" json:",omitempty"`
					} `xml:",omitempty" json:",omitempty"`
				} `xml:",omitempty" json:",omitempty"`
				REAB, PRVS, VERS_SPEC, IDDOKT string `xml:",omitempty" json:",omitempty"`
				ED_COL string `xml:",omitempty" json:",omitempty"`
				TARIF string`xml:",omitempty" json:",omitempty"`
				SUM_M string
				USL []struct {
					IDSERV, LPU, LPU_1, PODR, PROFIL, VID_VME, DET, DATE_IN, DATE_OUT, DS, P_OTK, CODE_USL, KOL_USL string `xml:",omitempty" json:",omitempty"`
					TARIF, SUMV_USL, PRVS, CODE_MD, NPL, COMENTU string `xml:",omitempty" json:",omitempty"`
				} `xml:",omitempty" json:",omitempty"`
				COMENTSL string`xml:",omitempty" json:",omitempty"`
			}
			IDSP string `xml:",omitempty" json:",omitempty"`
			SUMV string `xml:",omitempty" json:",omitempty"`
			OPLATA string `xml:",omitempty" json:",omitempty"`
			SUMP string `xml:",omitempty" json:",omitempty"`
			SANK []struct{
				S_CODE string `xml:",omitempty" json:",omitempty"`
				S_SUM string `xml:",omitempty" json:",omitempty"`
				S_TIP string `xml:",omitempty" json:",omitempty"`
				SL_ID, S_OSN, DATE_ACT, NUM_ACT, CODE_EXP, S_COM, S_IST string `xml:",omitempty" json:",omitempty"`
			} `xml:",omitempty" json:",omitempty"`
			SANK_IT string `xml:",omitempty" json:",omitempty"`
		}
	}
}

type L struct {
	ZGLV struct{
		VERSION, DATA, FILENAME, FILENAME1 string
	}
	PERS []struct{
		ID_PAC string
		FAM, IM, OT string `xml:",omitempty" json:",omitempty"`
		W, DR string `xml:",omitempty" json:",omitempty"`
		DOST []string `xml:",omitempty" json:",omitempty"`
		TEL, FAM_P, IM_P, OT_P, W_P, DR_P, MR, DOCTYPE, DOCSER string `xml:",omitempty" json:",omitempty"`
		DOCNUM, DOCDATE, DOCORG, SNILS, OKATOG, OKATOP, COMENTP string `xml:",omitempty" json:",omitempty"`
	}
}


type RsData struct {
	Rs Rs
	LFile L
}


func (rs_data *RsData) XmlDecode(file *bytes.Buffer, ftype string)  {
	decoder := xml.NewDecoder(file)
	decoder.CharsetReader = charset.NewReaderLabel
	if ftype == "L" {
		_ = decoder.Decode(&rs_data.LFile)
	} else {
		_ = decoder.Decode(&rs_data.Rs)
	}
}
