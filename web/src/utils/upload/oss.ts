import { getuploadingDir, gteossSTS } from "@/apis/commonality"
import { useGlobalStore } from "@/store/main"
import { GetuploadingDirReq } from "@/types/commonality/commonality"
import { FileUpload, OssSTSInfo } from "@/types/idnex"
import OSS from 'ali-oss'
import { fileHash, fileSuffix } from "./fileManipulation"


//Initialize sts
export const initOssSTS = async (_interface: string): Promise<OssSTSInfo> => {
    return new Promise((resolve, reject) => {
        // Get configuration from local localstore
        const globalStore = useGlobalStore()
        const conf = globalStore.ossData
        if (conf) {
            //  If the configuration exists and the expiration time is one minute and one second later, this configuration will be returned.
            const now = new Date().getTime() / 1000
            console.log(conf)
            if (conf.expirationTime - 600 > now) {
                resolve(conf)
                return
            }
        }
        // Request interface to return configuration data
        gteossSTS()
            .then((res) => {
                if (res.code == 200) {
                    if (!res.data) return false
                    let info = res.data
                    // Write configuration data to local store
                    globalStore.setOssInfo(<OssSTSInfo>{
                        region: info.region,
                        accessKeyId: info.access_key_id,
                        accessKeySecret: info.access_key_secret,
                        stsToken: info.sts_token,
                        bucket: info.bucket,
                        expirationTime: info.expiration_time
                    })
                    resolve(globalStore.ossData);
                } else {
                    reject(res)
                }
            })
            .catch((err) => {
                console.log(err);
                reject(err)
            })
    })
}


/**
 * Upload files to oss
 * @param file File object
 * @returns {Promise<{name:string,host:string}>}
 */
export const ossUpload = (file: File, uploadConfig: FileUpload, fragment: boolean): Promise<{ path: string }> => {
    return new Promise((resolve, reject) => {
        initOssSTS(uploadConfig.interface)
            .then(async (ossSts) => {
                //Get save path
                const response = await getuploadingDir(<GetuploadingDirReq>{
                    interface: uploadConfig.interface
                })
                let dir = response.data?.path
                const name = await fileHash(file) + fileSuffix(file.name)
                const key = `${dir}${name}`
                // Initialize Alibaba Cloud oss ​​client
                const client = new OSS({
                    region: ossSts.region,
                    accessKeyId: ossSts.accessKeyId,
                    accessKeySecret: ossSts.accessKeySecret,
                    stsToken: ossSts.stsToken,
                    bucket: ossSts.bucket,
                });
                console.log(fragment)
                if (!fragment) {
                    console.log("普通上传")
                    //In order to be able to display the progress bar, multi-part uploading is also performed.
                    var checkpoint = getCheckpoint(name);
                    const options = {
                        checkpoint: checkpoint,
                        progress: (p: any, cpt: any) => {
                            console.log("上传进度", p)
                            uploadConfig.progress = Math.round(p * 100)
                            saveCheckpoint(name, cpt);
                        },
                        mime: "text/plain",
                        // Set the number of concurrently uploaded shards.
                        parallel: 4,
                        // Set the shard size. The default value is 1 MB and the minimum value is 100 KB.
                        partSize: 200 * 1024,
                    };

                    try {
                        const res = await client.multipartUpload(`${dir}${name}`, file, {
                            ...options,
                        });
                        console.log(res);
                        deleteCheckpoint(name);
                        resolve({ path: key })
                    } catch (err) {
                        console.log(err);
                        deleteCheckpoint(name);
                        reject({ msg: 'upload failed' })
                    }
                } else {
                    console.log("Multipart upload")
                    //Multiple upload plus breakpoint resume
                    var checkpoint = getCheckpoint(name);
                    const options = {
                        checkpoint: checkpoint,
                        //Get the multipart upload progress, breakpoints and return values.
                        progress: (p: any, cpt: any) => {
                            saveCheckpoint(name, cpt);
                            console.log(cpt)
                            uploadConfig.progress = Math.round(p *100)
                        },
                        //Set the number of concurrently uploaded shards.
parallel: 4,
                        //Set the shard size. The default value is 1 MB and the minimum value is 100 KB.
                        partSize: 1 * 1024 * 1024,
                        mime: "text/plain",
                    };

                    try {
                        const res = await client.multipartUpload(`${dir}${name}`, file, {
                            ...options,
                        });
                        deleteCheckpoint(name);
                        resolve({ path: key })

                        console.log(res)
                    } catch (err) {
                        deleteCheckpoint(name);
                        console.log(err);
                        reject({ msg: 'upload failed' })
                    }
                }
            })
            .catch((err) => {
                console.log(err);
                reject({ msg: 'upload failed' })
            })
    })
}

// Save upload breakpoint
function saveCheckpoint(key: string, checkpoint: any) {
    localStorage.setItem(key, JSON.stringify(checkpoint));
}

// Get upload breakpoint
function getCheckpoint(key: string,) {
    var checkpoint = localStorage.getItem(key);
    if (checkpoint) {
        return JSON.parse(checkpoint);
    } else {
        return null;
    }
}

// Delete upload breakpoint
function deleteCheckpoint(key: string) {
    localStorage.removeItem(key);
}