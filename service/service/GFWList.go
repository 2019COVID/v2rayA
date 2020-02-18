package service

import (
	"V2RayA/extra/copyfile"
	"V2RayA/extra/download"
	"V2RayA/model/v2ray"
	"V2RayA/model/v2ray/asset"
	"V2RayA/persistence/configure"
	"V2RayA/tools/files"
	"V2RayA/tools/httpClient"
	"github.com/PuerkitoBio/goquery"
	gonanoid "github.com/matoous/go-nanoid"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var gfwlistupdatetime *time.Time
var gfwListMutex sync.Mutex

func GetRemoteGFWListUpdateTime(c *http.Client) (t time.Time, err error) {
	gfwListMutex.Lock()
	defer gfwListMutex.Unlock()
	if gfwlistupdatetime != nil {
		return *gfwlistupdatetime, nil
	}
	resp, err := httpClient.HttpGetUsingSpecificClient(c, "https://mirrors.0diis.com/github/ToutyRater/V2Ray-SiteDAT/contributors/master/geofiles/h2y.dat")
	if err != nil {
		return time.Time{}, errors.New("获取GFWList最新版本时间失败")
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	timeRaw, ok := doc.Find("relative-time").First().Attr("datetime")
	if !ok {
		log.Println(doc.Html())
		return time.Time{}, errors.New("获取最新GFWList版本失败")
	}
	t, err = time.Parse(time.RFC3339, timeRaw)
	gfwlistupdatetime = &t
	return t, err
}
func IsUpdate() (update bool, remoteTime time.Time, err error) {
	//c, err := httpClient.GetHttpClientAutomatically()
	//if err != nil {
	//	return
	//}
	remoteTime, err = GetRemoteGFWListUpdateTime(http.DefaultClient)
	if err != nil {
		return
	}

	if !asset.IsH2yExists() {
		//本地文件不存在，那远端必定比本地新
		return false, remoteTime, nil
	}
	//本地文件存在，检查本地版本是否比远端还新
	t, err := asset.GetH2yModTime()
	if err != nil {
		return
	}
	if t.After(remoteTime) {
		//那确实新
		update = true
		return
	}
	return
}

func UpdateLocalGFWList() (localGFWListVersionAfterUpdate string, err error) {
	id, _ := gonanoid.Nanoid()
	i := 0
	for {
		err = download.Pget("https://github.com/ToutyRater/V2Ray-SiteDAT/raw/master/geofiles/h2y.dat", "/tmp/h2y.dat."+id)
		if err != nil && i < 2 {
			//最多重试2次
			i++
			continue
		}
		break
	}
	if err != nil {
		return
	}
	err = copyfile.CopyFile("/tmp/h2y.dat."+id, asset.GetV2rayLocationAsset()+"/h2y.dat")
	if err != nil {
		return
	}
	err = os.Chmod(asset.GetV2rayLocationAsset()+"/h2y.dat", os.FileMode(0755))
	if err != nil {
		return
	}
	t, err := files.GetFileModTime(asset.GetV2rayLocationAsset() + "/h2y.dat")
	if err == nil {
		localGFWListVersionAfterUpdate = t.Format("2006-01-02")
	}
	return
}

func CheckAndUpdateGFWList() (localGFWListVersionAfterUpdate string, err error) {
	update, tRemote, err := IsUpdate()
	if err != nil {
		return
	}
	if update {
		return "", errors.New(
			"目前最新版本为" + tRemote.Format("2006-01-02") + "，您的本地文件已最新，无需更新",
		)
	}

	/* 更新h2y.dat */
	localGFWListVersionAfterUpdate, err = UpdateLocalGFWList()
	if err != nil {
		return
	}
	setting := configure.GetSettingNotNil()
	if v2ray.IsV2RayRunning() && //正在使用GFWList模式再重启
		(setting.Transparent == configure.TransparentGfwlist ||
			setting.Transparent == configure.TransparentClose && setting.PacMode == configure.GfwlistMode) {
		err = v2ray.UpdateV2rayWithConnectedServer()
	}
	return
}
