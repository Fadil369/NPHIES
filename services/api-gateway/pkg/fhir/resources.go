package fhir

import "time"

// Base FHIR Resource structure
type Resource struct {
	ResourceType string `json:"resourceType"`
	ID           string `json:"id,omitempty"`
	Meta         *Meta  `json:"meta,omitempty"`
}

type Meta struct {
	VersionID   string `json:"versionId,omitempty"`
	LastUpdated string `json:"lastUpdated,omitempty"`
	Profile     []string `json:"profile,omitempty"`
	Security    []Coding `json:"security,omitempty"`
	Tag         []Coding `json:"tag,omitempty"`
}

type Coding struct {
	System  string `json:"system,omitempty"`
	Version string `json:"version,omitempty"`
	Code    string `json:"code,omitempty"`
	Display string `json:"display,omitempty"`
}

type CodeableConcept struct {
	Coding []Coding `json:"coding,omitempty"`
	Text   string   `json:"text,omitempty"`
}

type Identifier struct {
	Use      string           `json:"use,omitempty"`
	Type     *CodeableConcept `json:"type,omitempty"`
	System   string           `json:"system,omitempty"`
	Value    string           `json:"value,omitempty"`
	Period   *Period          `json:"period,omitempty"`
	Assigner *Reference       `json:"assigner,omitempty"`
}

type Period struct {
	Start string `json:"start,omitempty"`
	End   string `json:"end,omitempty"`
}

type Reference struct {
	Reference string `json:"reference,omitempty"`
	Type      string `json:"type,omitempty"`
	Identifier *Identifier `json:"identifier,omitempty"`
	Display   string `json:"display,omitempty"`
}

type HumanName struct {
	Use    string   `json:"use,omitempty"`
	Text   string   `json:"text,omitempty"`
	Family string   `json:"family,omitempty"`
	Given  []string `json:"given,omitempty"`
	Prefix []string `json:"prefix,omitempty"`
	Suffix []string `json:"suffix,omitempty"`
	Period *Period  `json:"period,omitempty"`
}

type ContactPoint struct {
	System string  `json:"system,omitempty"`
	Value  string  `json:"value,omitempty"`
	Use    string  `json:"use,omitempty"`
	Rank   int     `json:"rank,omitempty"`
	Period *Period `json:"period,omitempty"`
}

type Address struct {
	Use        string   `json:"use,omitempty"`
	Type       string   `json:"type,omitempty"`
	Text       string   `json:"text,omitempty"`
	Line       []string `json:"line,omitempty"`
	City       string   `json:"city,omitempty"`
	District   string   `json:"district,omitempty"`
	State      string   `json:"state,omitempty"`
	PostalCode string   `json:"postalCode,omitempty"`
	Country    string   `json:"country,omitempty"`
	Period     *Period  `json:"period,omitempty"`
}

// FHIR Patient Resource
type Patient struct {
	ResourceType      string           `json:"resourceType"`
	ID                string           `json:"id,omitempty"`
	Meta              *Meta            `json:"meta,omitempty"`
	Identifier        []Identifier     `json:"identifier,omitempty"`
	Active            *bool            `json:"active,omitempty"`
	Name              []HumanName      `json:"name,omitempty"`
	Telecom           []ContactPoint   `json:"telecom,omitempty"`
	Gender            string           `json:"gender,omitempty"`
	BirthDate         string           `json:"birthDate,omitempty"`
	DeceasedBoolean   *bool            `json:"deceasedBoolean,omitempty"`
	DeceasedDateTime  string           `json:"deceasedDateTime,omitempty"`
	Address           []Address        `json:"address,omitempty"`
	MaritalStatus     *CodeableConcept `json:"maritalStatus,omitempty"`
	MultipleBirthBoolean *bool         `json:"multipleBirthBoolean,omitempty"`
	MultipleBirthInteger *int          `json:"multipleBirthInteger,omitempty"`
	Photo             []Attachment     `json:"photo,omitempty"`
	Contact           []PatientContact `json:"contact,omitempty"`
	Communication     []PatientCommunication `json:"communication,omitempty"`
	GeneralPractitioner []Reference   `json:"generalPractitioner,omitempty"`
	ManagingOrganization *Reference    `json:"managingOrganization,omitempty"`
	Link              []PatientLink    `json:"link,omitempty"`
}

type PatientContact struct {
	Relationship []CodeableConcept `json:"relationship,omitempty"`
	Name         *HumanName        `json:"name,omitempty"`
	Telecom      []ContactPoint    `json:"telecom,omitempty"`
	Address      *Address          `json:"address,omitempty"`
	Gender       string            `json:"gender,omitempty"`
	Organization *Reference        `json:"organization,omitempty"`
	Period       *Period           `json:"period,omitempty"`
}

