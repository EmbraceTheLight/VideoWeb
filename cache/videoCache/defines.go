package videoCache

import (
	"context"
)

type VideoBasic struct {
	Key       string
	VideoInfo map[string]any
}

type BarrageInfo struct {
	key         string
	barrageInfo map[string]any
}

type VideoTags struct {
	Key  string
	Tags []string
}

type VideoBarrages struct {
	Key      string
	Barrages []string
}

type VideoComments struct {
	Key      string
	Comments []string
}
type VideoCache struct {
	VBasic    VideoBasic
	VBarrages VideoBarrages
	VTags     VideoTags
	VComments VideoComments
}

func (video *VideoCache) initVideoCache() {
	video.VBasic = VideoBasic{VideoInfo: make(map[string]any)}
	video.VBarrages = VideoBarrages{Barrages: make([]string, 0)}
	video.VTags = VideoTags{Tags: make([]string, 0)}
	video.VComments = VideoComments{Comments: make([]string, 0)}
}

func (vbi *BarrageInfo) GetKey() string {
	return vbi.key
}

func (vbi *BarrageInfo) GetMap() map[string]any {
	return vbi.barrageInfo
}

func (vbasic *VideoBasic) Set(k string, v any) {
	vbasic.VideoInfo[k] = v
}

func (vbasic *VideoBasic) Get(k string) any {
	return vbasic.VideoInfo[k]
}

func (vbasic *VideoBasic) del(k string) {
	delete(vbasic.VideoInfo, k)
}

func (vbi *BarrageInfo) Set(k string, v any) {
	vbi.barrageInfo[k] = v
}

func (vbi *BarrageInfo) Get(k string) any {
	return vbi.barrageInfo[k]
}

func (vbi *BarrageInfo) del(k string) {
	delete(vbi.barrageInfo, k)
}

func (vt *VideoTags) Set(v any) {
	set(vt.Tags, v.(string))
}
func (vt *VideoTags) GetAll() []string { return vt.Tags }
func (vt *VideoTags) GetOne(ctx context.Context, member string) string {
	return getOne(ctx, "cache.videoCache.VideoTags->GetOne", vt.Key, member)
}
func (vt *VideoTags) Del(v any) {
	del(vt.Tags, v.(string))
}

func (vb *VideoBarrages) Set(v any) {
	set(vb.Barrages, v.(string))
}

func (vb *VideoBarrages) GetAll() []string {
	return vb.Barrages
}

func (vb *VideoBarrages) GetOne(ctx context.Context, member string) string {
	return getOne(ctx, "cache.videoCache.VideoBarrages->GetOne", vb.Key, member)
}

func (vb *VideoBarrages) Del(v any) {
	del(vb.Barrages, v.(string))
}

func (vc *VideoComments) Set(v any) {
	set(vc.Comments, v.(string))
}

func (vc *VideoComments) GetAll() []string {
	return vc.Comments
}

func (vc *VideoComments) GetOne(ctx context.Context, member string) string {
	return getOne(ctx, "cache.videoCache.VideoComments->GetOne", vc.Key, member)
}

func (vc *VideoComments) Del(v any) {
	del(vc.Comments, v.(string))
}
