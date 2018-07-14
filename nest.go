package conntrack

import (
	"errors"

	"github.com/mdlayher/netlink"
)

// Various errors which may occour when processing attributes
var (
	ErrAttrLength         = errors.New("Incorrect length of attribute")
	ErrAttrNotImplemented = errors.New("Attribute not implemented")
	ErrAttrNotExist       = errors.New("Type of attribute does not exist")
)

type attrCheckStruct struct {
	len, ct int
}

var attrCheck = map[ConnAttrType]attrCheckStruct{
	AttrOrigIPv4Src:             {len: 4, ct: ctaIPv4Src},
	AttrOrigIPv4Dst:             {len: 4, ct: ctaIPv4Dst},
	AttrReplIPv4Src:             {len: 4, ct: ctaIPv4Src},
	AttrReplIPv4Dst:             {len: 4, ct: ctaIPv4Dst},
	AttrOrigIPv6Src:             {len: 16, ct: ctaIPv6Src},
	AttrOrigIPv6Dst:             {len: 16, ct: ctaIPv6Dst},
	AttrReplIPv6Src:             {len: 16, ct: ctaIPv6Src},
	AttrReplIPv6Dst:             {len: 16, ct: ctaIPv6Dst},
	AttrOrigPortSrc:             {len: 2, ct: ctaUnspec},
	AttrOrigPortDst:             {len: 2, ct: ctaUnspec},
	AttrReplPortSrc:             {len: 2, ct: ctaUnspec},
	AttrReplPortDst:             {len: 2, ct: ctaUnspec},
	AttrIcmpType:                {len: 1, ct: ctaUnspec},
	AttrIcmpCode:                {len: 1, ct: ctaUnspec},
	AttrIcmpID:                  {len: 2, ct: ctaUnspec},
	AttrOrigL3Proto:             {len: 1, ct: ctaUnspec},
	AttrReplL3Proto:             {len: 1, ct: ctaUnspec},
	AttrOrigL4Proto:             {len: 1, ct: ctaProtoNum},
	AttrReplL4Proto:             {len: 1, ct: ctaProtoNum},
	AttrTCPState:                {len: 1, ct: ctaUnspec},
	AttrSNatIPv4:                {len: 4, ct: ctaUnspec},
	AttrDNatIPv4:                {len: 4, ct: ctaUnspec},
	AttrSNatPort:                {len: 2, ct: ctaUnspec},
	AttrDNatPort:                {len: 2, ct: ctaUnspec},
	AttrTimeout:                 {len: 4, ct: ctaTimeout},
	AttrMark:                    {len: 4, ct: ctaMark},
	AttrOrigCounterPackets:      {len: 8, ct: ctaUnspec},
	AttrReplCounterPackets:      {len: 8, ct: ctaUnspec},
	AttrOrigCounterBytes:        {len: 8, ct: ctaUnspec},
	AttrReplCounterBytes:        {len: 8, ct: ctaUnspec},
	AttrUse:                     {len: 4, ct: ctaUse},
	AttrID:                      {len: 4, ct: ctaID},
	AttrStatus:                  {len: 4, ct: ctaStatus},
	AttrTCPFlagsOrig:            {len: 1, ct: ctaUnspec},
	AttrTCPFlagsRepl:            {len: 1, ct: ctaUnspec},
	AttrTCPMaskOrig:             {len: 1, ct: ctaUnspec},
	AttrTCPMaskRepl:             {len: 1, ct: ctaUnspec},
	AttrMasterIPv4Src:           {len: 4, ct: ctaUnspec},
	AttrMasterIPv4Dst:           {len: 4, ct: ctaUnspec},
	AttrMasterIPv6Src:           {len: 16, ct: ctaUnspec},
	AttrMasterIPv6Dst:           {len: 16, ct: ctaUnspec},
	AttrMasterPortSrc:           {len: 2, ct: ctaUnspec},
	AttrMasterPortDst:           {len: 2, ct: ctaUnspec},
	AttrMasterL3Proto:           {len: 1, ct: ctaUnspec},
	AttrMasterL4Proto:           {len: 1, ct: ctaUnspec},
	AttrSecmark:                 {len: 4, ct: ctaSecmark},
	AttrOrigNatSeqCorrectionPos: {len: 4, ct: ctaUnspec},
	AttrOrigNatSeqOffsetBefore:  {len: 4, ct: ctaUnspec},
	AttrOrigNatSeqOffsetAfter:   {len: 4, ct: ctaUnspec},
	AttrReplNatSeqCorrectionPos: {len: 4, ct: ctaUnspec},
	AttrReplNatSeqOffsetBefore:  {len: 4, ct: ctaUnspec},
	AttrReplNatSeqOffsetAfter:   {len: 4, ct: ctaUnspec},
	AttrSctpState:               {len: 1, ct: ctaUnspec},
	AttrSctpVtagOrig:            {len: 4, ct: ctaUnspec},
	AttrSctpVtagRepl:            {len: 4, ct: ctaUnspec},
	AttrHelperName:              {len: 30, ct: ctaUnspec},
	AttrDccpState:               {len: 1, ct: ctaUnspec},
	AttrDccpRole:                {len: 1, ct: ctaUnspec},
	AttrDccpHandshakeSeq:        {len: 8, ct: ctaUnspec},
	AttrTCPWScaleOrig:           {len: 1, ct: ctaUnspec},
	AttrTCPWScaleRepl:           {len: 1, ct: ctaUnspec},
	AttrZone:                    {len: 2, ct: ctaZone},
	AttrSecCtx:                  {len: 30, ct: ctaUnspec},
	AttrTimestampStart:          {len: 8, ct: ctaUnspec},
	AttrTimestampStop:           {len: 8, ct: ctaUnspec},
	AttrHelperInfo:              {len: 30, ct: ctaUnspec},
	AttrConnlabels:              {len: 30, ct: ctaUnspec},
	AttrConnlabelsMask:          {len: 30, ct: ctaUnspec},
	AttrOrigzone:                {len: 2, ct: ctaUnspec},
	AttrReplzone:                {len: 2, ct: ctaUnspec},
	AttrSNatIPv6:                {len: 16, ct: ctaUnspec},
	AttrDNatIPv6:                {len: 16, ct: ctaUnspec},
}