type PatientCommunication struct {
	Language  CodeableConcept `json:"language"`
	Preferred *bool           `json:"preferred,omitempty"`
}

type PatientLink struct {
	Other Reference `json:"other"`
	Type  string    `json:"type"`
}

type Attachment struct {
	ContentType string     `json:"contentType,omitempty"`
	Language    string     `json:"language,omitempty"`
	Data        string     `json:"data,omitempty"`
	URL         string     `json:"url,omitempty"`
	Size        int        `json:"size,omitempty"`
	Hash        string     `json:"hash,omitempty"`
	Title       string     `json:"title,omitempty"`
	Creation    string     `json:"creation,omitempty"`
}

// FHIR Coverage Resource
type Coverage struct {
	ResourceType     string           `json:"resourceType"`
	ID               string           `json:"id,omitempty"`
	Meta             *Meta            `json:"meta,omitempty"`
	Identifier       []Identifier     `json:"identifier,omitempty"`
	Status           string           `json:"status"`
	Type             *CodeableConcept `json:"type,omitempty"`
	PolicyHolder     *Reference       `json:"policyHolder,omitempty"`
	Subscriber       *Reference       `json:"subscriber,omitempty"`
	SubscriberID     string           `json:"subscriberId,omitempty"`
	Beneficiary      Reference        `json:"beneficiary"`
	Dependent        string           `json:"dependent,omitempty"`
	Relationship     *CodeableConcept `json:"relationship,omitempty"`
	Period           *Period          `json:"period,omitempty"`
	Payor            []Reference      `json:"payor"`
	Class            []CoverageClass  `json:"class,omitempty"`
	Order            int              `json:"order,omitempty"`
	Network          string           `json:"network,omitempty"`
	CostToBeneficiary []CoverageCostToBeneficiary `json:"costToBeneficiary,omitempty"`
	Subrogation      *bool            `json:"subrogation,omitempty"`
	Contract         []Reference      `json:"contract,omitempty"`
}

type CoverageClass struct {
	Type  CodeableConcept `json:"type"`
	Value string          `json:"value"`
	Name  string          `json:"name,omitempty"`
}

type CoverageCostToBeneficiary struct {
	Type      *CodeableConcept `json:"type,omitempty"`
	ValueQuantity *Quantity    `json:"valueQuantity,omitempty"`
	ValueMoney    *Money       `json:"valueMoney,omitempty"`
	Exception []CoverageCostToBeneficiaryException `json:"exception,omitempty"`
}

type CoverageCostToBeneficiaryException struct {
	Type   CodeableConcept `json:"type"`
	Period *Period         `json:"period,omitempty"`
}

type Quantity struct {
	Value      float64 `json:"value,omitempty"`
	Comparator string  `json:"comparator,omitempty"`
	Unit       string  `json:"unit,omitempty"`
	System     string  `json:"system,omitempty"`
	Code       string  `json:"code,omitempty"`
}

type Money struct {
	Value    float64 `json:"value,omitempty"`
	Currency string  `json:"currency,omitempty"`
}

// FHIR Claim Resource
type Claim struct {
	ResourceType          string                `json:"resourceType"`
	ID                    string                `json:"id,omitempty"`
	Meta                  *Meta                 `json:"meta,omitempty"`
	Identifier            []Identifier          `json:"identifier,omitempty"`
	Status                string                `json:"status"`
	Type                  CodeableConcept       `json:"type"`
	SubType               *CodeableConcept      `json:"subType,omitempty"`
	Use                   string                `json:"use"`
	Patient               Reference             `json:"patient"`
	BillablePeriod        *Period               `json:"billablePeriod,omitempty"`
	Created               string                `json:"created"`
	Enterer               *Reference            `json:"enterer,omitempty"`
	Insurer               *Reference            `json:"insurer,omitempty"`
	Provider              Reference             `json:"provider"`
	Priority              CodeableConcept       `json:"priority"`
	FundsReserve          *CodeableConcept      `json:"fundsReserve,omitempty"`
	Related               []ClaimRelated        `json:"related,omitempty"`
	Prescription          *Reference            `json:"prescription,omitempty"`
	OriginalPrescription  *Reference            `json:"originalPrescription,omitempty"`
	Payee                 *ClaimPayee           `json:"payee,omitempty"`
	Referral              *Reference            `json:"referral,omitempty"`
	Facility              *Reference            `json:"facility,omitempty"`
	CareTeam              []ClaimCareTeam       `json:"careTeam,omitempty"`
	SupportingInfo        []ClaimSupportingInfo `json:"supportingInfo,omitempty"`
	Diagnosis             []ClaimDiagnosis      `json:"diagnosis,omitempty"`
	Procedure             []ClaimProcedure      `json:"procedure,omitempty"`
	Insurance             []ClaimInsurance      `json:"insurance"`
	Accident              *ClaimAccident        `json:"accident,omitempty"`
	Item                  []ClaimItem           `json:"item,omitempty"`
	Total                 *Money                `json:"total,omitempty"`
}

