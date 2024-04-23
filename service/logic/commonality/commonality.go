package commonality

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"os"
	"simple-video-net/global"
	receive "simple-video-net/interaction/receive/commonality"
	response "simple-video-net/interaction/response/commonality"
	"simple-video-net/models/sundry/upload"
	"simple-video-net/models/contribution/video"
	"simple-video-net/models/users"
	"simple-video-net/models/users/attention"
	"simple-video-net/utils/conversion"
	"simple-video-net/utils/location"
	"simple-video-net/utils/oss"
	"simple-video-net/utils/validator"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
)

var (
	//Temporary Location of document files
	Temporary = "assets/tmp"
)

func OssSTS() (results interface{}, err error) {
	info, err := oss.GteStsInfo()
	if err != nil {
		global.Logger.Errorf("Failed to obtain OssSts key. Error reason :%s", err.Error())
		return nil, fmt.Errorf("Failed to obtain")
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

	method := new(upload.upload)
	if !method.IsExistByField("interface", fileInterface) {
		return nil, fmt.Errorf("Upload interface does not exist")
	}
	if len(method.Path) == 0 {
		return nil, fmt.Errorf("Please contact the administrator to set the interface save path")
	}
	//Take out files
	index := strings.LastIndex(fileName, ".")
	suffix := fileName[index:]
	err = validator.CheckVideoSuffix(suffix)
	if err != nil {
		return nil, fmt.Errorf("Illegal suffixes!")
	}
	if err = os.MkdirAll(method.Path, 0775); err != nil {
		if err = os.MkdirAll(method.Path, 077); err != nil {
			global.Logger.Errorf("Failed to create file with error path. The creation path is：%s wrong reason : %s", method.Path, err.Error())
			return nil, fmt.Errorf("Failed to create save path")
		}
	}
	dst := method.Path + "/" + fileName
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		global.Logger.Errorf("Failed to save the file. The save path is：%s ,wrong reason : %s", dst, err.Error())
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

	method := new(upload.upload)
	if !method.IsExistByField("interface", fileInterface) {
		return nil, fmt.Errorf("Upload interface does not exist")
	}
	if len(method.Path) == 0 {
		return nil, fmt.Errorf("Please contact the administrator to set the interface save path")
	}
	if err = os.MkdirAll(Temporary, 0775); err != nil {
		if err = os.MkdirAll(Temporary, 077); err != nil {
			global.Logger.Errorf("Failed to create file with error path. The creation path is：%s", method.Path)
			return nil, fmt.Errorf("Failed to create save path")
		}
	}
	dst := Temporary + "/" + fileName
	err = ctx.SaveUploadedFile(file, dst)
	if err != nil {
		global.Logger.Errorf("Multipart upload failed to save. The save path is：%s ,wrong reason : %s ", dst, err.Error())
		return nil, fmt.Errorf("Upload Failed")
	} else {
		return dst, nil
	}
}

func UploadCheck(data *receive.UploadCheckStruct) (results interface{}, err error) {
	method := new(upload.upload)
	if !method.IsExistByField("interface", data.Interface) {
		return nil, fmt.Errorf("Upload method not configured")
	}
	list := make(receive.UploadSliceList, 0)
	path := method.Path + "/" + data.FileMd5
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		//File already exists
		global.Logger.Infof("upload files %s existed", data.FileMd5)
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
	method := new(upload.upload)
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
		global.Logger.Warnf("upload files %s Not all shards are uploaded", data.FileName)
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
	//merge operation
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
			global.Logger.Errorf("Merge shard append error Cause of error : %d", err)
		}
		// Turn off sharding
		if err := tmpFile.Close(); err != nil {
			global.Logger.Errorf("Close sharding error Cause of error : %d", err)
		}
		if err := os.Remove(tmpFile.Name()); err != nil {
			global.Logger.Errorf("The merge operation failed to delete temporary shards. Reason for the error. : %d", err)
		}
	}
	return dst, nil
}

func UploadingMethod(data *receive.UploadingMethodStruct) (results interface{}, err error) {
	method := new(upload.upload)
	if method.IsExistByField("interface", data.Method) {
		return response.UploadingMethodResponse(method.Method), nil
	} else {
		return nil, fmt.Errorf("Upload method not configured")
	}
}

func UploadingDir(data *receive.UploadingDirStruct) (results interface{}, err error) {
	method := new(upload.Upload)
	if method.IsExistByField("interface", data.Interface) {
		return response.UploadingDirResponse(method.Path, method.Quality), nil
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
				global.Logger.Errorf("user id %d Failed to obtain the following list, error reason : %s ", uid, err.Error())
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

func RegisterMedia(data *receive.RegisterMediaStruct) (results interface{}, err error) {
	path, _ := conversion.SwitchIngStorageFun(data.Type, data.Path)
	//Register media assets
	registerMediaBody, err := oss.RegisterMediaInfo(path, "video", time.Now().String())
	if err != nil {
		return nil, fmt.Errorf("Failed to register media assets")
	}
	return registerMediaBody.MediaId, nil
}