package commonality

import (
	"easy-video-net/global"
	receive "easy-video-net/interaction/receive/commonality"
	response "easy-video-net/interaction/response/commonality"
	"easy-video-net/models/config/uploadMethod"
	"easy-video-net/models/contribution/video"
	"easy-video-net/models/users"
	"easy-video-net/models/users/attention"
	"easy-video-net/utils/conversion"
	"easy-video-net/utils/location"
	"easy-video-net/utils/oss"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"os"
	"strings"
)

var (
	//Temporary Location of document files
	Temporary = "assets/tmp/"
)

func OssSTS() (results interface{}, err error) {
	info, err := oss.GteStsInfo()
	if err != nil {
		return nil, err
	}
	res, err := response.GteStsInfo(info)
	if err != nil {
		return nil, fmt.Errorf("Response Failure")
	}
	return res, nil
}

func Upload(file *multipart.FileHeader, ctx *gin.Context) (results interface{}, err error) {
	//If the file size exceeds maxMemory, a temporary file is used to store the file data in multipart/form.
	err = ctx.Request.ParseMultipartForm(128)
	if err != nil {
		return
	}
	mForm := ctx.Request.MultipartForm
	//Uploading documents name
	var fileName string
	fileName = strings.Join(mForm.Value["name"], fileName)
	var fileInterface string
	fileInterface = strings.Join(mForm.Value["interface"], fileInterface)

	method := new(uploadMethod.UploadMethod)
	if !method.IsExistByField("interface", fileInterface) {
		return nil, fmt.Errorf("Upload interface does not exist")
	}
	if len(method.Path) == 0 {
		return nil, fmt.Errorf("Please contact the administrator to set the interface save path")
	}
	//取出文件
	index := strings.LastIndex(file.Filename, ".")
	suffix := file.Filename[index:]
	switch suffix {
	case ".jpg", ".jpeg", ".png", ".ico", ".gif", ".wbmp", ".bmp", ".svg", ".webp", ".mp4":
	default:
		return nil, fmt.Errorf("Illegal suffixes!")
	}
	if !location.IsDir(method.Path) {
		if err = os.MkdirAll(method.Path, 077); err != nil {
			return nil, fmt.Errorf("Failed to create save path")
		}
	}
	dst := method.Path + "/" + fileName
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		global.Logger.Warn("update headPortrait err")
		return nil, fmt.Errorf("Upload Failed")
	} else {
		return dst, nil
	}
}

func UploadSlice(file *multipart.FileHeader, ctx *gin.Context) (results interface{}, err error) {
	//If the file size exceeds maxMemory, a temporary file is used to store the file data in multipart/form.
	err = ctx.Request.ParseMultipartForm(128)
	if err != nil {
		return
	}
	mForm := ctx.Request.MultipartForm
	//Uploading documents name
	var fileName string
	fileName = strings.Join(mForm.Value["name"], fileName)
	var fileInterface string
	fileInterface = strings.Join(mForm.Value["interface"], fileInterface)

	method := new(uploadMethod.UploadMethod)
	if !method.IsExistByField("interface", fileInterface) {
		return nil, fmt.Errorf("Upload interface does not exist")
	}
	if len(method.Path) == 0 {
		return nil, fmt.Errorf("Please contact the administrator to set the interface save path")
	}
	if !location.IsDir(Temporary) {
		if err = os.MkdirAll(Temporary, 077); err != nil {
			return nil, fmt.Errorf("Failed to create save path")
		}
	}
	dst := Temporary + "/" + fileName
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		global.Logger.Warn("Possible cause of the fragment upload failure: ", err)
		return nil, fmt.Errorf("Upload Failed")
	} else {
		return dst, nil
	}
}

func UploadCheck(data *receive.UploadCheckStruct) (results interface{}, err error) {
	method := new(uploadMethod.UploadMethod)
	if !method.IsExistByField("interface", data.Interface) {
		return nil, fmt.Errorf("Upload method not configured")
	}
	list := make(receive.UploadSliceList, 0)
	path := method.Path + "/" + data.FileMd5
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		//File already exists
		return response.UploadCheckResponse(true, list, path)
	}
	//Taking out unuploaded slices
	for _, v := range data.SliceList {
		if _, err := os.Stat(Temporary + "/" + v.Hash); os.IsNotExist(err) {
			list = append(list, receive.UploadSliceInfo{
				Index: v.Index,
				Hash:  v.Hash,
			})
		}
	}
	return response.UploadCheckResponse(false, list, "")
}

