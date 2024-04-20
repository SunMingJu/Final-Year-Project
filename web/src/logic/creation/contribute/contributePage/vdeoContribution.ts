import { getuploadingMethod, registerMedia } from "@/apis/commonality";
import { createVideoContribution, updateVideoContribution } from "@/apis/contribution";
import globalScss from "@/assets/styles/global/export.module.scss";
import { useEditVideoStore } from '@/store/creation';
import { useUserStore } from "@/store/main";
import { Ref, UnwrapNestedRefs, nextTick, reactive, ref } from "vue";
import { CreateVideoContributionReq, UpdateVideoContributionReq, uploadFileformation, vdeoContributionForm } from "@/types/creation/contribute/contributePage/vdeoContribution";
import { timetoRFC3339 } from "@/utils/conversion/timeConversion";
import { fileReader } from "@/utils/fun/fun";
import { uploadFile } from '@/utils/upload/upload';
import { validateVideoIntroduce, validateVideoTitle } from "@/utils/validate/validate";
import { ElInput, FormInstance, UploadProps, UploadRawFile, UploadRequestOptions } from 'element-plus';
import Swal from 'sweetalert2';
import { Ref, UnwrapNestedRefs, nextTick, reactive, ref } from "vue";
import { Router, useRouter } from 'vue-router';

export const useVdeoContributionProp = () => {
    const userStore = useUserStore()
    const editVideoStore = useEditVideoStore()
    const formRef = ref<FormInstance>()
    const ruleFormRef = ref<FormInstance>()
    const router = useRouter()
    const video = ref() //Uploaded video information
    const form = reactive(<vdeoContributionForm>{
        id: 0,
        isShow: false,
        title: '',
        type: false,
        labelInputVisible: false,
        labelText: "",
        label: [],
        introduce: "",
        videoDuration: 0,
    })
    const uploadFileformation = reactive(<uploadFileformation>{
        progress: 0,
        FileUrl: '',
        uploadUrl: "",
        interface: "videoContribution",
        uploadType: "",
        action: "#",
    })

    const uploadCoveration = reactive(<uploadFileformation>{
        progress: 0,
        FileUrl: '',
        uploadUrl: "",
        interface: "videoContributionCover",
        uploadType: "",
        action: "#",
    })
    const labelInputRef = ref<InstanceType<typeof ElInput>>()
    return {
        ruleFormRef,
        userStore,
        formRef,
        form,
        router,
        uploadFileformation,
        uploadCoveration,
        labelInputRef,
        video,
        editVideoStore
    }
}

//上Video upload processing
export const useHandleFileMethod = (uploadFileformation: uploadFileformation, form: vdeoContributionForm, video: Ref) => {

    const handleFileSuccess: UploadProps['onSuccess'] = async (
        response,
        uploadFile
    ) => {
        uploadFileformation.FileUrl = URL.createObjectURL(uploadFile.raw!)
        //video ready event
        video.value.onloadedmetadata = () => {
            //Modify video duration
            form.videoDuration = Math.round(video.value.duration)
        }
        const readerInfo = await fileReader(uploadFile.raw!)
        video.value.src = readerInfo?.result
    }

    const handleFileError: UploadProps['onError'] = (
        response,
    ) => {
        console.log("upload failed")
        Swal.fire({
            title: "upload failed",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",

        })
        console.log(response)

    }

    //Processing before uploading
    const beforeFileUpload: UploadProps['beforeUpload'] = async (rawFile: UploadRawFile) => {
        return true
    }

    //Modify default request
    const RedefineUploadFile = async (params: UploadRequestOptions) => {
        try {
            //Larger than 30mb fragments
            let fragment = params.file.size > 30 * 1024 * 1024 ? true : false
            form.isShow = !form.isShow
            const response = await uploadFile(uploadFileformation, params.file, fragment)
            console.log(response)
            uploadFileformation.uploadUrl = response.path
            //Register media resources
            let media = await registerMedia(<RegisterMediaReq>{
                type: uploadFileformation.uploadType,
                path: response.path
            })
            uploadFileformation.media = media.data
        } catch (err) {
            console.log(err)
            form.isShow = false
            Swal.fire({
                title: "Failed to upload video",
                heightAuto: false,
                confirmButtonColor: globalScss.colorButtonTheme,
                icon: "error",
            })
        }
    }

    return {
        handleFileSuccess,
        beforeFileUpload,
        handleFileError,
        RedefineUploadFile
    }

}

//Upload cover processing
export const useHandleCoverMethod = (uploadCoveration: uploadFileformation, form: vdeoContributionForm) => {

    const handleFileSuccess: UploadProps['onSuccess'] = (
        response,
        uploadFile
    ) => {
        uploadCoveration.FileUrl = URL.createObjectURL(uploadFile.raw!)
    }

    const handleFileError: UploadProps['onError'] = (
        response,
    ) => {
        console.log("upload failed")
        Swal.fire({
            title: "upload failed",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",

        })
        console.log(response)

    }

    //Processing before uploading
    const beforeFileUpload: UploadProps['beforeUpload'] = async (rawFile: UploadRawFile) => {
        return await new Promise<boolean>((resolve, reject) => {
            //Determine size
            if (rawFile.size / 1024 / 1024 > 2) {
                Swal.fire({
                    title: "Cover size cannot be larger than 2 m",
                    heightAuto: false,
                    icon: "error",

                })
                reject(false);
            }
            //Determine size
            let reader = new FileReader();
            reader.readAsDataURL(rawFile);
            reader.onload = () => {
                // Let the src of the img tag in the page point to the read path
                let img = new Image();
                img.src = reader.result as string;
                img.onload = () => {
                    console.log(img.width);
                    console.log(img.height);
                    if (img.width < 960 || img.height < 600) {
                        Swal.fire({
                            title: "Please upload images above 960*600 size",
                            heightAuto: false,
                            confirmButtonColor: globalScss.colorButtonTheme,
                            icon: "error",
                        });
                        reject(false);
                    } else {
                        resolve(true);
                    }
                };
            };
        })
    }

    //Modify default request
    const RedefineUploadFile = async (params: UploadRequestOptions) => {
        try {
            const response = await uploadFile(uploadCoveration, params.file)
            console.log(response)
            uploadCoveration.uploadUrl = response.path
        } catch (err) {
            console.log(err)
            Swal.fire({
                title: "Failed to upload video cover",
                heightAuto: false,
                confirmButtonColor: globalScss.colorButtonTheme,
                icon: "error",
            })
        }
    }

    return {
        handleFileSuccess,
        beforeFileUpload,
        handleFileError,
        RedefineUploadFile
    }

}