type ClaimRelated struct {
	Claim       *Reference       `json:"claim,omitempty"`
	Relationship *CodeableConcept `json:"relationship,omitempty"`
	Reference   *Identifier      `json:"reference,omitempty"`
}

type ClaimPayee struct {
	Type  CodeableConcept `json:"type"`
	Party *Reference      `json:"party,omitempty"`
}

type ClaimCareTeam struct {
	Sequence     int              `json:"sequence"`
	Provider     Reference        `json:"provider"`
	Responsible  *bool            `json:"responsible,omitempty"`
	Role         *CodeableConcept `json:"role,omitempty"`
	Qualification *CodeableConcept `json:"qualification,omitempty"`
}

type ClaimSupportingInfo struct {
	Sequence         int              `json:"sequence"`
	Category         CodeableConcept  `json:"category"`
	Code             *CodeableConcept `json:"code,omitempty"`
	TimingDate       string           `json:"timingDate,omitempty"`
	TimingPeriod     *Period          `json:"timingPeriod,omitempty"`
	ValueBoolean     *bool            `json:"valueBoolean,omitempty"`
	ValueString      string           `json:"valueString,omitempty"`
	ValueQuantity    *Quantity        `json:"valueQuantity,omitempty"`
	ValueAttachment  *Attachment      `json:"valueAttachment,omitempty"`
	ValueReference   *Reference       `json:"valueReference,omitempty"`
	Reason           *CodeableConcept `json:"reason,omitempty"`
}

type ClaimDiagnosis struct {
	Sequence           int               `json:"sequence"`
	DiagnosisCodeableConcept *CodeableConcept `json:"diagnosisCodeableConcept,omitempty"`
	DiagnosisReference *Reference        `json:"diagnosisReference,omitempty"`
	Type               []CodeableConcept `json:"type,omitempty"`
	OnAdmission        *CodeableConcept  `json:"onAdmission,omitempty"`
	PackageCode        *CodeableConcept  `json:"packageCode,omitempty"`
}

type ClaimProcedure struct {
	Sequence           int               `json:"sequence"`
	Type               []CodeableConcept `json:"type,omitempty"`
	Date               string            `json:"date,omitempty"`
	ProcedureCodeableConcept *CodeableConcept `json:"procedureCodeableConcept,omitempty"`
	ProcedureReference *Reference        `json:"procedureReference,omitempty"`
	UDI                []Reference       `json:"udi,omitempty"`
}

type ClaimInsurance struct {
	Sequence       int        `json:"sequence"`
	Focal          bool       `json:"focal"`
	Identifier     *Identifier `json:"identifier,omitempty"`
	Coverage       Reference   `json:"coverage"`
	BusinessArrangement string `json:"businessArrangement,omitempty"`
	PreAuthRef     []string    `json:"preAuthRef,omitempty"`
	ClaimResponse  *Reference  `json:"claimResponse,omitempty"`
}

type ClaimAccident struct {
	Date            string           `json:"date"`
	Type            *CodeableConcept `json:"type,omitempty"`
	LocationAddress *Address         `json:"locationAddress,omitempty"`
	LocationReference *Reference     `json:"locationReference,omitempty"`
}

type ClaimItem struct {
	Sequence                int                     `json:"sequence"`
	CareTeamSequence        []int                   `json:"careTeamSequence,omitempty"`
	DiagnosisSequence       []int                   `json:"diagnosisSequence,omitempty"`
	ProcedureSequence       []int                   `json:"procedureSequence,omitempty"`
	InformationSequence     []int                   `json:"informationSequence,omitempty"`
	Revenue                 *CodeableConcept        `json:"revenue,omitempty"`
	Category                *CodeableConcept        `json:"category,omitempty"`
	ProductOrService        CodeableConcept         `json:"productOrService"`
	Modifier                []CodeableConcept       `json:"modifier,omitempty"`
	ProgramCode             []CodeableConcept       `json:"programCode,omitempty"`
	ServicedDate            string                  `json:"servicedDate,omitempty"`
	ServicedPeriod          *Period                 `json:"servicedPeriod,omitempty"`
	LocationCodeableConcept *CodeableConcept        `json:"locationCodeableConcept,omitempty"`
	LocationAddress         *Address                `json:"locationAddress,omitempty"`
	LocationReference       *Reference              `json:"locationReference,omitempty"`
	Quantity                *Quantity               `json:"quantity,omitempty"`
	UnitPrice               *Money                  `json:"unitPrice,omitempty"`
	Factor                  float64                 `json:"factor,omitempty"`
	Net                     *Money                  `json:"net,omitempty"`
	UDI                     []Reference             `json:"udi,omitempty"`
	BodySite                *CodeableConcept        `json:"bodySite,omitempty"`
	SubSite                 []CodeableConcept       `json:"subSite,omitempty"`
	Encounter               []Reference             `json:"encounter,omitempty"`
	Detail                  []ClaimItemDetail       `json:"detail,omitempty"`
}

