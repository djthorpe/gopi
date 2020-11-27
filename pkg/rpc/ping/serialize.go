package ping

import (
	"strconv"
	"time"

	gopi "github.com/djthorpe/gopi/v3"
	ptypes "github.com/golang/protobuf/ptypes"
)

/////////////////////////////////////////////////////////////////////
// TYPES

type version struct {
	pb *VersionResponse
}

/////////////////////////////////////////////////////////////////////
// PUBLIC METHODS - VERSION

func toProtoVersion(v gopi.Version) *VersionResponse {
	if v == nil {
		return nil
	}
	tag, branch, hash := v.Version()
	ts, _ := ptypes.TimestampProto(v.BuildTime())
	return &VersionResponse{
		Name:      v.Name(),
		Tag:       tag,
		Branch:    branch,
		Hash:      hash,
		Buildtime: ts,
		Goversion: v.GoVersion(),
	}
}

func fromProtoVersion(pb *VersionResponse) gopi.Version {
	if pb == nil {
		return nil
	} else {
		return &version{pb}
	}
}

func (this *version) Name() string {
	return this.pb.Name
}

func (this *version) Version() (string, string, string) {
	return this.pb.Tag, this.pb.Branch, this.pb.Hash
}

func (this *version) GoVersion() string {
	return this.pb.Goversion
}

func (this *version) BuildTime() time.Time {
	if ts, err := ptypes.Timestamp(this.pb.Buildtime); err != nil {
		return time.Time{}
	} else {
		return ts
	}
}

///////////////////////////////////////////////////////////////////////////////
// STRINGIFY

func (this *version) String() string {
	str := "<version"
	str += " name=" + strconv.Quote(this.Name())
	tag, branch, hash := this.Version()
	if tag != "" {
		str += " tag=" + strconv.Quote(tag)
	}
	if branch != "" {
		str += " branch=" + strconv.Quote(branch)
	}
	if hash != "" {
		str += " hash=" + strconv.Quote(hash)
	}
	if goversion := this.GoVersion(); goversion != "" {
		str += " goversion=" + strconv.Quote(goversion)
	}
	if buildtime := this.BuildTime(); buildtime.IsZero() == false {
		str += " buildtime=" + strconv.Quote(buildtime.Format(time.RFC3339))
	}
	return str + ">"
}
