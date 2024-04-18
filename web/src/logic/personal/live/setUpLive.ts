import { getuploadingMethod } from "@/apis/commonality";
import { getLiveDataRequist, saveLiveDataRequist } from "@/apis/personal";
import globalScss from "@/assets/styles/global/export.module.scss";
import { useUserStore } from "@/store/main";
import { GetUploadingMethodReq, GetUploadingMethodRes } from "@/types/commonality/commonality";
import { GetLiveDataRes, LiveInformation, SaveLiveDataReq } from "@/types/personal/live/setUp";
import { getLocation } from "@/utils/conversion/stringConversion";
import { uploadFile } from '@/utils/upload/upload';
import { validateLiveTitle } from "@/utils/validate/validate";
import { FormInstance, UploadProps, UploadRequestOptions } from 'element-plus';
import Swal from 'sweetalert2';
import { reactive, ref } from "vue";
import useClipboard from "vue-clipboard3";

export const useLiveInfoProp = () => {
    const userStore = useUserStore()
    const saveDateFormRef = ref<FormInstance>()
const liveInformationForm = reactive<LiveInformation>({
        FileUrl: '',
        uploadUrl: "",
        interface: "liveCover",
        title: "",
        uploadType: "",
        action: "#",
        Progress: 0
    });

    //Define the original data of the request result
    const rawData = reactive<GetLiveDataRes>({
        title: "",
        img: ""
    })

    return {
        userStore,
        liveInformationForm,
        saveDateFormRef,
        rawData
    }
}
export const useHandleFileMethod = (liveInformationForm: LiveInformation) => {

    const handleFileSuccess: UploadProps['onSuccess'] = (
        response,
        uploadFile
    ) => {
        liveInformationForm.FileUrl = URL.createObjectURL(uploadFile.raw!)
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


    const beforeFileUpload: UploadProps['beforeUpload'] = async (rawFile) => {
        return await new Promise<boolean>((resolve, reject) => {
            //Judge the size
            if (rawFile.size /1024 /1024 > 2) {
                Swal.fire({
                    title: "Cover size cannot be larger than 2M",
                    heightAuto: false,
                    icon: "error",
})
                reject(false);
            }
            //Judge size
            let reader = new FileReader();
            reader.readAsDataURL(rawFile);
            reader.onload = () => {
                //Let the src of the img tag in the page point to the read path
                let img = new Image();
                img.src = reader.result as string;
                img.onload = () => {
                    console.log(img.width);
                    console.log(img.height);
if (img.width < 960 || img.height < 600) {
                        Swal.fire({
                            title: "Please upload pictures above 960*600 size",
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

    const RedefineUploadFile = async (params: UploadRequestOptions) => {
        try {
            const response = await uploadFile(liveInformationForm, params.file)
            liveInformationForm.uploadUrl = response.path
            console.log(response)
        } catch (err) {
            console.log(err)
            Swal.fire({
                title: "Failed to obtain upload node",
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

export const useSaveData = async (liveInformationForm: LiveInformation, formEl: FormInstance | undefined, rawData: GetLiveDataRes) => {
    if (!formEl) return;

    await formEl.validate(async (valid, fields) => {
        if (valid) {
try {
                if (liveInformationForm.uploadUrl == rawData.img && liveInformationForm.title == rawData.title) throw "Unmodified information";
                if (!liveInformationForm.uploadUrl) throw "Please upload pictures first"
                const requistData = <SaveLiveDataReq>{
                    type: liveInformationForm.uploadType,
                    imgUrl: liveInformationForm.uploadUrl,
                    title: liveInformationForm.title,
                }
const data = await saveLiveDataRequist(requistData)
                console.log(data)
                Swal.fire({
                    title: "Modification successful",
                    confirmButtonColor: globalScss.colorButtonTheme,
                    heightAuto: false,
                    icon: "success",

                })
                console.log("Upload successful")
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

export const useInit = async (liveInformationForm: LiveInformation, rawData: GetLiveDataRes) => {
    try {
        //Get user information
        const data = (await getLiveDataRequist()).data as GetLiveDataRes;
liveInformationForm.FileUrl = data.img
        const imgPathInfo = getLocation(data.img)

        //How to return the full path from the backend to get the path after domain name
        if (imgPathInfo?.pathname) {
            let pathname = imgPathInfo?.pathname.slice(1)
            liveInformationForm.uploadUrl = pathname
            rawData.img = pathname
        }
        //Save original data
        rawData.title = data.title
        liveInformationForm.title = data.title
        //Get the request method of the current interface
const updataMenhod = (await getuploadingMethod(<GetUploadingMethodReq>{
            method: liveInformationForm.interface
        })).data as GetUploadingMethodRes
        liveInformationForm.uploadType = updataMenhod.type
        console.log(updataMenhod)

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
export const useCopy = async (text: string) => {
    try {
        const { toClipboard } = useClipboard();
        await toClipboard(text); //Realize copying
        Swal.fire({
            title: "Copy successfully",
            confirmButtonColor: globalScss.colorButtonTheme,
            heightAuto: false,
            icon: "success",

        })
    } catch (e) {
        console.error(e);
    }
};

//form validation
export const useRules = () => {
    const liveInformationRules = reactive({
title: [{ validator: validateLiveTitle, trigger: 'change' }],
    });
    return {
        liveInformationRules,
    };
};