func nestSubTuple(tupleType uint16, sub []netlink.Attribute) ([]byte, error) {
	attr, err := netlink.MarshalAttributes(sub)
	if err != nil {
		return nil, err
	}
	var tuple netlink.Attribute
	tuple.Type = tupleType | nlafNested
	tuple.Length = uint16(len(attr) + 4)
	tuple.Data = attr

	return tuple.MarshalBinary()
}

func nestDirTuple(dir int, ipTuple, protoTuple []netlink.Attribute) ([]byte, error) {
	ipSub, err := nestSubTuple(1, ipTuple)
	if err != nil {
		return nil, err
	}
	protoSub, err := nestSubTuple(2, protoTuple)
	if err != nil {
		return nil, err
	}
	var tuple netlink.Attribute
	if dir == 1 {
		tuple.Type = ctaTupleOrig | nlafNested
	} else {
		tuple.Type = ctaTupleReply | nlafNested
	}
	tuple.Length = uint16(len(ipSub) + len(protoSub) + 4)
	tuple.Data = ipSub
	tuple.Data = append(tuple.Data, protoSub...)

	return tuple.MarshalBinary()
}

func nestAttributes(filters []ConnAttr) ([]byte, error) {
	var attributes []byte
	var attrs []netlink.Attribute
	var tupleOrig, tupleRepl []netlink.Attribute
	var protoOrig, protoRepl []netlink.Attribute

	for _, filter := range filters {
		if _, ok := attrCheck[filter.Type]; !ok {
			return nil, ErrAttrNotExist
		}
		if attrCheck[filter.Type].ct == ctaUnspec {
			return nil, ErrAttrNotImplemented
		}
		if len(filter.Data) != attrCheck[filter.Type].len {
			return nil, ErrAttrLength
		}
		if filter.Type == AttrOrigIPv4Src ||
			filter.Type == AttrOrigIPv4Dst ||
			filter.Type == AttrOrigIPv6Src ||
			filter.Type == AttrOrigIPv6Dst {
			tupleOrig = append(tupleOrig, netlink.Attribute{Type: uint16(attrCheck[filter.Type].ct), Data: filter.Data})
		} else if filter.Type == AttrReplIPv4Src ||
			filter.Type == AttrReplIPv4Dst ||
			filter.Type == AttrReplIPv6Src ||
			filter.Type == AttrReplIPv6Dst {
			tupleRepl = append(tupleRepl, netlink.Attribute{Type: uint16(attrCheck[filter.Type].ct), Data: filter.Data})
		} else if filter.Type == AttrOrigL4Proto {
			protoOrig = append(protoOrig, netlink.Attribute{Type: uint16(attrCheck[filter.Type].ct), Data: filter.Data})
		} else if filter.Type == AttrReplL4Proto {
			protoRepl = append(protoRepl, netlink.Attribute{Type: uint16(attrCheck[filter.Type].ct), Data: filter.Data})
		} else {
			attrs = append(attrs, netlink.Attribute{Type: uint16(attrCheck[filter.Type].ct), Data: filter.Data})
		}
	}

	if len(tupleOrig) != 0 {
		data, err := nestDirTuple(1, tupleOrig, protoOrig)
		if err != nil {
			return nil, err
		}
		attributes = append(attributes, data...)
	}
	if len(tupleRepl) != 0 {
		data, err := nestDirTuple(0, tupleRepl, protoRepl)
		if err != nil {
			return nil, err
		}
		attributes = append(attributes, data...)
	}

	regular, err := netlink.MarshalAttributes(attrs)
	if err != nil {
		return nil, err
	}
	attributes = append(attributes, regular...)
	return attributes, nil
}