export const userLabelHandlMethod = (form: vdeoContributionForm, labelInputRef: Ref) => {
    const handleClose = (tag: string) => {
        form.label.splice(form.label.indexOf(tag), 1)
    }

    const showInput = () => {
        form.labelInputVisible = true
        nextTick(() => {
            labelInputRef.value!.input!.focus()
        })
    }

    const handleInputConfirm = () => {
        if (form.labelText) {
            form.label.push(form.labelText)
        }
        form.labelInputVisible = false
        form.labelText = ''
    }

    return {
        handleClose,
        showInput,
        handleInputConfirm
    }

}
export const useSaveData = async (form: vdeoContributionForm, uploadFileformation: uploadFileformation, uploadCoveration: uploadFileformation, formEl: FormInstance | undefined, router: Router, props: any) => {
    if (!formEl) return;
    await formEl.validate(async (valid, fields) => {
        if (valid) {
            try {
                if (!uploadCoveration.uploadUrl) throw "Please upload the cover first"
                //Determine operation type
                if (props.type == "edit") {
                    //Update video
                    var updateRequistData = <UpdateVideoContributionReq>{
                        id: form.id,
                        cover: uploadCoveration.uploadUrl,
                        coverUploadType: uploadCoveration.uploadType,
                        title: form.title,
                        reprinted: form.type,
                        label: form.label,
                        introduce: form.introduce,
                    }
                    await updateVideoContribution(updateRequistData)
                } else {
                    //Create video
                    if (!uploadFileformation.uploadUrl) throw "Upload not completed"
                    var createRequistData = <CreateVideoContributionReq>{
                        id: form.id,
                        video: uploadFileformation.uploadUrl,
                        videoUploadType: uploadFileformation.uploadType,
                        cover: uploadCoveration.uploadUrl,
                        coverUploadType: uploadCoveration.uploadType,
                        title: form.title,
                        reprinted: form.type,
                        label: form.label,
                        introduce: form.introduce,
                        videoDuration: form.videoDuration,
                        media: uploadFileformation.media
                    }
                    await createVideoContribution(createRequistData)
                }
                let swalTitle = props.type == "edit" ? "update completed" : "Posted successfully"
                Swal.fire({
                    title: swalTitle,
                    confirmButtonColor: globalScss.colorButtonTheme,
                    heightAuto: false,
                    icon: "success",
                    preConfirm: () => {
                        router.push({ name: "Creation" })
                    }
                })
            } catch (err) {
                console.log(err)
                Swal.fire({
                    title: err as string,
                    confirmButtonColor: globalScss.colorButtonTheme,
                    heightAuto: false,
                    icon: "warning",

                })
            }
        } else {
            console.log('error submit!', fields)
        }
    })
}




export const useInit = async (uploadFileformation: uploadFileformation, uploadCoveration: uploadFileformation, form: UnwrapNestedRefs<vdeoContributionForm>, props: any, editVideoStore: any) => {
    try {
        //Get the request method of the current interface
        const updataMenhod = (await getuploadingMethod(<GetUploadingMethodReq>{
            method: uploadFileformation.interface
        })).data as GetUploadingMethodRes
        uploadFileformation.uploadType = updataMenhod.type
        const updataMenhodCover = (await getuploadingMethod(<GetUploadingMethodReq>{
            method: uploadCoveration.interface
        })).data as GetUploadingMethodRes

        //Determine current mode
        if (props.type == "edit") {
            //edit mode
            form.isShow = true
            uploadFileformation.progress = 100
            form.id = editVideoStore.editVideoData.videoID
            form.title = editVideoStore.editVideoData.title
            form.label = editVideoStore.editVideoData.label
            form.type = editVideoStore.editVideoData.reprinted
            form.introduce = editVideoStore.editVideoData.introduce
            uploadCoveration.FileUrl = editVideoStore.editVideoData.cover
            uploadCoveration.uploadUrl = editVideoStore.editVideoData.cover_url
            uploadCoveration.uploadType = editVideoStore.editVideoData.coverUploadType
        }



        uploadCoveration.uploadType = updataMenhod.type

    } catch (err) {
        console.log(err)
        Swal.fire({
            title: "Failed to get upload method",
            heightAuto: false,
            confirmButtonColor: globalScss.colorButtonTheme,
            icon: "error",
        })
    }
}

//form validation
export const useRules = () => {
    const videoContributionRules = reactive({
        title: [{ validator: validateVideoTitle, trigger: 'change' }],
        introduce: [{ validator: validateVideoIntroduce, trigger: 'change' }]
    });
    return {
        videoContributionRules,
    };
}; 