type ClaimItemDetail struct {
	Sequence         int                     `json:"sequence"`
	Revenue          *CodeableConcept        `json:"revenue,omitempty"`
	Category         *CodeableConcept        `json:"category,omitempty"`
	ProductOrService CodeableConcept         `json:"productOrService"`
	Modifier         []CodeableConcept       `json:"modifier,omitempty"`
	ProgramCode      []CodeableConcept       `json:"programCode,omitempty"`
	Quantity         *Quantity               `json:"quantity,omitempty"`
	UnitPrice        *Money                  `json:"unitPrice,omitempty"`
	Factor           float64                 `json:"factor,omitempty"`
	Net              *Money                  `json:"net,omitempty"`
	UDI              []Reference             `json:"udi,omitempty"`
	SubDetail        []ClaimItemDetailSubDetail `json:"subDetail,omitempty"`
}

type ClaimItemDetailSubDetail struct {
	Sequence         int                     `json:"sequence"`
	Revenue          *CodeableConcept        `json:"revenue,omitempty"`
	Category         *CodeableConcept        `json:"category,omitempty"`
	ProductOrService CodeableConcept         `json:"productOrService"`
	Modifier         []CodeableConcept       `json:"modifier,omitempty"`
	ProgramCode      []CodeableConcept       `json:"programCode,omitempty"`
	Quantity         *Quantity               `json:"quantity,omitempty"`
	UnitPrice        *Money                  `json:"unitPrice,omitempty"`
	Factor           float64                 `json:"factor,omitempty"`
	Net              *Money                  `json:"net,omitempty"`
	UDI              []Reference             `json:"udi,omitempty"`
}

// FHIR Bundle for search results
type Bundle struct {
	ResourceType string       `json:"resourceType"`
	ID           string       `json:"id,omitempty"`
	Meta         *Meta        `json:"meta,omitempty"`
	Identifier   *Identifier  `json:"identifier,omitempty"`
	Type         string       `json:"type"`
	Timestamp    string       `json:"timestamp,omitempty"`
	Total        int          `json:"total,omitempty"`
	Link         []BundleLink `json:"link,omitempty"`
	Entry        []BundleEntry `json:"entry,omitempty"`
	Signature    *Signature   `json:"signature,omitempty"`
}

type BundleLink struct {
	Relation string `json:"relation"`
	URL      string `json:"url"`
}

type BundleEntry struct {
	Link     []BundleLink    `json:"link,omitempty"`
	FullURL  string          `json:"fullUrl,omitempty"`
	Resource interface{}     `json:"resource,omitempty"`
	Search   *BundleEntrySearch `json:"search,omitempty"`
	Request  *BundleEntryRequest `json:"request,omitempty"`
	Response *BundleEntryResponse `json:"response,omitempty"`
}

type BundleEntrySearch struct {
	Mode  string  `json:"mode,omitempty"`
	Score float64 `json:"score,omitempty"`
}

type BundleEntryRequest struct {
	Method          string `json:"method"`
	URL             string `json:"url"`
	IfNoneMatch     string `json:"ifNoneMatch,omitempty"`
	IfModifiedSince string `json:"ifModifiedSince,omitempty"`
	IfMatch         string `json:"ifMatch,omitempty"`
	IfNoneExist     string `json:"ifNoneExist,omitempty"`
}

type BundleEntryResponse struct {
	Status       string `json:"status"`
	Location     string `json:"location,omitempty"`
	Etag         string `json:"etag,omitempty"`
	LastModified string `json:"lastModified,omitempty"`
	Outcome      interface{} `json:"outcome,omitempty"`
}

type Signature struct {
	Type         []Coding    `json:"type"`
	When         string      `json:"when"`
	Who          Reference   `json:"who"`
	OnBehalfOf   *Reference  `json:"onBehalfOf,omitempty"`
	TargetFormat string      `json:"targetFormat,omitempty"`
	SigFormat    string      `json:"sigFormat,omitempty"`
	Data         string      `json:"data,omitempty"`
}