import { getuploadingDir, uploadCheck, uploadFile, uploadMerge, UploadSliceFile } from "@/apis/commonality";
import { GetuploadingDirReq, UploadCheckReq, UploadMergeReq, UploadSliceList } from "@/types/commonality/commonality";
import { FileSliceUpload, FileUpload } from "@/types/idnex";
import Compressor from "compressorjs";
import { ref, watch } from "vue";
import { fileHash, fileSuffix } from "./fileManipulation";
import { getSliceFile } from "./getSliceFile";
/**
 * Upload files locally
 * @param file File object
 * @returns {Promise<{name:string,host:string}>}
 */
export const localUpload = async (file: File, uploadConfig: FileUpload, dir: string, fragment?: boolean): Promise<any> => {
    return new Promise(async (resolve, reject) => {
        if (!fragment) {
            //Upload directly
            // Calculate file Hash to avoid redundant file uploads. The purpose of doing this is to occupy as little space as possible
            const name = await fileHash(file) + fileSuffix(file.name)
            const formData = new FormData()
            const key = `${name}`
            formData.append('interface', uploadConfig.interface)
            formData.append('name', name)
            formData.append('file', file)
            try {
                const response = await uploadFile(formData, uploadConfig)
                resolve({ path: response.data as string })
                console.log(response)
            } catch (err) {
                console.log(err)
                reject({ msg: 'upload failed' })
            }
        } else {
            const uploadCheckFun = async () => {
                // Calculate file Hash to avoid redundant file uploads. The purpose of doing this is to occupy as little space as possible
                const name = await fileHash(file) + fileSuffix(file.name)
                //total slices
                let sliceArr = await getSliceFile(file, 1, name)
                let sliceList = ref(<UploadSliceList>[])
                sliceArr.filter((item) => {
                    sliceList.value.push({
                        index: item.index,
                        hash: item.hash
                    })
            const uploadCheckResponse = await uploadCheck(<UploadCheckReq>{
                    file_md5: name,
                    interface: uploadConfig.interface,
                    slice_list: sliceArr
                })
                if (uploadCheckResponse.data?.is_upload) {
                    uploadConfig.progress = 100
                    return resolve({ path: uploadCheckResponse.data?.path })
                }
               let notUploadIndex: number[] = []
                uploadCheckResponse.data?.list.filter((item) => {
                    notUploadIndex.push(item.index)
                })
                //Get unuploaded shards
                const notUploadSlice = sliceArr.filter((item) => {
                    return (notUploadIndex.includes(item.index))
                })
                //Set the upload progress to 100%
                sliceArr = sliceArr.filter((item) => {
                    if (!notUploadIndex.includes(item.index)) {
                        item.progress = 100
                    }
                    return item
                })
                console.log("So the slices that need to be uploaded", sliceArr)
                    console.log("Unuploaded slices", notUploadSlice)
                    
                    let promiseArr = []

                for (let i = 0; i < notUploadSlice.length; i++) {
                    const formData = new FormData()
                    formData.append('interface', uploadConfig.interface)
                    formData.append('name', notUploadSlice[i].hash)
                    formData.append('file', notUploadSlice[i].file)
                    const p = new Promise<void>(async (resolve, reject) => {
                        const slice = ref(<FileSliceUpload>{
                            index: i,
                            progress: 0, //Upload progress
                            size: notUploadSlice[i].size
                        })
                        let w = watch(() => { slice.value.progress }, () => {
                            //calculate
                            sliceArr.filter((item, index, arr) => {
                                if (item.index === notUploadSlice[i].index) {
                                    sliceArr[index].progress = slice.value.progress
                                    updataProgress(sliceArr, uploadConfig)
                                }
                                return item;
                            })
                            if (slice.value.progress === 100) {
                                w()
                                resolve()
                                return
                            }
                        }, { deep: true })
                        await UploadSliceFile(formData, slice.value)
                        .catch((error) => {
                            reject(error)
                        })
                    })
                    promiseArr.push(p)
                }
                    try {
                    await Promise.all(promiseArr)
                    console.log('All parts upload completed')
                    //All shards were successfully uploaded and merged
                    const uploadMergeResponse = await uploadMerge(<UploadMergeReq>{
                        file_name: name,
                        interface: uploadConfig.interface,
                        slice_list: sliceArr
                    })
                    uploadConfig.progress = 100
                    return resolve({ path: uploadMergeResponse.data })
                } catch (err) {
                    return new Promise((_, reject) => {
                        console.log('There are unuploaded fragments')
                    uploadCheckFun()
                    })
                }

            }

            const updataProgress = (sliceArr: Array<any>, uploadConfig: FileUpload) => {
                const totalSize = file.size
                let loadSize = 0

                sliceArr.filter((item) => {
                    //Calculate the upload size of each piece
                    loadSize += (item.sliceSizeInByte * item.progress) / 100
                })
                let progress = Math.round(loadSize / totalSize * 100)
                uploadConfig.progress = progress
            }
            return uploadCheckFun()
        }

    });


} 