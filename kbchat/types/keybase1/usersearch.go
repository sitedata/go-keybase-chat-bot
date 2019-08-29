// Auto-generated types using avdl-compiler v1.4.1 (https://github.com/keybase/node-avdl-compiler)
//   Input file: ../client/protocol/avdl/keybase1/usersearch.avdl

package keybase1

import (
	"errors"
)

type APIUserServiceIDWithContact string

func (o APIUserServiceIDWithContact) DeepCopy() APIUserServiceIDWithContact {
	return o
}

type APIUserKeybaseResult struct {
	Username   string  `codec:"username" json:"username"`
	Uid        UID     `codec:"uid" json:"uid"`
	PictureUrl *string `codec:"pictureUrl,omitempty" json:"picture_url,omitempty"`
	FullName   *string `codec:"fullName,omitempty" json:"full_name,omitempty"`
	RawScore   float64 `codec:"rawScore" json:"raw_score"`
	Stellar    *string `codec:"stellar,omitempty" json:"stellar,omitempty"`
	IsFollowee bool    `codec:"isFollowee" json:"is_followee"`
}

func (o APIUserKeybaseResult) DeepCopy() APIUserKeybaseResult {
	return APIUserKeybaseResult{
		Username: o.Username,
		Uid:      o.Uid.DeepCopy(),
		PictureUrl: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.PictureUrl),
		FullName: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.FullName),
		RawScore: o.RawScore,
		Stellar: (func(x *string) *string {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.Stellar),
		IsFollowee: o.IsFollowee,
	}
}

type APIUserServiceResult struct {
	ServiceName APIUserServiceIDWithContact `codec:"serviceName" json:"service_name"`
	Username    string                      `codec:"username" json:"username"`
	PictureUrl  string                      `codec:"pictureUrl" json:"picture_url"`
	Bio         string                      `codec:"bio" json:"bio"`
	Location    string                      `codec:"location" json:"location"`
	FullName    string                      `codec:"fullName" json:"full_name"`
	Confirmed   *bool                       `codec:"confirmed,omitempty" json:"confirmed,omitempty"`
}

func (o APIUserServiceResult) DeepCopy() APIUserServiceResult {
	return APIUserServiceResult{
		ServiceName: o.ServiceName.DeepCopy(),
		Username:    o.Username,
		PictureUrl:  o.PictureUrl,
		Bio:         o.Bio,
		Location:    o.Location,
		FullName:    o.FullName,
		Confirmed: (func(x *bool) *bool {
			if x == nil {
				return nil
			}
			tmp := (*x)
			return &tmp
		})(o.Confirmed),
	}
}

type APIUserServiceSummary struct {
	ServiceName APIUserServiceIDWithContact `codec:"serviceName" json:"service_name"`
	Username    string                      `codec:"username" json:"username"`
}

func (o APIUserServiceSummary) DeepCopy() APIUserServiceSummary {
	return APIUserServiceSummary{
		ServiceName: o.ServiceName.DeepCopy(),
		Username:    o.Username,
	}
}

type ImpTofuSearchResult struct {
	Assertion       string `codec:"assertion" json:"assertion"`
	AssertionValue  string `codec:"assertionValue" json:"assertionValue"`
	AssertionKey    string `codec:"assertionKey" json:"assertionKey"`
	Label           string `codec:"label" json:"label"`
	PrettyName      string `codec:"prettyName" json:"prettyName"`
	KeybaseUsername string `codec:"keybaseUsername" json:"keybaseUsername"`
}

func (o ImpTofuSearchResult) DeepCopy() ImpTofuSearchResult {
	return ImpTofuSearchResult{
		Assertion:       o.Assertion,
		AssertionValue:  o.AssertionValue,
		AssertionKey:    o.AssertionKey,
		Label:           o.Label,
		PrettyName:      o.PrettyName,
		KeybaseUsername: o.KeybaseUsername,
	}
}

type APIUserSearchResult struct {
	Score           float64                                               `codec:"score" json:"score"`
	Keybase         *APIUserKeybaseResult                                 `codec:"keybase,omitempty" json:"keybase,omitempty"`
	Service         *APIUserServiceResult                                 `codec:"service,omitempty" json:"service,omitempty"`
	Contact         *ProcessedContact                                     `codec:"contact,omitempty" json:"contact,omitempty"`
	Imptofu         *ImpTofuSearchResult                                  `codec:"imptofu,omitempty" json:"imptofu,omitempty"`
	ServicesSummary map[APIUserServiceIDWithContact]APIUserServiceSummary `codec:"servicesSummary" json:"services_summary"`
	RawScore        float64                                               `codec:"rawScore" json:"rawScore"`
}

func (o APIUserSearchResult) DeepCopy() APIUserSearchResult {
	return APIUserSearchResult{
		Score: o.Score,
		Keybase: (func(x *APIUserKeybaseResult) *APIUserKeybaseResult {
			if x == nil {
				return nil
			}
			tmp := (*x).DeepCopy()
			return &tmp
		})(o.Keybase),
		Service: (func(x *APIUserServiceResult) *APIUserServiceResult {
			if x == nil {
				return nil
			}
			tmp := (*x).DeepCopy()
			return &tmp
		})(o.Service),
		Contact: (func(x *ProcessedContact) *ProcessedContact {
			if x == nil {
				return nil
			}
			tmp := (*x).DeepCopy()
			return &tmp
		})(o.Contact),
		Imptofu: (func(x *ImpTofuSearchResult) *ImpTofuSearchResult {
			if x == nil {
				return nil
			}
			tmp := (*x).DeepCopy()
			return &tmp
		})(o.Imptofu),
		ServicesSummary: (func(x map[APIUserServiceIDWithContact]APIUserServiceSummary) map[APIUserServiceIDWithContact]APIUserServiceSummary {
			if x == nil {
				return nil
			}
			ret := make(map[APIUserServiceIDWithContact]APIUserServiceSummary, len(x))
			for k, v := range x {
				kCopy := k.DeepCopy()
				vCopy := v.DeepCopy()
				ret[kCopy] = vCopy
			}
			return ret
		})(o.ServicesSummary),
		RawScore: o.RawScore,
	}
}

