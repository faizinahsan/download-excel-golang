package domain

type UpdateProfile struct {
	Id                  int32
	Name                string
	Email               string
	OldDataMobile       string
	NewDataMobile       string
	Created             string
	ListDataChanges     string
	ReferenceNumber     string
	OldPhone            string `json:"old_phone"`
	NewPhone            string `json:"new_phone"`
	OldEmail            string `json:"old_email"`
	NewEmail            string `json:"new_email"`
	OldOccupation       string `json:"old_occupation"`
	NewOccupation       string `json:"new_occupation"`
	OldJob              string `json:"old_job"`
	NewJob              string `json:"new_job"`
	OldCompanyName      string `json:"old_companyname"`
	NewCompanyName      string `json:"new_companyname"`
	OldCompanyAddr      string `json:"old_companyaddr"`
	NewCompanyAddr      string `json:"new_companyaddr"`
	OldCompanyDesc      string `json:"old_companydesc"`
	NewCompanyDesc      string `json:"new_companydesc"`
	OldIncomeDesc       string `json:"old_incomedesc"`
	NewIncomeDesc       string `json:"new_incomedesc"`
	OldEmergencyPhone   string `json:"old_emergencyphone"`
	NewEmergencyPhone   string `json:"new_emergencyphone"`
	OldEmergencyContact string `json:"old_emergencycontact"`
	NewEmergencyContact string `json:"new_emergencycontact"`
	OldStreet           string `json:"old_street"`
	NewStreet           string `json:"new_street"`
	OldKecamatan        string `json:"old_kecamatan"`
	NewKecamatan        string `json:"new_kecamatan"`
	OldKelurahan        string `json:"old_kelurahan"`
	NewKelurahan        string `json:"new_kelurahan"`
	OldRT               string `json:"old_rt"`
	NewRT               string `json:"new_rt"`
	OldRW               string `json:"old_rw"`
	NewRW               string `json:"new_rw"`
	OldCityName         string `json:"old_cityName"`
	NewCityName         string `json:"new_cityName"`
	OldProvinceName     string `json:"old_provinceName"`
	NewProvinceName     string `json:"new_provinceName"`
	OldPostalCode       string `json:"old_postalcode"`
	NewPostalCode       string `json:"new_postalcode"`
}
