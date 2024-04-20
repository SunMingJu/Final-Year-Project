import { getuploadingDir } from "@/apis/commonality";
import { GetuploadingDirReq } from "@/types/commonality/commonality";
import { FileUpload } from "@/types/idnex";
import { compressFile, isImageFile } from "./fileManipulation";
import { localUpload } from "./local";
import { ossUpload } from "./oss";
export const uploadFile = async (config: FileUpload, rawFile: File, fragment?: boolean): Promise<{ path: string }> => {
    let res
    //Default false
    if (fragment == undefined) fragment = false
    //Get the save path and image quality
    const response = await getuploadingDir(<GetuploadingDirReq>{
        interface: config.interface
    })
    let dir = response.data?.path as string
    let quality = response.data?.quality as number
    //Compress Pictures
    if (isImageFile(rawFile)) {
        try {
            const compressedFile = await compressFile(rawFile, quality as number)
            // Operations after successful compression
            rawFile = compressedFile as File
        } catch (err) {
            // What to do after compression fails
            console.log('Compression failed!', err);
        }
    }
    switch (config.uploadType) {
        case "aliyunOss":
            res = ossUpload(rawFile, config, dir, fragment)
            break;
        case "local":
            res = localUpload(rawFile, config, dir, fragment)
            break;
    }
    return res as Promise<{ path: string }>
}

