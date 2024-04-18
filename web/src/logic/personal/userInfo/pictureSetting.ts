import { getuploadingMethod } from "@/apis/commonality";
import { updateAvatarRequist } from "@/apis/personal";
import globalScss from "@/assets/styles/global/export.module.scss";
import { useUserStore } from "@/store/main";
import { GetUploadingMethodReq, GetUploadingMethodRes } from "@/types/commonality/commonality";
import { uploadAvatar } from "@/types/personal/userInfo/pictureSetting";
import { UpdateAvatarReq } from "@/types/personal/userInfo/userInfo";
import { uploadFile } from "@/utils/upload/upload";
import type { UploadProps, UploadRequestOptions } from 'element-plus';
import Swal from 'sweetalert2';
import { reactive } from "vue";

export const useAvatarProp = () => {
    const userStore = useUserStore()
    const uploadAvatarForm = reactive<uploadAvatar>({
        FileUrl: '',
        uploadUrl: "",
        interface: "userAvatar",
        uploadType: "",
        action: "#",
        progress: 0
    });

    return {
        userStore,
        uploadAvatarForm
    }
}


export const useHandleFileMethod = (uploadAvatarForm: uploadAvatar) => {

    const handleFileSuccess: UploadProps['onSuccess'] = (
        response,
        uploadFile
    ) => {
        uploadAvatarForm.FileUrl = URL.createObjectURL(uploadFile.raw!)
    }


    const handleFileError: UploadProps['onError'] = (
        response,
    ) => {
        console.log("upload failed")
        Swal.fire({
            title: "upload failed",
            heightAuto: false,
            icon: "error",

        })
        console.log(response)

    }

    const beforeFileUpload: UploadProps['beforeUpload'] = (rawFile) => {
        if (rawFile.size / 1024 / 1024 > 2) {
            Swal.fire({
                title: "File does not conform to the format",
                heightAuto: false,
                icon: "error",

            })
            return false
        }
        console.log(rawFile.text)
        return true
    }


    const redefineUploadFile = async (params: UploadRequestOptions) => {
        try {
            const response = await uploadFile(uploadAvatarForm, params.file)
            uploadAvatarForm.uploadUrl = response.path
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
        redefineUploadFile
    }

}

export const useInit = async (uploadAvatarForm: uploadAvatar) => {
    try {
        const updataMenhod = (await getuploadingMethod(<GetUploadingMethodReq>{
            method: uploadAvatarForm.interface
        })).data as GetUploadingMethodRes
        uploadAvatarForm.uploadType = updataMenhod.type
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


export const useUpdateAvatar = async (store: any, uploadAvatarForm: uploadAvatar) => {
    try {
        if (!uploadAvatarForm.uploadUrl) throw "Please upload pictures first"
        const requistData = <UpdateAvatarReq>{
            imgUrl: uploadAvatarForm.uploadUrl,
            type: uploadAvatarForm.uploadType
        }
        const data = await updateAvatarRequist(requistData)
        console.log(data)
        uploadAvatarForm.FileUrl = ""
        store.userInfoData.photo = String(data.data) ?? ""
        Swal.fire({
            title: "Replacement successful",
            heightAuto: false,
            icon: "success",

        })
        console.log("Upload successful")
    } catch (err) {
        console.log(err)
        Swal.fire({
            title: "Please upload pictures",
            confirmButtonColor: globalScss.colorButtonTheme,
            heightAuto: false,
            icon: "warning",

        })
    }
} 