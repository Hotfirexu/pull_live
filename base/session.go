package base

const (
	SessionTypeCustomizePub      SessionType = SessionProtocolCustomize<<8 | SessionBaseTypePub
	SessionTypeRtmpServerSession SessionType = SessionProtocolRtmp<<8 | SessionBaseTypePubSub
	SessionTypeRtmpPush          SessionType = SessionProtocolRtmp<<8 | SessionBaseTypePush
	SessionTypeRtmpPull          SessionType = SessionProtocolRtmp<<8 | SessionBaseTypePull
	SessionTypeRtspPub           SessionType = SessionProtocolRtsp<<8 | SessionBaseTypePub
	SessionTypeRtspSub           SessionType = SessionProtocolRtsp<<8 | SessionBaseTypeSub
	SessionTypeRtspPush          SessionType = SessionProtocolRtsp<<8 | SessionBaseTypePush
	SessionTypeRtspPull          SessionType = SessionProtocolRtsp<<8 | SessionBaseTypePull
	SessionTypeFlvSub            SessionType = SessionProtocolFlv<<8 | SessionBaseTypeSub
	SessionTypeFlvPull           SessionType = SessionProtocolFlv<<8 | SessionBaseTypePull
	SessionTypeTsSub             SessionType = SessionProtocolTs<<8 | SessionBaseTypeSub
	SessionTypePsPub             SessionType = SessionProtocolPs<<8 | SessionBaseTypePub

	SessionProtocolCustomize = 1
	SessionProtocolRtmp      = 2
	SessionProtocolRtsp      = 3
	SessionProtocolFlv       = 4
	SessionProtocolTs        = 5
	SessionProtocolPs        = 6

	SessionBaseTypePubSub = 1
	SessionBaseTypePub    = 2
	SessionBaseTypeSub    = 3
	SessionBaseTypePush   = 4
	SessionBaseTypePull   = 5

	SessionProtocolCustomizeStr = "CUSTOMIZE"
	SessionProtocolRtmpStr      = "RTMP"
	SessionProtocolRtspStr      = "RTSP"
	SessionProtocolFlvStr       = "FLV"
	SessionProtocolTsStr        = "TS"
	SessionProtocolPsStr        = "PS"

	SessionBaseTypePubSubStr = "PUBSUB"
	SessionBaseTypePubStr    = "PUB"
	SessionBaseTypeSubStr    = "SUB"
	SessionBaseTypePushStr   = "PUSH"
	SessionBaseTypePullStr   = "PULL"
)