type NonUserDetails struct {
	IsNonUser      bool                  `codec:"isNonUser" json:"isNonUser"`
	AssertionValue string                `codec:"assertionValue" json:"assertionValue"`
	AssertionKey   string                `codec:"assertionKey" json:"assertionKey"`
	Description    string                `codec:"description" json:"description"`
	Contact        *ProcessedContact     `codec:"contact,omitempty" json:"contact,omitempty"`
	Service        *APIUserServiceResult `codec:"service,omitempty" json:"service,omitempty"`
	SiteIcon       []SizedImage          `codec:"siteIcon" json:"siteIcon"`
	SiteIconFull   []SizedImage          `codec:"siteIconFull" json:"siteIconFull"`
}

func (o NonUserDetails) DeepCopy() NonUserDetails {
	return NonUserDetails{
		IsNonUser:      o.IsNonUser,
		AssertionValue: o.AssertionValue,
		AssertionKey:   o.AssertionKey,
		Description:    o.Description,
		Contact: (func(x *ProcessedContact) *ProcessedContact {
			if x == nil {
				return nil
			}
			tmp := (*x).DeepCopy()
			return &tmp
		})(o.Contact),
		Service: (func(x *APIUserServiceResult) *APIUserServiceResult {
			if x == nil {
				return nil
			}
			tmp := (*x).DeepCopy()
			return &tmp
		})(o.Service),
		SiteIcon: (func(x []SizedImage) []SizedImage {
			if x == nil {
				return nil
			}
			ret := make([]SizedImage, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.SiteIcon),
		SiteIconFull: (func(x []SizedImage) []SizedImage {
			if x == nil {
				return nil
			}
			ret := make([]SizedImage, len(x))
			for i, v := range x {
				vCopy := v.DeepCopy()
				ret[i] = vCopy
			}
			return ret
		})(o.SiteIconFull),
	}
}

type ImpTofuSearchType int

const (
	ImpTofuSearchType_PHONE ImpTofuSearchType = 0
	ImpTofuSearchType_EMAIL ImpTofuSearchType = 1
)

func (o ImpTofuSearchType) DeepCopy() ImpTofuSearchType { return o }

var ImpTofuSearchTypeMap = map[string]ImpTofuSearchType{
	"PHONE": 0,
	"EMAIL": 1,
}

var ImpTofuSearchTypeRevMap = map[ImpTofuSearchType]string{
	0: "PHONE",
	1: "EMAIL",
}

func (e ImpTofuSearchType) String() string {
	if v, ok := ImpTofuSearchTypeRevMap[e]; ok {
		return v
	}
	return ""
}

type ImpTofuQuery struct {
	T__     ImpTofuSearchType `codec:"t" json:"t"`
	Phone__ *PhoneNumber      `codec:"phone,omitempty" json:"phone,omitempty"`
	Email__ *EmailAddress     `codec:"email,omitempty" json:"email,omitempty"`
}

func (o *ImpTofuQuery) T() (ret ImpTofuSearchType, err error) {
	switch o.T__ {
	case ImpTofuSearchType_PHONE:
		if o.Phone__ == nil {
			err = errors.New("unexpected nil value for Phone__")
			return ret, err
		}
	case ImpTofuSearchType_EMAIL:
		if o.Email__ == nil {
			err = errors.New("unexpected nil value for Email__")
			return ret, err
		}
	}
	return o.T__, nil
}

func (o ImpTofuQuery) Phone() (res PhoneNumber) {
	if o.T__ != ImpTofuSearchType_PHONE {
		panic("wrong case accessed")
	}
	if o.Phone__ == nil {
		return
	}
	return *o.Phone__
}

func (o ImpTofuQuery) Email() (res EmailAddress) {
	if o.T__ != ImpTofuSearchType_EMAIL {
		panic("wrong case accessed")
	}
	if o.Email__ == nil {
		return
	}
	return *o.Email__
}

func NewImpTofuQueryWithPhone(v PhoneNumber) ImpTofuQuery {
	return ImpTofuQuery{
		T__:     ImpTofuSearchType_PHONE,
		Phone__: &v,
	}
}

func NewImpTofuQueryWithEmail(v EmailAddress) ImpTofuQuery {
	return ImpTofuQuery{
		T__:     ImpTofuSearchType_EMAIL,
		Email__: &v,
	}
}

func (o ImpTofuQuery) DeepCopy() ImpTofuQuery {
	return ImpTofuQuery{
		T__: o.T__.DeepCopy(),
		Phone__: (func(x *PhoneNumber) *PhoneNumber {
			if x == nil {
				return nil
			}
			tmp := (*x).DeepCopy()
			return &tmp
		})(o.Phone__),
		Email__: (func(x *EmailAddress) *EmailAddress {
			if x == nil {
				return nil
			}
			tmp := (*x).DeepCopy()
			return &tmp
		})(o.Email__),
	}
}
