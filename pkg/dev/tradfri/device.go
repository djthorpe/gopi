package tradfri

import (
	"fmt"
	"strconv"
	"time"
)

////////////////////////////////////////////////////////////////////////////////
// TYPES

type device struct {
	Name_        string `json:"9001"`
	Created_     int64  `json:"9002"`
	Updated_     int64  `json:"9020"`
	Id_          uint   `json:"9003"`
	Active_      uint   `json:"9019"`
	Type_        uint   `json:"5750"`
	NeedsUpdate_ uint   `json:"9054"`

	Metadata_ struct {
		Vendor       string `json:"0"`
		Product      string `json:"1"`
		Serial       string `json:"2"`
		Version      string `json:"3"`
		PowerSource  int    `json:"6"`
		BatteryLevel int    `json:"9"`
	} `json:"3"`

	//Lights_ []lightbulb `json:"3311"`
}

////////////////////////////////////////////////////////////////////////////////
// IMPLEMENTATION

func (this *device) Name() string {
	return this.Name_
}

func (this *device) Id() uint {
	return this.Id_
}

func (this *device) Type() uint {
	return this.Type_
}

func (this *device) Created() time.Time {
	return time.Unix(this.Created_, 0)
}

func (this *device) Updated() time.Time {
	return time.Unix(this.Updated_, 0)
}

func (this *device) Active() bool {
	return this.Active_ != 0
}

func (this *device) Vendor() string {
	return this.Metadata_.Vendor
}

func (this *device) Product() string {
	return this.Metadata_.Product
}

func (this *device) Version() string {
	return this.Metadata_.Version
}

////////////////////////////////////////////////////////////////////////////////
// EQUALS

func (this *device) Equals(other *device) bool {
	if this.Name_ != other.Name_ {
		return false
	}
	if this.Created_ != other.Created_ {
		return false
	}
	if this.Updated_ != other.Updated_ {
		return false
	}
	if this.Id_ != other.Id_ {
		return false
	}
	if this.Active_ != other.Active_ {
		return false
	}
	if this.Type_ != other.Type_ {
		return false
	}
	if this.NeedsUpdate_ != other.NeedsUpdate_ {
		return false
	}
	if this.Metadata_.Vendor != other.Metadata_.Vendor {
		return false
	}
	if this.Metadata_.Product != other.Metadata_.Product {
		return false
	}
	if this.Metadata_.Serial != other.Metadata_.Serial {
		return false
	}
	if this.Metadata_.Version != other.Metadata_.Version {
		return false
	}
	if this.Metadata_.PowerSource != other.Metadata_.PowerSource {
		return false
	}
	if this.Metadata_.BatteryLevel != other.Metadata_.BatteryLevel {
		return false
	}
	/*
		if len(this.Lights_) != len(other.Lights_) {
			return false
		}
		for i, light := range this.Lights_ {
			if light.Equals(other.Lights_[i]) == false {
				return false
			}
		}
	*/
	// Otherwise, all equal
	return true
}

////////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *device) String() string {
	str := "<tradfri.device"
	str += " id=" + fmt.Sprint(this.Id())
	str += " type=" + fmt.Sprint(this.Type())
	str += " name=" + strconv.Quote(this.Name())
	str += " active=" + fmt.Sprint(this.Active())

	/*
		switch this.Type() {
		case mutablehome.IKEA_DEVICE_TYPE_LIGHT:
			str += " lights=" + fmt.Sprint(this.Lights())
		}

		if this.Metadata_.BatteryLevel != 0 {
			str += " batterylevel=" + fmt.Sprint(this.Metadata_.BatteryLevel)
		}

		if created := this.Created(); created.IsZero() == false {
			str += " created=" + created.Format(time.RFC3339)
		}
		if updated := this.Updated(); updated.IsZero() == false {
			str += " updated=" + updated.Format(time.RFC3339)
		}
		if vendor := this.Vendor(); vendor != "" {
			str += " vendor=" + strconv.Quote(vendor)
		}
		if product := this.Product(); product != "" {
			str += " product=" + strconv.Quote(product)
		}
		if version := this.Version(); version != "" {
			str += " version=" + strconv.Quote(version)
		}
	*/

	return str + ">"
}