func UploadMerge(data *receive.UploadMergeStruct) (results interface{}, err error) {
	method := new(uploadMethod.UploadMethod)
	if !method.IsExistByField("interface", data.Interface) {
		return nil, fmt.Errorf("Upload method not configured")
	}
	dst := method.Path + "/" + data.FileName
	list := make(receive.UploadSliceList, 0)
	path := method.Path + "/" + data.FileName
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		//File already exists directly return
		return dst, nil
	}
	//Taking out unuploaded slices
	for _, v := range data.SliceList {
		if _, err := os.Stat(Temporary + "/" + v.Hash); os.IsNotExist(err) {
			list = append(list, receive.UploadSliceInfo{
				Index: v.Index,
				Hash:  v.Hash,
			})
		}
	}
	if len(list) > 0 {
		return nil, fmt.Errorf("Segmentation not fully uploaded")
	}
	//Perform a merge operation
	cf, err := os.Create(dst)
	defer func(f *os.File) {
		if err := f.Close(); err != nil {
			global.Logger.Errorf("Failure to free memory for merge operation %d", err)
		}
	}(cf)
	if err != nil {
		return nil, fmt.Errorf("fail to save")
	}
	fileInfo, err := os.OpenFile(dst, os.O_APPEND, os.ModeSetuid)
	defer func(fileInfo *os.File) {
		if err := fileInfo.Close(); err != nil {
			global.Logger.Errorf("Closure of resources err : %d", err)
		}
	}(fileInfo)
	//合并操作
	for _, v := range data.SliceList {
		tmpFile, err := os.OpenFile(Temporary+"/"+v.Hash, os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		b, err := ioutil.ReadAll(tmpFile)
		if err != nil {
			fmt.Println(err)
		}
		if _, err := fileInfo.Write(b); err != nil {
			global.Logger.Errorf("Merge slice append error err : %d", err)
		}
		// 关闭分片
		if err := tmpFile.Close(); err != nil {
			global.Logger.Errorf("Close Split Error err : %d", err)
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			global.Logger.Errorf("Merge operation to delete temporary slice failed err : %d", err)
		}
	}
	return dst, nil
}

func UploadingMethod(data *receive.UploadingMethodStruct) (results interface{}, err error) {
	method := new(uploadMethod.UploadMethod)
	if method.IsExistByField("interface", data.Method) {
		return response.UploadingMethodResponse(method.Method), nil
	} else {
		return nil, fmt.Errorf("Upload method not configured")
	}
}

func UploadingDir(data *receive.UploadingDirStruct) (results interface{}, err error) {
	method := new(uploadMethod.UploadMethod)
	if method.IsExistByField("interface", data.Interface) {
		return response.UploadingDirResponse(method.Path), nil
	} else {
		return nil, fmt.Errorf("Upload method not configured")
	}
}

func GetFullPathOfImage(data *receive.GetFullPathOfImageMethodStruct) (results interface{}, err error) {
	path, err := conversion.SwitchIngStorageFun(data.Type, data.Path)
	if err != nil {
		return nil, err
	}
	return path, nil
}

func Search(data *receive.SearchStruct, uid uint) (results interface{}, err error) {
	switch data.Type {
	case "video":
		//Video serach
		list := new(video.VideosContributionList)
		err = list.Search(data.PageInfo)
		if err != nil {
			return nil, fmt.Errorf("Enquiry Failure")
		}
		res, err := response.SearchVideoResponse(list)
		if err != nil {
			return nil, fmt.Errorf("Response Failure")
		}
		return res, nil
		break
	case "user":
		list := new(users.UserList)
		err := list.Search(data.PageInfo)
		if err != nil {
			return nil, fmt.Errorf("Enquiry Failure")
		}
		aids := make([]uint, 0)
		if uid != 0 {
			//In case of user login
			al := new(attention.AttentionsList)
			err = al.GetAttentionList(uid)
			if err != nil {
				return nil, fmt.Errorf("Failed to get follow list")
			}
			for _, v := range *al {
				aids = append(aids, v.AttentionID)
			}
		}
		res, err := response.SearchUserResponse(list, aids)
		return res, nil
		break
	default:
		return nil, fmt.Errorf("Unmatched types")
	}
	return
}
