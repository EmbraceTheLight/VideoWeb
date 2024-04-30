package test

import (
	"VideoWeb/Utilities"
	"VideoWeb/logic"
	"fmt"
	"path"
	"strconv"
	"testing"
)

func TestVideoTime(t *testing.T) {

	time, _ := logic.GetVideoDuration(
		"/home/zey/ZeyGO/project/VideoWeb/resources/Videos/2024-04-07T011915.avi")
	//if err != nil {
	//	t.Error(err)
	//}
	fmt.Println(Utilities.SecondToTime(time))
}

func TestPath(t *testing.T) {
	p := "/home/zey/ZeyGO/project/VideoWeb/resources/Videos/2024-04-07T011915.avi"
	fmt.Println(path.Dir(p))
	conv, err := strconv.Unquote(`"` + "<?xml version=\\\"1.0\\\" encoding=\\\"utf-8\\\"?>\\n<MPD xmlns:xsi=\\\"http://www.w3.org/2001/XMLSchema-instance\\\"\\n\\txmlns=\\\"urn:mpeg:dash:schema:mpd:2011\\\"\\n\\txmlns:xlink=\\\"http://www.w3.org/1999/xlink\\\"\\n\\txsi:schemaLocation=\\\"urn:mpeg:DASH:schema:MPD:2011 http://standards.iso.org/ittf/PubliclyAvailableStandards/MPEG-DASH_schema_files/DASH-MPD.xsd\\\"\\n\\tprofiles=\\\"urn:mpeg:dash:profile:isoff-live:2011\\\"\\n\\ttype=\\\"static\\\"\\n\\tmediaPresentationDuration=\\\"PT1M31.1S\\\"\\n\\tmaxSegmentDuration=\\\"PT5.0S\\\"\\n\\tminBufferTime=\\\"PT24.6S\\\">\\n\\t<ProgramInformation>\\n\\t</ProgramInformation>\\n\\t<ServiceDescription id=\\\"0\\\">\\n\\t</ServiceDescription>\\n\\t<Period id=\\\"0\\\" start=\\\"PT0.0S\\\">\\n\\t\\t<AdaptationSet id=\\\"0\\\" contentType=\\\"video\\\" startWithSAP=\\\"1\\\" segmentAlignment=\\\"true\\\" bitstreamSwitching=\\\"true\\\" frameRate=\\\"30/1\\\" maxWidth=\\\"1280\\\" maxHeight=\\\"720\\\" par=\\\"16:9\\\" lang=\\\"und\\\">\\n\\t\\t\\t<Representation id=\\\"0\\\" mimeType=\\\"video/mp4\\\" codecs=\\\"avc1.64001f\\\" bandwidth=\\\"1956505\\\" width=\\\"1280\\\" height=\\\"720\\\" sar=\\\"1:1\\\">\\n\\t\\t\\t\\t<SegmentTemplate timescale=\\\"15360\\\" initialization=\\\"init-stream$RepresentationID$.m4s\\\" media=\\\"chunk-stream$RepresentationID$-$Number%05d$.m4s\\\" startNumber=\\\"1\\\">\\n\\t\\t\\t\\t\\t<SegmentTimeline>\\n\\t\\t\\t\\t\\t\\t<S t=\\\"0\\\" d=\\\"98816\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"135168\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"104960\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"112128\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"84992\\\" r=\\\"1\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"83456\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"84480\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"88576\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"82944\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"86528\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"189440\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"163840\\\" />\\n\\t\\t\\t\\t\\t</SegmentTimeline>\\n\\t\\t\\t\\t</SegmentTemplate>\\n\\t\\t\\t</Representation>\\n\\t\\t</AdaptationSet>\\n\\t\\t<AdaptationSet id=\\\"1\\\" contentType=\\\"audio\\\" startWithSAP=\\\"1\\\" segmentAlignment=\\\"true\\\" bitstreamSwitching=\\\"true\\\" lang=\\\"und\\\">\\n\\t\\t\\t<Representation id=\\\"1\\\" mimeType=\\\"audio/mp4\\\" codecs=\\\"mp4a.40.2\\\" bandwidth=\\\"128578\\\" audioSamplingRate=\\\"44100\\\">\\n\\t\\t\\t\\t<AudioChannelConfiguration schemeIdUri=\\\"urn:mpeg:dash:23003:3:audio_channel_configuration:2011\\\" value=\\\"2\\\" />\\n\\t\\t\\t\\t<SegmentTemplate timescale=\\\"44100\\\" initialization=\\\"init-stream$RepresentationID$.m4s\\\" media=\\\"chunk-stream$RepresentationID$-$Number%05d$.m4s\\\" startNumber=\\\"1\\\">\\n\\t\\t\\t\\t\\t<SegmentTimeline>\\n\\t\\t\\t\\t\\t\\t<S t=\\\"0\\\" d=\\\"220160\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"221184\\\" r=\\\"16\\\" />\\n\\t\\t\\t\\t\\t\\t<S d=\\\"40133\\\" />\\n\\t\\t\\t\\t\\t</SegmentTimeline>\\n\\t\\t\\t\\t</SegmentTemplate>\\n\\t\\t\\t</Representation>\\n\\t\\t</AdaptationSet>\\n\\t</Period>\\n</MPD>\\n" + `"`)
	//s := "as\nfg"
	//conv := strconv.Quote
	fmt.Println(conv)
	//conv, err = strconv.Unquote(conv)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Println(conv)
}
